package services

import (
	"context"
	"errors"
	"fmt"
	sentimentv1 "sentiment-service/internal/gen/sentiment/v1"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sentiment-service/internal/grpc"
	"sentiment-service/internal/models"
	"sentiment-service/internal/mq"
	"sentiment-service/internal/repositories"
)

// SentimentService 定义了情感分析服务的操作
type SentimentService struct {
	repository  repositories.SentimentRepository
	grpcClient  *grpc.SentimentClient
	mqClient    *mq.SentimentMQ
	callbackURL string
}

// NewSentimentService 创建情感分析服务
func NewSentimentService(
	repo repositories.SentimentRepository,
	grpcAddr string,
	amqpURL string,
	taskQueue string,
	resultQueue string,
) (*SentimentService, error) {
	// 记录初始化信息
	logrus.Infof("正在初始化情感分析服务...")
	logrus.Infof("gRPC地址: %s", grpcAddr)
	logrus.Infof("RabbitMQ地址: %s", amqpURL)
	logrus.Infof("任务队列: %s", taskQueue)
	logrus.Infof("结果队列: %s", resultQueue)

	// 创建gRPC客户端
	grpcClient, err := grpc.NewSentimentClient(grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("初始化gRPC客户端失败: %v", err)
	}
	logrus.Info("gRPC客户端初始化成功")

	// 创建消息队列客户端
	mqClient, err := mq.NewSentimentMQ(
		amqpURL,
		taskQueue,
		resultQueue,
	)
	if err != nil {
		grpcClient.Close()
		return nil, fmt.Errorf("初始化MQ客户端失败: %v", err)
	}
	logrus.Info("RabbitMQ客户端初始化成功")

	return &SentimentService{
		repository:  repo,
		grpcClient:  grpcClient,
		mqClient:    mqClient,
		callbackURL: "",
	}, nil
}

// AnalyzeSentiment 同步分析文本情感（使用gRPC）
func (s *SentimentService) AnalyzeSentiment(
	ctx context.Context,
	text string,
	language string,
	storeResult bool,
	metadata map[string]string,
) (*models.SentimentResult, error) {
	if text == "" {
		return nil, errors.New("文本不能为空")
	}

	logrus.WithFields(logrus.Fields{
		"text_length": len(text),
		"language":    language,
		"store":       storeResult,
	}).Debug("开始分析情感")

	// 生成请求ID
	requestID := uuid.New().String()

	// 调用gRPC服务
	response, _ := s.grpcClient.AnalyzeSentiment(
		ctx,
		text,
		language,
		requestID,
	)

	var result *models.SentimentResult

	// 转换gRPC响应到模型
	result = &models.SentimentResult{
		Text:             text,
		Sentiment:        response.Sentiment,
		Score:            response.Score,
		ConfidenceScores: response.ConfidenceScores,
		Keywords:         response.Keywords,
		RequestID:        response.RequestId,
		Timestamp:        time.Now(),
	}

	// 如果请求存储结果
	if storeResult {
		if err := s.storeAnalysisResult(ctx, result, language, metadata); err != nil {
			logrus.WithError(err).Error("存储分析结果失败")
		}
	}

	return result, nil
}

// AnalyzeSentimentAsync 异步分析文本情感（使用消息队列）
func (s *SentimentService) AnalyzeSentimentAsync(
	ctx context.Context,
	text string,
	language string,
	storeResult bool,
	metadata map[string]string,
) (string, error) {
	if text == "" {
		return "", errors.New("文本不能为空")
	}

	// 创建回调函数
	callback := func(result *models.SentimentResult) {
		if storeResult {
			// 使用背景上下文，因为回调可能在请求上下文结束后发生
			storeCtx := context.Background()
			if err := s.storeAnalysisResult(storeCtx, result, language, metadata); err != nil {
				logrus.WithError(err).Error("存储异步分析结果失败")
			}
		}
	}

	// 发布到消息队列
	requestID, err := s.mqClient.PublishTask(ctx, text, language, callback)
	if err != nil {
		return "", fmt.Errorf("发布异步任务失败: %v", err)
	}

	return requestID, nil
}

