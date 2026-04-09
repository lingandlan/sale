package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

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

// Claims JWT Claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// Auth 认证中间件
func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing token")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid token format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		// 存入 Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireRole 角色鉴权中间件
func RequireRole(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "forbidden")
			c.Abort()
			return
		}

		role := roleVal.(int)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient permission")
		c.Abort()
	}
}

// GetUserID 获取当前用户 ID
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	return userID.(int64)
}

// GetUsername 获取当前用户名
func GetUsername(c *gin.Context) string {
	username, _ := c.Get("username")
	return username.(string)
}

// GetRole 获取当前用户角色
func GetRole(c *gin.Context) int {
	role, _ := c.Get("role")
	return role.(int)
}
