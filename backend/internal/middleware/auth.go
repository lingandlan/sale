package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/errmsg"
	"marketplace/backend/pkg/response"
)

// AuthMiddleware JWT 认证中间件
type AuthMiddleware struct {
	jwtSecret []byte
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(jwtSecret)}
}

// Auth 认证中间件
func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, errmsg.Get("auth.token_missing"))
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, errmsg.Get("auth.token_format_error"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 Token（使用service.Claims结构）
		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, errmsg.Get("auth.token_invalid"))
			c.Abort()
			return
		}

		// 存入 Context
		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireRole 角色鉴权中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, errmsg.Get("common.forbidden"))
			c.Abort()
			return
		}

		role := roleVal.(string)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, errmsg.Get("common.forbidden"))
		c.Abort()
	}
}

// GetUserID 获取当前用户 ID
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	return userID.(int64)
}

// GetPhone 获取当前用户手机号
func GetPhone(c *gin.Context) string {
	phone, _ := c.Get("phone")
	return phone.(string)
}

// GetRole 获取当前用户角色
func GetRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}
