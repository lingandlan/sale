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
	"marketplace/backend/pkg/mall"

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

	// 3.1 初始化 GORM 数据库（用于充值模块）
	gormDB, err := repository.NewGormDB(&cfg.Database)
	if err != nil {
		log.Fatal("connect gorm database failed", zap.Error(err))
	}
	log.Info("gorm database connected")

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
	rechargeRepo := repository.NewRechargeRepository(gormDB)

	// 7. 初始化 WSY 商城客户端
	wsyClient := mall.NewWSYClient(cfg.Mall)

	// 8. 初始化 Service
	authSvc := service.NewAuthService(&cfg.JWT, redisClient, userRepo)
	userSvc := service.NewUserService(userRepo)
	memberSvc := service.NewMemberService(wsyClient)
	rechargeSvc := service.NewRechargeService(rechargeRepo, memberSvc)

	// 9. 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	userHandler := handler.NewUserHandler(userSvc, memberSvc)
	adminHandler := handler.NewAdminHandler(userSvc, casbinSvc)
	rechargeHandler := handler.NewRechargeHandler(rechargeSvc, userRepo)

	// 10. 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)
	rbacMiddleware := middleware.NewRBACMiddleware(casbinSvc)

	// 11. 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 12. 创建 Gin 实例
	r := gin.New()

	// 13. 注册中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.ZapLogger())
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// 14. 注册路由
	router.SetupRouter(r, authHandler, userHandler, adminHandler, rechargeHandler, authMiddleware, rbacMiddleware, redisClient)

	// 15. 创建 HTTP Server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 16. 等待中断信号
	quit := make(chan os.Signal, 1)

	// 17. 启动服务器（goroutine）
	go func() {
		log.Info(fmt.Sprintf("server listening on :%d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", zap.Error(err))
			quit <- syscall.SIGTERM
		}
	}()

	// 18. 等待中断信号优雅关闭
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	// 17. 关闭数据库连接
	db.Close()

	if sqlGormDB, err := gormDB.DB(); err == nil && sqlGormDB != nil {
		sqlGormDB.Close()
	}

	// 18. 关闭 Redis 连接
	redisClient.Close()

	log.Info("server exited")
}
