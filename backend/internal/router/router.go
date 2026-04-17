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
	adminHandler *handler.AdminHandler,
	rechargeHandler *handler.RechargeHandler,
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
		// ========== 公开接口 ==========
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// ========== Dashboard接口 ==========
		dashboard := v1.Group("/dashboard")
		dashboard.Use(authMiddleware.Auth())
		{
			dashboard.GET("/statistics", rechargeHandler.GetDashboardStatistics)
			dashboard.GET("/todos", rechargeHandler.GetDashboardTodos)
			dashboard.GET("/recharge-trends", rechargeHandler.GetDashboardRechargeTrends)
		}

		// ========== B端充值申请 ==========
		bRecharge := v1.Group("/recharge/b-apply")
		bRecharge.Use(authMiddleware.Auth())
		{
			bRecharge.POST("", rechargeHandler.CreateBRechargeApplication)
		}

		// ========== B端充值审批 ==========
		bApproval := v1.Group("/recharge/b-approval")
		bApproval.Use(authMiddleware.Auth())
		{
			bApproval.GET("", rechargeHandler.GetRechargeApplicationList)
			bApproval.GET("/:id", rechargeHandler.GetRechargeApplicationDetail)
			bApproval.POST("/action", rechargeHandler.ApprovalRechargeApplication)
		}

		// ========== C端充值 ==========
		cRecharge := v1.Group("/recharge/c-entry")
		cRecharge.Use(authMiddleware.Auth())
		{
			cRecharge.POST("", rechargeHandler.CreateCRecharge)
			cRecharge.GET("", rechargeHandler.GetCRechargeList)
			cRecharge.GET("/search-member", userHandler.SearchMember)
			cRecharge.GET("/:id", rechargeHandler.GetCRechargeDetail)
		}

		// ========== 充值记录 ==========
		records := v1.Group("/recharge/records")
		records.Use(authMiddleware.Auth())
		{
			records.GET("", rechargeHandler.GetCRechargeList)
			records.GET("/:id", rechargeHandler.GetCRechargeDetail)
		}

		// ========== 门店卡管理 ==========
		card := v1.Group("/card")
		card.Use(authMiddleware.Auth())
		{
			card.GET("/verify/:cardNo", rechargeHandler.VerifyCard)
			card.POST("/consume", rechargeHandler.ConsumeCard)
			card.GET("/available", rechargeHandler.GetAvailableCards)
			card.GET("/list", rechargeHandler.GetCardList)
			card.GET("/detail/:cardNo", rechargeHandler.GetCardDetail)
			card.GET("/stats", rechargeHandler.GetCardStats)
			card.GET("/inventory-stats", rechargeHandler.GetCardInventoryStats)
			card.POST("/batch-import", rechargeHandler.BatchImportCards)
			card.POST("/allocate", rechargeHandler.AllocateCards)
			card.POST("/bind", rechargeHandler.BindCardToUser)
			card.POST("/:cardNo/freeze", rechargeHandler.FreezeCard)
			card.POST("/:cardNo/unfreeze", rechargeHandler.UnfreezeCard)
			card.POST("/:cardNo/void", rechargeHandler.VoidCard)
		}

		// ========== 充值中心管理 ==========
		center := v1.Group("/center")
		center.Use(authMiddleware.Auth())
		{
			center.GET("", rechargeHandler.GetCenters)
			center.GET("/:id", rechargeHandler.GetCenterDetail)
			center.POST("", rechargeHandler.CreateCenter)
			center.PUT("/:id", rechargeHandler.UpdateCenter)
			center.DELETE("/:id", rechargeHandler.DeleteCenter)
		}

		// ========== 操作员管理 ==========
		operator := v1.Group("/operator")
		operator.Use(authMiddleware.Auth())
		{
			operator.GET("", rechargeHandler.GetOperators)
			operator.POST("", rechargeHandler.CreateOperator)
			operator.PUT("/:id", rechargeHandler.UpdateOperator)
			operator.DELETE("/:id", rechargeHandler.DeleteOperator)
		}

		// ========== 系统设置 ==========
		system := v1.Group("/system")
		system.Use(authMiddleware.Auth(), rbacMiddleware.Auth())
		{
			system.GET("/config", func(c *gin.Context) {
				c.JSON(200, gin.H{"data": gin.H{
					"systemName":     "太积堂充值与门店管理系统",
					"rechargeRatio":  1,
					"minRecharge":    100,
					"cardExpiryDays": 365,
				}})
			})
			system.PUT("/config", func(c *gin.Context) {
				var req map[string]interface{}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"message": "配置已更新"})
			})
		}

		// ========== 用户管理 ==========
		user := v1.Group("/user")
		user.Use(authMiddleware.Auth(), rbacMiddleware.Auth())
		{
			user.GET("/info", userHandler.GetUserInfo)
			user.PUT("/info", userHandler.UpdateUserInfo)
			user.POST("/change-password", userHandler.ChangePassword)
		}

		// ========== 管理员接口 ==========
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.Auth(), rbacMiddleware.Auth())
		{
			// 用户管理
			admin.GET("/users", adminHandler.ListUsers)
			admin.POST("/users", adminHandler.CreateUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.POST("/users/:id/reset-password", adminHandler.ResetPassword)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)
		}
	}
}
