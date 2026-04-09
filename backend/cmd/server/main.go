package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/handler"
	"marketplace/backend/internal/middleware"
	"marketplace/backend/internal/rbac"
	"marketplace/backend/internal/repository"
	"marketplace/backend/internal/router"
	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	// 2. 初始化日志
	logCfg := &logger.Config{
		Mode:        cfg.Log.Mode,
		Level:       cfg.Log.Level,
		ServiceName: cfg.Log.ServiceName,
	}
	if err := logger.Init(logCfg); err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	log := logger.GetLogger()
	log.Info("starting server...")

	// 3. 初始化数据库
	db, err := repository.NewDB(&cfg.Database)
	if err != nil {
		log.Fatal("connect database failed", zap.Error(err))
	}
	log.Info("database connected")

	// 4. 初始化 Redis
	redisClient, err := repository.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatal("connect redis failed", zap.Error(err))
	}
	log.Info("redis connected")

	// 5. 初始化 Casbin RBAC
	casbinSvc, err := rbac.NewCasbinService(db, &cfg.Database)
	if err != nil {
		log.Fatal("init casbin failed", zap.Error(err))
	}
	log.Info("casbin initialized")

	// 6. 初始化 Repository
	userRepo := repository.NewUserRepository(db)

	// 6. 初始化 Service
	authSvc := service.NewAuthService(&cfg.JWT, redisClient, userRepo)
	userSvc := service.NewUserService(userRepo)

	// 7. 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	userHandler := handler.NewUserHandler(userSvc)

	// 8. 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)
	rbacMiddleware := middleware.NewRBACMiddleware(casbinSvc)

	// 9. 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 10. 创建 Gin 实例
	r := gin.New()

	// 11. 注册中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.ZapLogger())
	r.Use(middleware.CORS())

	// 13. 注册路由
	router.SetupRouter(r, authHandler, userHandler, authMiddleware, rbacMiddleware)

	// 13. 创建 HTTP Server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 14. 启动服务器（ goroutine ）
	go func() {
		log.Info(fmt.Sprintf("server listening on :%d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	// 15. 等待中断信号优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	// 16. 关闭数据库连接
	db.Close()

	// 17. 关闭 Redis 连接
	redisClient.Close()

	log.Info("server exited")
}
