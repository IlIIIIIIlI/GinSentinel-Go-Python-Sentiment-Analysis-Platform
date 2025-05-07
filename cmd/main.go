package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"sentiment-service/internal/app"
)

func main() {
	// 捕获终止信号，用于优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// 启动服务（在单独的 goroutine 中）
	errChan := make(chan error, 1)
	go func() {
		if err := app.Start(); err != nil {
			errChan <- err
		}
	}()

	// 等待终止信号或错误
	select {
	case <-c:
		logrus.Info("收到终止信号，正在优雅关闭...")
	case err := <-errChan:
		logrus.WithError(err).Fatal("服务启动失败")
	}

	// 这里可以添加额外的清理代码
	logrus.Info("服务已关闭")
}
