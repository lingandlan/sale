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

				// 返回 500 - 不暴露内部错误详情
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{
					Code:    500,
					Message: "internal server error",
				})

				// 打印日志（含完整堆栈）
				fmt.Printf("[PANIC] %v\n%s\n", err, stack)

				c.Abort()
			}
		}()

		c.Next()
	}
}
