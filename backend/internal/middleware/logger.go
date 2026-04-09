package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"marketplace/backend/pkg/logger"
)

// ZapLogger Zap 日志中间件
func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		log := logger.GetLogger()
		if c.Writer.Status() >= 500 {
			log.Error("request error",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("cost", cost),
				zap.String("ip", c.ClientIP()),
				zap.Any("errors", c.Errors),
			)
		} else if c.Writer.Status() >= 400 {
			log.Warn("request warning",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("cost", cost),
				zap.String("ip", c.ClientIP()),
			)
		} else {
			log.Info("request",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("cost", cost),
				zap.String("ip", c.ClientIP()),
			)
		}
	}
}