// BatchAnalyzeSentiment 批量分析多个文本的情感（使用gRPC流）
func (s *SentimentService) BatchAnalyzeSentiment(
	ctx context.Context,
	texts []string,
	language string,
	storeResults bool,
	metadata map[string]string,
) (*models.BatchSentimentResult, error) {
	if len(texts) == 0 {
		return nil, errors.New("文本列表不能为空")
	}

	logrus.WithFields(logrus.Fields{
		"text_count": len(texts),
		"language":   language,
		"store":      storeResults,
	}).Debug("开始批量分析情感")

	// 创建批处理ID
	batchID := uuid.New().String()

	// 创建批处理结果
	batchResult := &models.BatchSentimentResult{
		BatchID: batchID,
		Results: make([]models.SentimentResult, 0, len(texts)),
	}

	// 存储批处理记录（如果请求）
	var batchAnalysis *models.BatchAnalysis
	var analysisIDs []string

	if storeResults {
		// 创建批处理记录
		batchAnalysis = &models.BatchAnalysis{
			ID:     batchID,
			Count:  len(texts),
			Status: "pending",
		}

		// 添加用户ID（如果可用）
		if userID, ok := metadata["user_id"]; ok {
			batchAnalysis.UserID = userID
		}

		// 存储批处理
		if err := s.repository.CreateBatchAnalysis(ctx, batchAnalysis); err != nil {
			logrus.WithError(err).Error("存储批处理分析记录失败")
		} else {
			logrus.WithField("batch_id", batchID).Debug("批处理分析记录已存储")
		}
	}

	// 尝试使用gRPC流
	useGrpcStream := true
	stream, err := s.grpcClient.BatchAnalyzeStream(ctx)
	if err != nil {
		logrus.WithError(err).Warn("创建gRPC流失败，使用本地分析")
		useGrpcStream = false
	}

	if useGrpcStream {
		// 发送所有请求
		textToRequestID := make(map[string]string, len(texts))
		for _, text := range texts {
			requestID := uuid.New().String()
			textToRequestID[text] = requestID

			// 发送请求
			// 修复: 使用正确的SentimentRequest类型
			err := stream.Send(&sentimentv1.SentimentRequest{
				Text:      text,
				Language:  language,
				RequestId: requestID,
			})
			if err != nil {
				logrus.WithError(err).Error("发送流请求失败")
				useGrpcStream = false
				break
			}
		}

		// 如果所有请求都已发送
		if useGrpcStream {
			// 关闭发送方向
			if err := stream.CloseSend(); err != nil {
				logrus.WithError(err).Error("关闭发送流失败")
				useGrpcStream = false
			} else {
				// 接收所有响应
				for i := 0; i < len(texts); i++ {
					resp, err := stream.Recv()
					if err != nil {
						logrus.WithError(err).Error("接收流响应失败")
						useGrpcStream = false
						break
					}

					// 查找对应的文本
					var text string
					for t, reqID := range textToRequestID {
						if reqID == resp.RequestId {
							text = t
							break
						}
					}

					// 创建结果
					result := models.SentimentResult{
						Text:             text,
						Sentiment:        resp.Sentiment,
						Score:            resp.Score,
						ConfidenceScores: resp.ConfidenceScores,
						Keywords:         resp.Keywords,
						RequestID:        resp.RequestId,
						Timestamp:        time.Now(),
					}

					// 添加到批量结果
					batchResult.Results = append(batchResult.Results, result)

					// 存储结果（如果需要）
					if storeResults {
						analysisID, err := s.storeAnalysisResultForBatch(ctx, &result, language, metadata, batchID)
						if err != nil {
							logrus.WithError(err).Error("存储批量分析结果失败")
						} else {
							analysisIDs = append(analysisIDs, analysisID)
						}
					}
				}
			}
		}
	}

	// 存储批处理项目（如果我们有批处理ID和分析ID）
	if storeResults && len(analysisIDs) > 0 {
		if err := s.repository.AddBatchItems(ctx, batchID, analysisIDs); err != nil {
			logrus.WithError(err).Error("存储批处理项目失败")
		} else {
			// 更新批处理状态为已完成
			if err := s.repository.UpdateBatchStatus(ctx, batchID, "completed"); err != nil {
				logrus.WithError(err).Error("更新批处理状态失败")
			}
			logrus.WithField("batch_id", batchID).Debug("批处理项目已存储")
		}
	}

	return batchResult, nil
}

