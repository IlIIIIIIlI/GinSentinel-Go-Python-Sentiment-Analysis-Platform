package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `yaml:"app" mapstructure:"app"`
	Database  DatabaseConfig  `yaml:"database" mapstructure:"database"`
	Algorithm AlgorithmConfig `yaml:"algorithm" mapstructure:"algorithm"`
	Log       LogConfig       `yaml:"log" mapstructure:"log"`
	Redis     RedisConfig     `yaml:"redis" mapstructure:"redis"`
	RabbitMQ  RabbitMQConfig  `yaml:"rabbitmq" mapstructure:"rabbitmq"`
}

var Conf *Config

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 设置配置文件路径和名称
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 将配置文件内容解析到 Conf 变量中
	Conf = &Config{}
	err = viper.Unmarshal(Conf)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

type AppConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	Addr string `mapstructure:"addr"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver" mapstructure:"driver"`
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
	TimeZone string `yaml:"timezone" mapstructure:"timezone"`
}

type AlgorithmConfig struct {
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
}

type LogConfig struct {
	Format       string `yaml:"format" mapstructure:"format"`
	Level        string `yaml:"level" mapstructure:"level"`
	ReportCaller bool   `yaml:"reportCaller" mapstructure:"reportCaller"`
}

// Redis配置
type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

// RabbitMQ配置
type RabbitMQConfig struct {
	URL         string `yaml:"url" mapstructure:"url"`
	TaskQueue   string `yaml:"task_queue" mapstructure:"task_queue"`
	ResultQueue string `yaml:"result_queue" mapstructure:"result_queue"`
}
