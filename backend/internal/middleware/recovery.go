package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"marketplace/backend/pkg/logger"
	"marketplace/backend/pkg/response"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := debug.Stack()

				// 记录 panic 日志（含完整堆栈）
				logger.Error("panic recovered",
					zap.String("error", fmt.Sprintf("%v", err)),
					zap.String("stack", string(stack)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回 500 - 不暴露内部错误详情
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{
					Code:    500,
					Message: "internal server error",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
