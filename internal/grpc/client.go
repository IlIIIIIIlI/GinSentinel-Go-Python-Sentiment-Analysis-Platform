package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// This import path will be updated by Buf's code generation
	pb "sentiment-service/internal/gen/sentiment/v1"
)

// SentimentClient 是情感分析gRPC客户端
type SentimentClient struct {
	conn   *grpc.ClientConn
	client pb.SentimentAnalyzerClient
}

// NewSentimentClient 创建新的gRPC客户端
func NewSentimentClient(serverAddr string) (*SentimentClient, error) {
	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 建立连接
	conn, err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("无法连接到gRPC服务器: %v", err)
	}

	client := pb.NewSentimentAnalyzerClient(conn)

	return &SentimentClient{
		conn:   conn,
		client: client,
	}, nil
}

// AnalyzeSentiment 使用gRPC调用分析文本
func (c *SentimentClient) AnalyzeSentiment(ctx context.Context, text, language, requestID string) (*pb.SentimentResponse, error) {
	request := &pb.SentimentRequest{
		Text:      text,
		Language:  language,
		RequestId: requestID,
	}

	return c.client.AnalyzeSentiment(ctx, request)
}

// BatchAnalyzeStream 创建批量分析流
func (c *SentimentClient) BatchAnalyzeStream(ctx context.Context) (pb.SentimentAnalyzer_BatchAnalyzeSentimentClient, error) {
	return c.client.BatchAnalyzeSentiment(ctx)
}

// Close 关闭gRPC连接
func (c *SentimentClient) Close() error {
	return c.conn.Close()
}
