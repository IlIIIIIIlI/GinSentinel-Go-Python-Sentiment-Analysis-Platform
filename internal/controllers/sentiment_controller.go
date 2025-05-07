package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"sentiment-service/internal/services"
)

// SentimentController 处理情感分析相关的HTTP请求
type SentimentController struct {
	sentimentService *services.SentimentService
}

// NewSentimentController 创建一个新的情感分析控制器
func NewSentimentController(sentimentService *services.SentimentService) *SentimentController {
	return &SentimentController{sentimentService: sentimentService}
}

// AnalyzeSentiment 分析单个文本的情感
// @Summary 分析单个文本的情感
// @Description 对提供的文本进行情感分析，返回情感得分和分类
// @Tags sentiment
// @Accept json
// @Produce json
// @Param request body AnalyzeSentimentRequest true "分析请求"
// @Success 200 {object} SentimentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sentiment/analyze [post]
func (sc *SentimentController) AnalyzeSentiment(c *gin.Context) {
	var request AnalyzeSentimentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "无效的请求参数: " + err.Error()})
		return
	}

	// 验证必要参数
	if request.Text == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "文本不能为空"})
		return
	}

	result, err := sc.sentimentService.AnalyzeSentiment(
		c.Request.Context(),
		request.Text,
		request.Language,
		request.StoreResult,
		request.Metadata,
	)
	if err != nil {
		logrus.WithError(err).Error("分析情感失败")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "处理请求失败"})
		return
	}

	// 构建响应
	response := SentimentResponse{
		Text:             result.Text,
		Sentiment:        result.Sentiment,
		Score:            result.Score,
		ConfidenceScores: result.ConfidenceScores,
		Keywords:         result.Keywords,
		RequestID:        result.RequestID,
		Timestamp:        result.Timestamp.Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// AnalyzeSentimentAsync 异步分析单个文本的情感
// @Summary 异步分析单个文本的情感
// @Description 异步处理提供的文本进行情感分析，返回请求ID用于后续查询
// @Tags sentiment
// @Accept json
// @Produce json
// @Param request body AnalyzeSentimentRequest true "分析请求"
// @Success 202 {object} AsyncResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sentiment/analyze/async [post]
func (sc *SentimentController) AnalyzeSentimentAsync(c *gin.Context) {
	var request AnalyzeSentimentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "无效的请求参数: " + err.Error()})
		return
	}

	// 验证必要参数
	if request.Text == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "文本不能为空"})
		return
	}

	// 调用异步分析服务
	requestID, err := sc.sentimentService.AnalyzeSentimentAsync(
		c.Request.Context(),
		request.Text,
		request.Language,
		request.StoreResult,
		request.Metadata,
	)
	if err != nil {
		logrus.WithError(err).Error("提交异步分析失败")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "处理请求失败"})
		return
	}

	// 构建响应
	response := AsyncResponse{
		RequestID: requestID,
		Status:    "pending",
		Message:   "分析请求已提交，正在处理中",
	}

	c.JSON(http.StatusAccepted, response)
}

