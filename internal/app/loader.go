package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"

	"sentiment-service/internal/api/v1"
	"sentiment-service/internal/app/config"
	"sentiment-service/internal/app/initializer"
	"sentiment-service/internal/middleware"
)

const (
	banner = `
███████╗███████╗███╗   ██╗████████╗██╗███╗   ███╗███████╗███╗   ██╗████████╗
██╔════╝██╔════╝████╗  ██║╚══██╔══╝██║████╗ ████║██╔════╝████╗  ██║╚══██╔══╝
███████╗█████╗  ██╔██╗ ██║   ██║   ██║██╔████╔██║█████╗  ██╔██╗ ██║   ██║   
╚════██║██╔══╝  ██║╚██╗██║   ██║   ██║██║╚██╔╝██║██╔══╝  ██║╚██╗██║   ██║   
███████║███████╗██║ ╚████║   ██║   ██║██║ ╚═╝ ██║███████╗██║ ╚████║   ██║   
╚══════╝╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝   
    `
)

// Start 启动服务
func Start() error {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("配置文件加载错误: %v", err)
	}

	// 从环境变量获取配置并覆盖配置文件中的值
	loadEnvironmentVariables()

	// 记录关键配置值
	logrus.WithFields(logrus.Fields{
		"algorithm_endpoint": config.Conf.Algorithm.Endpoint,
		"rabbitmq_url":       config.Conf.RabbitMQ.URL,
		"db_host":            config.Conf.Database.Host,
	}).Info("加载的配置信息")

	// 初始化所有模块
	if err := InitializeAll(); err != nil {
		return fmt.Errorf("模块初始化错误: %v", err)
	}

	// 设置Gin模式
	if config.Conf.App.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 应用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// 设置路由
	v1.SetupRoutes(r, initializer.DB)

	// 打印启动信息
	printStartupInfo()

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Conf.App.Addr, config.Conf.App.Port)
	logrus.Infof("服务器启动，监听地址 %s", addr)

	return r.Run(addr)
}

// loadEnvironmentVariables 从环境变量获取配置
func loadEnvironmentVariables() {
	// 数据库配置
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Conf.Database.Host = dbHost
		logrus.Infof("从环境变量加载数据库主机: %s", dbHost)
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		// 这里可以添加端口转换为整数的逻辑
		logrus.Infof("从环境变量加载数据库端口: %s", dbPort)
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Conf.Database.User = dbUser
		logrus.Infof("从环境变量加载数据库用户: %s", dbUser)
	}

	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Conf.Database.Password = dbPassword
		logrus.Info("从环境变量加载数据库密码")
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Conf.Database.DBName = dbName
		logrus.Infof("从环境变量加载数据库名称: %s", dbName)
	}

	// 算法服务端点
	if endpoint := os.Getenv("ALGORITHM_ENDPOINT"); endpoint != "" {
		config.Conf.Algorithm.Endpoint = endpoint
		logrus.Infof("从环境变量加载算法服务端点: %s", endpoint)
	}

	// RabbitMQ配置
	if rabbitMQURL := os.Getenv("RABBITMQ_URL"); rabbitMQURL != "" {
		config.Conf.RabbitMQ.URL = rabbitMQURL
		logrus.Infof("从环境变量加载RabbitMQ URL: %s", rabbitMQURL)
	}
}

// InitializeAll 初始化所有模块
func InitializeAll() error {
	// 初始化日志
	if err := initializer.InitializeLogger(); err != nil {
		return fmt.Errorf("日志初始化错误: %v", err)
	}

	// 初始化数据库
	if err := initializer.InitializeDB(); err != nil {
		return fmt.Errorf("数据库初始化错误: %v", err)
	}

	// 初始化Redis（如果需要）
	if err := initializer.InitializeRedis(); err != nil {
		logrus.Warnf("Redis初始化错误（非致命）: %v", err)
		// 这里不返回错误，因为Redis可能不是必需的
	}

	return nil
}

// 打印启动信息
func printStartupInfo() {
	fmt.Println(banner)
	fmt.Printf("\n情感分析服务 (Gin版本)\n")
	fmt.Printf("------------------------------------\n")
	fmt.Printf("模式: %s\n", config.Conf.App.Mode)
	fmt.Printf("地址: %s\n", config.Conf.App.Addr)
	fmt.Printf("端口: %d\n", config.Conf.App.Port)
	fmt.Printf("数据库: PostgreSQL @ %s:%d\n", config.Conf.Database.Host, config.Conf.Database.Port)
	fmt.Printf("数据库名: %s\n", config.Conf.Database.DBName)
	fmt.Printf("算法服务: %s\n", config.Conf.Algorithm.Endpoint)
	fmt.Printf("消息队列: %s\n", config.Conf.RabbitMQ.URL)
	fmt.Printf("------------------------------------\n\n")
}
