package initializer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"sentiment-service/internal/app/config"
	"sentiment-service/internal/models"
)

// DB 全局数据库连接
var DB *gorm.DB

// InitializeDB 初始化数据库连接
func InitializeDB() error {
	// 构建DSN
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.User,
		config.Conf.Database.Password,
		config.Conf.Database.DBName,
		config.Conf.Database.SSLMode,
		config.Conf.Database.TimeZone,
	)

	logrus.Debugf("数据库连接DSN: %s", dsn)

	// 配置GORM日志
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  getGormLogLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// 创建数据库连接
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 获取底层SQL DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("无法获取数据库实例: %v", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	logrus.Info("数据库连接成功")

	// 运行自动迁移
	if err := migrateDatabase(DB); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	logrus.Info("数据库初始化完成")
	return nil
}

// migrateDatabase 运行数据库迁移
func migrateDatabase(db *gorm.DB) error {
	logrus.Info("正在执行数据库迁移...")

	// 自动迁移模型
	err := db.AutoMigrate(
		&models.SentimentAnalysis{},
		&models.AnalysisMetadata{},
		&models.BatchAnalysis{},
		&models.BatchItem{},
	)

	if err != nil {
		return err
	}

	logrus.Info("数据库迁移完成")
	return nil
}

// getGormLogLevel 根据应用日志级别获取GORM日志级别
func getGormLogLevel() logger.LogLevel {
	switch config.Conf.Log.Level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Info
	}
}
