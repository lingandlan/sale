package router

import (
	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/handler"
	"marketplace/backend/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	rbacMiddleware *middleware.RBACMiddleware,
) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// 需要认证的接口（JWT + RBAC）
		user := v1.Group("/user")
		user.Use(authMiddleware.Auth(), rbacMiddleware.Auth())
		{
			user.GET("/info", userHandler.GetUserInfo)
			user.PUT("/info", userHandler.UpdateUserInfo)
			user.GET("/:id", userHandler.GetUserByID)
		}

		// 管理员接口（JWT + RBAC）
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.Auth(), rbacMiddleware.Auth())
		{
			admin.GET("/users", userHandler.ListUsers)
		}
	}
}
