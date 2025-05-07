package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"sentiment-service/internal/app/config"
	"sentiment-service/internal/controllers"
	"sentiment-service/internal/repositories"
	"sentiment-service/internal/services"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 创建存储库
	repo := repositories.NewSentimentRepository(db)

	// 获取配置信息
	grpcEndpoint := config.Conf.Algorithm.Endpoint
	rabbitmqURL := config.Conf.RabbitMQ.URL
	taskQueue := config.Conf.RabbitMQ.TaskQueue
	resultQueue := config.Conf.RabbitMQ.ResultQueue

	// 记录配置信息
	logrus.WithFields(logrus.Fields{
		"grpc_endpoint": grpcEndpoint,
		"rabbitmq_url":  rabbitmqURL,
		"task_queue":    taskQueue,
		"result_queue":  resultQueue,
	}).Info("正在设置服务连接")

	// 创建服务
	service, err := services.NewSentimentService(
		repo,
		grpcEndpoint,
		rabbitmqURL,
		taskQueue,
		resultQueue,
	)
	if err != nil {
		logrus.Fatalf("初始化情感分析服务失败: %v", err)
	}

	// 创建控制器
	controller := controllers.NewSentimentController(service)

	// 设置API组
	api := r.Group("/api/v1")
	{
		// 情感分析API
		sentiment := api.Group("/sentiment")
		{
			// 同步分析接口（使用gRPC）
			sentiment.POST("/analyze", controller.AnalyzeSentiment)

			// 异步分析接口（使用RabbitMQ）
			sentiment.POST("/analyze/async", controller.AnalyzeSentimentAsync)

			// 批量分析接口（使用gRPC流）
			sentiment.POST("/batch", controller.BatchAnalyzeSentiment)

			// 历史记录查询
			sentiment.GET("/history", controller.GetAnalysisHistory)
		}

		// 健康检查API
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// API版本信息
		api.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"service": "Sentiment Analysis API",
				"version": "1.0.0",
				"engine":  "Gin Framework",
			})
		})
	}

	logrus.Info("API路由已设置")
}
