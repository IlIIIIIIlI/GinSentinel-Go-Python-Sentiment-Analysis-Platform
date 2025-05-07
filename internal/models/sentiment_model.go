package models

import (
	"time"

	"gorm.io/gorm"
)

// SentimentAnalysis 表示数据库中的情感分析记录
type SentimentAnalysis struct {
	ID        string             `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Text      string             `gorm:"type:text;not null" json:"text"`
	Sentiment string             `gorm:"type:varchar(20);not null" json:"sentiment"`     // positive, negative, neutral
	Score     float64            `gorm:"type:decimal(5,4);not null" json:"score"`        // 范围通常为 -1.0 到 1.0
	UserID    string             `gorm:"type:varchar(50);index" json:"user_id"`          // 可选的用户标识
	Metadata  []AnalysisMetadata `gorm:"foreignKey:AnalysisID" json:"metadata"`          // 关联的元数据
	Language  string             `gorm:"type:varchar(10)" json:"language"`               // 语言代码（例如，"en", "zh"）
	Keywords  string             `gorm:"type:text" json:"keywords"`                      // 逗号分隔的关键词
	RequestID string             `gorm:"type:varchar(50);uniqueIndex" json:"request_id"` // 唯一请求标识
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}

// AnalysisMetadata 表示与分析关联的元数据
type AnalysisMetadata struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AnalysisID string         `gorm:"type:uuid;not null;index" json:"analysis_id"` // 引用 SentimentAnalysis.ID
	Key        string         `gorm:"type:varchar(50);not null" json:"key"`
	Value      string         `gorm:"type:text;not null" json:"value"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// BatchAnalysis 表示批处理请求
type BatchAnalysis struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:varchar(50);index" json:"user_id"`
	Count     int            `gorm:"type:int;not null" json:"count"`          // 批处理中的分析数量
	Status    string         `gorm:"type:varchar(20);not null" json:"status"` // pending, completed, failed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BatchItem 表示批处理请求中的单个项目
type BatchItem struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	BatchID    string         `gorm:"type:uuid;not null;index" json:"batch_id"` // 引用 BatchAnalysis.ID
	AnalysisID string         `gorm:"type:uuid;not null" json:"analysis_id"`    // 引用 SentimentAnalysis.ID
	Order      int            `gorm:"type:int;not null" json:"order"`           // 批处理中的顺序
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 覆盖 SentimentAnalysis 的表名
func (SentimentAnalysis) TableName() string {
	return "sentiment_analyses"
}

// TableName 覆盖 AnalysisMetadata 的表名
func (AnalysisMetadata) TableName() string {
	return "analysis_metadata"
}

// TableName 覆盖 BatchAnalysis 的表名
func (BatchAnalysis) TableName() string {
	return "batch_analyses"
}

// TableName 覆盖 BatchItem 的表名
func (BatchItem) TableName() string {
	return "batch_items"
}

// 以下是服务层使用的结构体，不映射到数据库

// SentimentResult 包含情感分析的结果
type SentimentResult struct {
	Text             string
	Sentiment        string
	Score            float64
	ConfidenceScores map[string]float64
	Keywords         []string
	RequestID        string
	Timestamp        time.Time
}

// BatchSentimentResult 包含多个情感分析的结果
type BatchSentimentResult struct {
	Results []SentimentResult
	BatchID string
}

// AnalysisHistoryResult 包含历史情感分析
type AnalysisHistoryResult struct {
	Records    []AnalysisRecord
	TotalCount int
}

// AnalysisRecord 表示存储的情感分析记录
type AnalysisRecord struct {
	ID        string
	Text      string
	Sentiment string
	Score     float64
	Timestamp time.Time
	Metadata  map[string]string
}