// GetAnalysisHistory 获取过去的情感分析历史
func (s *SentimentService) GetAnalysisHistory(
	ctx context.Context,
	userId string,
	startTime time.Time,
	endTime time.Time,
	limit int,
	offset int,
) (*models.AnalysisHistoryResult, error) {
	// 构建搜索参数
	params := repositories.FindAnalysesParams{
		UserID:    userId,
		StartTime: &startTime,
		EndTime:   &endTime,
		Limit:     limit,
		Offset:    offset,
	}

	logrus.WithFields(logrus.Fields{
		"user_id":    userId,
		"start_time": startTime,
		"end_time":   endTime,
		"limit":      limit,
		"offset":     offset,
	}).Debug("获取情感分析历史")

	// 从存储库获取分析
	analyses, count, err := s.repository.FindAnalyses(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建结果
	result := &models.AnalysisHistoryResult{
		Records:    make([]models.AnalysisRecord, len(analyses)),
		TotalCount: int(count),
	}

	// 转换每个分析
	for i, analysis := range analyses {
		// 转换元数据为映射
		metadata := make(map[string]string)
		for _, meta := range analysis.Metadata {
			metadata[meta.Key] = meta.Value
		}

		// 添加到结果
		result.Records[i] = models.AnalysisRecord{
			ID:        analysis.ID,
			Text:      analysis.Text,
			Sentiment: analysis.Sentiment,
			Score:     analysis.Score,
			Timestamp: analysis.CreatedAt,
			Metadata:  metadata,
		}
	}

	return result, nil
}

// 存储分析结果
func (s *SentimentService) storeAnalysisResult(
	ctx context.Context,
	result *models.SentimentResult,
	language string,
	metadata map[string]string,
) error {
	// 将关键词转换为逗号分隔的字符串
	keywordsStr := strings.Join(result.Keywords, ",")

	// 创建分析记录
	analysis := &models.SentimentAnalysis{
		Text:      result.Text,
		Sentiment: result.Sentiment,
		Score:     result.Score,
		Language:  language,
		Keywords:  keywordsStr,
		RequestID: result.RequestID,
	}

	// 添加用户ID（如果存在）
	if userID, ok := metadata["user_id"]; ok {
		analysis.UserID = userID
	}

	// 添加元数据
	for k, v := range metadata {
		analysis.Metadata = append(analysis.Metadata, models.AnalysisMetadata{
			Key:   k,
			Value: v,
		})
	}

	// 存储到数据库
	if err := s.repository.CreateAnalysis(ctx, analysis); err != nil {
		logrus.WithError(err).Error("存储情感分析记录失败")
		return err
	}

	logrus.WithField("analysis_id", analysis.ID).Debug("情感分析记录已存储")
	return nil
}

// 为批处理存储分析结果，并返回分析ID
func (s *SentimentService) storeAnalysisResultForBatch(
	ctx context.Context,
	result *models.SentimentResult,
	language string,
	metadata map[string]string,
	batchID string,
) (string, error) {
	// 将关键词转换为逗号分隔的字符串
	keywordsStr := strings.Join(result.Keywords, ",")

	// 创建分析记录
	analysis := &models.SentimentAnalysis{
		Text:      result.Text,
		Sentiment: result.Sentiment,
		Score:     result.Score,
		Language:  language,
		Keywords:  keywordsStr,
		RequestID: result.RequestID,
	}

	// 添加用户ID（如果存在）
	if userID, ok := metadata["user_id"]; ok {
		analysis.UserID = userID
	}

	// 添加元数据
	for k, v := range metadata {
		analysis.Metadata = append(analysis.Metadata, models.AnalysisMetadata{
			Key:   k,
			Value: v,
		})
	}

	// 添加批处理ID作为元数据
	analysis.Metadata = append(analysis.Metadata, models.AnalysisMetadata{
		Key:   "batch_id",
		Value: batchID,
	})

	// 存储到数据库
	if err := s.repository.CreateAnalysis(ctx, analysis); err != nil {
		logrus.WithError(err).Error("存储情感分析记录失败")
		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"analysis_id": analysis.ID,
		"batch_id":    batchID,
	}).Debug("批处理情感分析记录已存储")

	return analysis.ID, nil
}
