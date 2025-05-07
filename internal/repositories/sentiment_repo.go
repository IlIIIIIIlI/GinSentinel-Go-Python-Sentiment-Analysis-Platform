package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"sentiment-service/internal/models"
)

// SentimentRepository 定义了情感分析数据存储操作的接口
type SentimentRepository interface {
	// CreateAnalysis 创建一个新的情感分析记录
	CreateAnalysis(ctx context.Context, analysis *models.SentimentAnalysis) error

	// GetAnalysisById 根据ID获取情感分析记录
	GetAnalysisById(ctx context.Context, id string) (*models.SentimentAnalysis, error)

	// GetAnalysisByRequestId 根据请求ID获取情感分析记录
	GetAnalysisByRequestId(ctx context.Context, requestId string) (*models.SentimentAnalysis, error)

	// FindAnalyses 获取情感分析记录，可选过滤条件
	FindAnalyses(ctx context.Context, params FindAnalysesParams) ([]*models.SentimentAnalysis, int64, error)

	// CreateBatchAnalysis 创建一个新的批处理分析记录
	CreateBatchAnalysis(ctx context.Context, batch *models.BatchAnalysis) error

	// AddBatchItems 向批处理分析添加项目
	AddBatchItems(ctx context.Context, batchId string, analysisIds []string) error

	// UpdateBatchStatus 更新批处理分析的状态
	UpdateBatchStatus(ctx context.Context, batchId string, status string) error
}

// FindAnalysesParams 定义了搜索分析记录的参数
type FindAnalysesParams struct {
	UserID    string
	StartTime *time.Time
	EndTime   *time.Time
	Sentiment string
	Limit     int
	Offset    int
}

// sentimentRepository 实现了SentimentRepository接口
type sentimentRepository struct {
	db *gorm.DB
}

// NewSentimentRepository 创建一个新的情感分析仓库
func NewSentimentRepository(db *gorm.DB) SentimentRepository {
	return &sentimentRepository{db: db}
}

// CreateAnalysis 创建一个新的情感分析记录
func (r *sentimentRepository) CreateAnalysis(ctx context.Context, analysis *models.SentimentAnalysis) error {
	if analysis.ID == "" {
		analysis.ID = uuid.New().String()
	}

	if analysis.RequestID == "" {
		analysis.RequestID = uuid.New().String()
	}

	logrus.WithFields(logrus.Fields{
		"id":         analysis.ID,
		"request_id": analysis.RequestID,
		"sentiment":  analysis.Sentiment,
	}).Debug("创建情感分析记录")

	return r.db.WithContext(ctx).Create(analysis).Error
}

// GetAnalysisById 根据ID获取情感分析记录
func (r *sentimentRepository) GetAnalysisById(ctx context.Context, id string) (*models.SentimentAnalysis, error) {
	var analysis models.SentimentAnalysis

	err := r.db.WithContext(ctx).
		Preload("Metadata").
		First(&analysis, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithField("id", id).Debug("未找到情感分析记录")
			return nil, nil
		}
		return nil, err
	}

	return &analysis, nil
}

// GetAnalysisByRequestId 根据请求ID获取情感分析记录
func (r *sentimentRepository) GetAnalysisByRequestId(ctx context.Context, requestId string) (*models.SentimentAnalysis, error) {
	var analysis models.SentimentAnalysis

	err := r.db.WithContext(ctx).
		Preload("Metadata").
		First(&analysis, "request_id = ?", requestId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithField("request_id", requestId).Debug("未找到情感分析记录")
			return nil, nil
		}
		return nil, err
	}

	return &analysis, nil
}

// FindAnalyses 获取情感分析记录，可选过滤条件
func (r *sentimentRepository) FindAnalyses(ctx context.Context, params FindAnalysesParams) ([]*models.SentimentAnalysis, int64, error) {
	var analyses []*models.SentimentAnalysis
	var count int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(&models.SentimentAnalysis{})

	// 应用过滤器
	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}

	if params.StartTime != nil {
		query = query.Where("created_at >= ?", params.StartTime)
	}

	if params.EndTime != nil {
		query = query.Where("created_at <= ?", params.EndTime)
	}

	if params.Sentiment != "" {
		query = query.Where("sentiment = ?", params.Sentiment)
	}

	// 获取总数
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	logrus.WithFields(logrus.Fields{
		"user_id":   params.UserID,
		"sentiment": params.Sentiment,
		"count":     count,
	}).Debug("查询情感分析记录")

	// 应用分页
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// 执行查询并预加载
	err = query.
		Preload("Metadata").
		Order("created_at DESC").
		Find(&analyses).Error

	if err != nil {
		return nil, 0, err
	}

	return analyses, count, nil
}

// CreateBatchAnalysis 创建一个新的批处理分析记录
func (r *sentimentRepository) CreateBatchAnalysis(ctx context.Context, batch *models.BatchAnalysis) error {
	if batch.ID == "" {
		batch.ID = uuid.New().String()
	}

	if batch.Status == "" {
		batch.Status = "pending"
	}

	logrus.WithFields(logrus.Fields{
		"id":     batch.ID,
		"count":  batch.Count,
		"status": batch.Status,
	}).Debug("创建批处理分析记录")

	return r.db.WithContext(ctx).Create(batch).Error
}

// AddBatchItems 向批处理分析添加项目
func (r *sentimentRepository) AddBatchItems(ctx context.Context, batchId string, analysisIds []string) error {
	// 使用事务确保所有项目都被原子性添加
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, analysisId := range analysisIds {
			item := models.BatchItem{
				BatchID:    batchId,
				AnalysisID: analysisId,
				Order:      i,
			}

			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		// 更新批处理计数
		return tx.Model(&models.BatchAnalysis{}).
			Where("id = ?", batchId).
			Update("count", len(analysisIds)).
			Error
	})
}

// UpdateBatchStatus 更新批处理分析的状态
func (r *sentimentRepository) UpdateBatchStatus(ctx context.Context, batchId string, status string) error {
	logrus.WithFields(logrus.Fields{
		"batch_id": batchId,
		"status":   status,
	}).Debug("更新批处理状态")

	return r.db.WithContext(ctx).
		Model(&models.BatchAnalysis{}).
		Where("id = ?", batchId).
		Update("status", status).
		Error
}
