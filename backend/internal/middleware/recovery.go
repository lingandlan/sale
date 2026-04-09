package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"marketplace/backend/pkg/response"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := debug.Stack()

				// 返回 500
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{
					Code:    500,
					Message: "internal server error",
					Detail:  fmt.Sprintf("%v", err),
				})

				// 打印日志
				// logger.Error("panic recovered",
				// 	logger.String("error", fmt.Sprintf("%v", err)),
				// 	logger.String("stack", string(stack)),
				// )

				_ = stack // 避免未使用警告
				c.Abort()
			}
		}()

		c.Next()
	}
}
