package initializer

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"sentiment-service/internal/app/config"
)

// InitializeLogger 初始化日志系统
func InitializeLogger() error {
	// 设置日志格式
	switch config.Conf.Log.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// 设置日志级别
	switch config.Conf.Log.Level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// 设置调用者报告
	logrus.SetReportCaller(config.Conf.Log.ReportCaller)

	// 确保日志目录存在
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 配置日志轮转
	logPath := filepath.Join(logDir, "sentiment-service.%Y%m%d.log")
	writer, err := rotatelogs.New(
		logPath,
		rotatelogs.WithLinkName(filepath.Join(logDir, "sentiment-service.log")),
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留7天
		rotatelogs.WithRotationTime(24*time.Hour), // 每天轮转一次
	)
	if err != nil {
		return err
	}

	// 设置多输出，同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, writer)
	logrus.SetOutput(multiWriter)

	logrus.Info("日志系统初始化完成")
	return nil
}
