package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"sentiment-service/internal/models"
)

// SentimentMQ 管理RabbitMQ情感分析队列
type SentimentMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	// 任务队列名称
	taskQueue string

	// 结果队列名称
	resultQueue string

	// 结果回调
	resultCallbacks map[string]ResultCallback
}

// ResultCallback 结果回调函数类型
type ResultCallback func(*models.SentimentResult)

// NewSentimentMQ 创建新的RabbitMQ客户端
func NewSentimentMQ(amqpURL, taskQueue, resultQueue string) (*SentimentMQ, error) {
	// 记录连接信息
	logrus.Infof("正在连接到RabbitMQ: %s", amqpURL)

	// 尝试多次连接，处理RabbitMQ可能启动较慢的情况
	var conn *amqp.Connection
	var err error
	maxAttempts := 5

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// 连接RabbitMQ
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			break
		}
		logrus.Warnf("连接到RabbitMQ失败(尝试 %d/%d): %v", attempt, maxAttempts, err)

		// 最后一次尝试失败，返回错误
		if attempt == maxAttempts {
			return nil, fmt.Errorf("连接到RabbitMQ失败: %v", err)
		}

		// 等待后重试
		time.Sleep(time.Duration(attempt) * 2 * time.Second)
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建通道失败: %v", err)
	}

	// 声明队列
	_, err = channel.QueueDeclare(
		taskQueue, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 非阻塞
		nil,       // 参数
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("声明任务队列失败: %v", err)
	}

	_, err = channel.QueueDeclare(
		resultQueue, // 队列名称
		true,        // 持久化
		false,       // 自动删除
		false,       // 独占
		false,       // 非阻塞
		nil,         // 参数
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("声明结果队列失败: %v", err)
	}

	logrus.Info("成功连接到RabbitMQ并声明队列")

	mq := &SentimentMQ{
		conn:            conn,
		channel:         channel,
		taskQueue:       taskQueue,
		resultQueue:     resultQueue,
		resultCallbacks: make(map[string]ResultCallback),
	}

	// 启动结果消费者
	go mq.consumeResults()

	return mq, nil
}

// PublishTask 发布情感分析任务
func (mq *SentimentMQ) PublishTask(
	ctx context.Context,
	text,
	language string,
	callback ResultCallback,
) (string, error) {
	// 生成请求ID
	requestID := uuid.New().String()

	// 创建任务
	task := map[string]interface{}{
		"text":       text,
		"language":   language,
		"request_id": requestID,
		"timestamp":  time.Now().Unix(),
	}

	// 序列化任务
	body, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("序列化任务失败: %v", err)
	}

	// 发布消息
	err = mq.channel.Publish(
		"",           // 交换机
		mq.taskQueue, // 路由键
		false,        // 强制
		false,        // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return "", fmt.Errorf("发布任务失败: %v", err)
	}

	// 注册回调
	if callback != nil {
		mq.resultCallbacks[requestID] = callback
	}

	return requestID, nil
}

// consumeResults 消费结果队列
func (mq *SentimentMQ) consumeResults() {
	// 消费消息
	msgs, err := mq.channel.Consume(
		mq.resultQueue, // 队列
		"",             // 消费者
		true,           // 自动应答
		false,          // 独占
		false,          // 不等待
		false,          // 参数
		nil,            // 参数
	)
	if err != nil {
		logrus.Errorf("开始消费结果失败: %v", err)
		return
	}

	// 处理消息
	for msg := range msgs {
		var result map[string]interface{}
		if err := json.Unmarshal(msg.Body, &result); err != nil {
			logrus.Errorf("解析结果消息失败: %v", err)
			continue
		}

		requestID, ok := result["request_id"].(string)
		if !ok {
			logrus.Error("结果消息缺少request_id字段")
			continue
		}

		// 转换为结果模型
		sentimentResult := convertToSentimentResult(result)

		// 调用回调函数
		if callback, exists := mq.resultCallbacks[requestID]; exists {
			callback(sentimentResult)
			delete(mq.resultCallbacks, requestID)
		}
	}
}

// 转换为模型
func convertToSentimentResult(data map[string]interface{}) *models.SentimentResult {
	result := &models.SentimentResult{
		RequestID: data["request_id"].(string),
		Timestamp: time.Now(),
	}

	if text, ok := data["text"].(string); ok {
		result.Text = text
	}

	if sentiment, ok := data["sentiment"].(string); ok {
		result.Sentiment = sentiment
	}

	if score, ok := data["score"].(float64); ok {
		result.Score = score
	}

	if confScores, ok := data["confidence_scores"].(map[string]interface{}); ok {
		result.ConfidenceScores = make(map[string]float64)
		for k, v := range confScores {
			if score, ok := v.(float64); ok {
				result.ConfidenceScores[k] = score
			}
		}
	}

	if keywords, ok := data["keywords"].([]interface{}); ok {
		result.Keywords = make([]string, 0, len(keywords))
		for _, v := range keywords {
			if keyword, ok := v.(string); ok {
				result.Keywords = append(result.Keywords, keyword)
			}
		}
	}

	return result
}

// Close 关闭连接
func (mq *SentimentMQ) Close() error {
	if mq.channel != nil {
		mq.channel.Close()
	}
	if mq.conn != nil {
		return mq.conn.Close()
	}
	return nil
}
