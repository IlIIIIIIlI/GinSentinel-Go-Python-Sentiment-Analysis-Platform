package middleware

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger 记录HTTP请求的日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 创建请求ID（如果不存在）
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Header("X-Request-ID", requestID)
		}

		// 准备日志字段
		logger := logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"client_ip":  c.ClientIP(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"user_agent": c.Request.UserAgent(),
			"referer":    c.Request.Referer(),
		})

		// 记录请求开始
		logger.Info("API请求开始")

		// 继续处理请求
		c.Next()

		// 记录请求结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 添加状态码和延迟
		logger = logger.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"latency_ms":  float64(latency.Nanoseconds()) / 1_000_000.0,
			"error_count": len(c.Errors),
		})

		// 根据状态码确定日志级别
		statusCode := c.Writer.Status()
		switch {
		case statusCode >= 500:
			logger.Error("API请求结束")
		case statusCode >= 400:
			logger.Warn("API请求结束")
		default:
			logger.Info("API请求结束")
		}
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 使用UUID替代自定义随机字符串生成方法
	// 这样更可靠，不会有索引越界的风险
	return uuid.New().String()
}

// 注意：不再使用有问题的randomString函数
// 以下是安全的随机字符串生成函数，仅作参考
func safeRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 创建一个安全的随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 预分配结果切片
	result := make([]byte, length)

	// 安全地生成随机字符串
	for i := 0; i < length; i++ {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}