// BatchAnalyzeSentiment 批量分析多个文本的情感
// @Summary 批量分析多个文本的情感
// @Description 对多个提供的文本进行情感分析，返回每个文本的情感分析结果
// @Tags sentiment
// @Accept json
// @Produce json
// @Param request body BatchAnalyzeSentimentRequest true "批量分析请求"
// @Success 200 {object} BatchSentimentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sentiment/batch [post]
func (sc *SentimentController) BatchAnalyzeSentiment(c *gin.Context) {
	var request BatchAnalyzeSentimentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "无效的请求参数: " + err.Error()})
		return
	}

	// 验证必要参数
	if len(request.Texts) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "文本列表不能为空"})
		return
	}

	result, err := sc.sentimentService.BatchAnalyzeSentiment(
		c.Request.Context(),
		request.Texts,
		request.Language,
		request.StoreResults,
		request.Metadata,
	)
	if err != nil {
		logrus.WithError(err).Error("批量分析情感失败")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "处理请求失败"})
		return
	}

	// 构建响应
	response := BatchSentimentResponse{
		BatchID: result.BatchID,
		Results: make([]SentimentResponse, len(result.Results)),
	}

	// 转换每个结果
	for i, res := range result.Results {
		response.Results[i] = SentimentResponse{
			Text:             res.Text,
			Sentiment:        res.Sentiment,
			Score:            res.Score,
			ConfidenceScores: res.ConfidenceScores,
			Keywords:         res.Keywords,
			RequestID:        res.RequestID,
			Timestamp:        res.Timestamp.Unix(),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetAnalysisHistory 获取情感分析历史记录
// @Summary 获取情感分析历史记录
// @Description 返回指定用户的情感分析历史记录
// @Tags sentiment
// @Accept json
// @Produce json
// @Param user_id query string false "用户ID"
// @Param start_time query int false "开始时间戳"
// @Param end_time query int false "结束时间戳"
// @Param limit query int false "每页结果数量" default(50)
// @Param offset query int false "分页偏移量" default(0)
// @Success 200 {object} GetAnalysisHistoryResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sentiment/history [get]
func (sc *SentimentController) GetAnalysisHistory(c *gin.Context) {
	// 获取查询参数
	userID := c.Query("user_id")

	// 解析时间参数
	startTimeStr := c.DefaultQuery("start_time", "0")
	endTimeStr := c.DefaultQuery("end_time", "0")

	startTime := time.Unix(parseTimestamp(startTimeStr), 0)
	endTime := time.Unix(parseTimestamp(endTimeStr), 0)

	// 解析分页参数
	limit := parseIntParam(c.DefaultQuery("limit", "50"))
	offset := parseIntParam(c.DefaultQuery("offset", "0"))

	// 调用服务获取历史记录
	result, err := sc.sentimentService.GetAnalysisHistory(
		c.Request.Context(),
		userID,
		startTime,
		endTime,
		limit,
		offset,
	)
	if err != nil {
		logrus.WithError(err).Error("获取情感分析历史失败")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "处理请求失败"})
		return
	}

	// 构建响应
	response := GetAnalysisHistoryResponse{
		TotalCount: result.TotalCount,
		Records:    make([]SentimentRecord, len(result.Records)),
	}

	// 转换每个记录
	for i, record := range result.Records {
		response.Records[i] = SentimentRecord{
			ID:        record.ID,
			Text:      record.Text,
			Sentiment: record.Sentiment,
			Score:     record.Score,
			Timestamp: record.Timestamp.Unix(),
			Metadata:  record.Metadata,
		}
	}

	c.JSON(http.StatusOK, response)
}

// 以下是请求和响应结构体定义

// AnalyzeSentimentRequest 单个文本情感分析请求
type AnalyzeSentimentRequest struct {
	Text        string            `json:"text" binding:"required"`
	Language    string            `json:"language"`
	StoreResult bool              `json:"store_result"`
	Metadata    map[string]string `json:"metadata"`
}

// BatchAnalyzeSentimentRequest 批量文本情感分析请求
type BatchAnalyzeSentimentRequest struct {
	Texts        []string          `json:"texts" binding:"required"`
	Language     string            `json:"language"`
	StoreResults bool              `json:"store_results"`
	Metadata     map[string]string `json:"metadata"`
}

// SentimentResponse 情感分析响应
type SentimentResponse struct {
	Text             string             `json:"text"`
	Sentiment        string             `json:"sentiment"`
	Score            float64            `json:"score"`
	ConfidenceScores map[string]float64 `json:"confidence_scores"`
	Keywords         []string           `json:"keywords"`
	RequestID        string             `json:"request_id"`
	Timestamp        int64              `json:"timestamp"`
}

// AsyncResponse 异步分析响应
type AsyncResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// BatchSentimentResponse 批量情感分析响应
type BatchSentimentResponse struct {
	Results []SentimentResponse `json:"results"`
	BatchID string              `json:"batch_id"`
}

// SentimentRecord 情感分析历史记录
type SentimentRecord struct {
	ID        string            `json:"id"`
	Text      string            `json:"text"`
	Sentiment string            `json:"sentiment"`
	Score     float64           `json:"score"`
	Timestamp int64             `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// GetAnalysisHistoryResponse 获取历史记录响应
type GetAnalysisHistoryResponse struct {
	Records    []SentimentRecord `json:"records"`
	TotalCount int               `json:"total_count"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}

// 辅助函数

// parseTimestamp 解析时间戳字符串为int64
func parseTimestamp(timestampStr string) int64 {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return 0
	}
	return timestamp
}

// parseIntParam 解析整数参数字符串为int
func parseIntParam(paramStr string) int {
	val, err := strconv.Atoi(paramStr)
	if err != nil || val < 0 {
		return 0
	}
	return val
}
