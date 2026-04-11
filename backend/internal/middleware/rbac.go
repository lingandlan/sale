package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/rbac"
	"marketplace/backend/pkg/response"
)

// RBACMiddleware RBAC 权限中间件
type RBACMiddleware struct {
	casbinSvc *rbac.CasbinService
}

// NewRBACMiddleware 创建 RBAC 中间件
func NewRBACMiddleware(casbinSvc *rbac.CasbinService) *RBACMiddleware {
	return &RBACMiddleware{casbinSvc: casbinSvc}
}

// Auth RBAC 认证中间件
// 必须在 JWT Auth 中间件之后使用
func (m *RBACMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户 ID 和角色
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 超级管理员拥有所有权限
		role, _ := c.Get("role")
		if role.(string) == "super_admin" {
			c.Next()
			return
		}

		// 获取用户角色列表
		roles, err := m.casbinSvc.GetRolesForUser(formatUserID(userID))
		if err != nil {
			response.InternalError(c, "权限检查失败")
			c.Abort()
			return
		}

		// 检查是否拥有有效角色
		if len(roles) == 0 {
			response.Forbidden(c, "无有效角色")
			c.Abort()
			return
		}

		// 检查权限（支持多个角色）
		hasPermission := false
		for _, r := range roles {
			ok, err := m.casbinSvc.EnforcePath(r, method, path)
			if err != nil {
				continue
			}
			if ok {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoles 要求特定角色（可指定多个，满足其一即可）
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		for _, r := range roles {
			if userRole.(string) == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "角色不匹配")
		c.Abort()
		return
	}
}

// formatUserID 格式化用户 ID 为字符串
func formatUserID(userID interface{}) string {
	switch v := userID.(type) {
	case int64:
		return fmt.Sprintf("%d", v)
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	default:
		return ""
	}
}

// stringToRole 将字符串角色转为整型
func stringToRole(role string) int {
	switch role {
	case "admin":
		return 2
	case "merchant":
		return 1
	case "user":
		return 0
	default:
		return -1
	}
}

// roleToString 将整型角色转为字符串
func roleToString(role int) string {
	switch role {
	case 2:
		return "admin"
	case 1:
		return "merchant"
	case 0:
		return "user"
	default:
		return "unknown"
	}
}
