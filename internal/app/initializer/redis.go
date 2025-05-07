package initializer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"sentiment-service/internal/app/config"
)

// Redis 全局Redis客户端
var Redis *redis.Client

// InitializeRedis 初始化Redis连接
func InitializeRedis() error {
	// 创建Redis客户端
	addr := fmt.Sprintf("%s:%d", config.Conf.Redis.Host, config.Conf.Redis.Port)
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB, // 注意这里使用大写的DB
	})

	// 测试连接
	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接测试失败: %v", err)
	}

	logrus.Info("Redis初始化完成")
	return nil
}
