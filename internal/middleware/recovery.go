package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery 从 panic 中恢复并记录堆栈跟踪
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈跟踪
				stack := debug.Stack()
				stackStr := string(stack)

				// 清理堆栈跟踪以方便日志记录
				stackTraceLines := strings.Split(stackStr, "\n")
				cleanedStack := []string{}

				// 只保留一些有用的行，避免日志过长
				for i, line := range stackTraceLines {
					if i < 15 { // 只保留前15行，通常包含最重要的信息
						cleanedStack = append(cleanedStack, line)
					}
				}

				// 如果堆栈跟踪很长，添加一个省略指示
				if len(stackTraceLines) > 15 {
					cleanedStack = append(cleanedStack, "... (more stack frames omitted)")
				}

				finalStack := strings.Join(cleanedStack, "\n")

				// 记录错误和堆栈跟踪
				logrus.WithFields(logrus.Fields{
					"error":       err,
					"method":      c.Request.Method,
					"path":        c.Request.URL.Path,
					"client_ip":   c.ClientIP(),
					"stack_trace": finalStack,
				}).Error("服务器发生 panic")

				// 构建友好的错误消息
				errorMessage := "服务器内部错误"
				if gin.Mode() == gin.DebugMode {
					// 在调试模式下提供更多详细信息
					errorMessage = fmt.Sprintf("服务器错误: %v", err)
				}

				// 返回统一的错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": errorMessage,
				})
			}
		}()

		c.Next()
	}
}
