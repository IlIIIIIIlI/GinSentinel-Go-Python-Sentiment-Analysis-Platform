package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AppError 自定义应用错误
type AppError struct {
	HTTPStatus int    // HTTP 状态码（如 404）
	Code       int    // 业务错误码（如 40401）
	Message    string // 错误消息（如 "User not found"）
}

// Error 实现error接口
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError 创建一个新的应用错误
func NewAppError(httpStatus, code int, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// ErrorHandler 处理应用中的错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			var appErr *AppError

			if errors.As(err, &appErr) {
				// 自定义业务错误
				c.JSON(appErr.HTTPStatus, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
			} else {
				// 非自定义错误（如 Go 原生 error）
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "内部服务器错误",
				})
			}

			c.Abort()
		}
	}
}

// 预定义一些常见错误
var (
	ErrBadRequest = func(message string) *AppError {
		return NewAppError(http.StatusBadRequest, 400, message)
	}

	ErrUnauthorized = func(message string) *AppError {
		return NewAppError(http.StatusUnauthorized, 401, message)
	}

	ErrForbidden = func(message string) *AppError {
		return NewAppError(http.StatusForbidden, 403, message)
	}

	ErrNotFound = func(message string) *AppError {
		return NewAppError(http.StatusNotFound, 404, message)
	}

	ErrInternalServer = func(message string) *AppError {
		return NewAppError(http.StatusInternalServerError, 500, message)
	}
)
