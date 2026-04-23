package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"marketplace/backend/pkg/response"
)

// RateLimitConfig 限速配置
type RateLimitConfig struct {
	Limit  int           // 窗口内允许的最大请求数
	Window time.Duration // 滑动窗口大小
}

// RateLimitByKey 基于自定义 key 的滑动窗口限速中间件
func RateLimitByKey(rdb *redis.Client, keyFunc func(*gin.Context) string, cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s", keyFunc(c))
		ctx := c.Request.Context()

		now := time.Now().UnixNano()
		windowStart := now - cfg.Window.Nanoseconds()

		pipe := rdb.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
		countCmd := pipe.ZCard(ctx, key)
		pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
		pipe.Expire(ctx, key, cfg.Window)
		_, err := pipe.Exec(ctx)

		if err != nil {
			// Redis 故障时放行，不影响正常业务
			c.Next()
			return
		}

		if countCmd.Val() >= int64(cfg.Limit) {
			c.JSON(http.StatusTooManyRequests, response.ErrorResponse{
				Code:    429,
				Message: "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimit 登录限速：每 IP 每分钟 10 次
func LoginRateLimit(rdb *redis.Client) gin.HandlerFunc {
	return RateLimitByKey(rdb, func(c *gin.Context) string {
		return fmt.Sprintf("login:%s", c.ClientIP())
	}, RateLimitConfig{
		Limit:  10,
		Window: time.Minute,
	})
}

// APIRateLimit 通用 API 限速：每用户每分钟 60 次
func APIRateLimit(rdb *redis.Client) gin.HandlerFunc {
	return RateLimitByKey(rdb, func(c *gin.Context) string {
		userID, exists := c.Get("user_id")
		if !exists {
			return fmt.Sprintf("api:anon:%s", c.ClientIP())
		}
		return fmt.Sprintf("api:user:%v", userID)
	}, RateLimitConfig{
		Limit:  60,
		Window: time.Minute,
	})
}
