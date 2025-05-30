// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: sentiment/v1/sentiment.proto

package sentimentv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 情感分析请求
type SentimentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text      string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Language  string `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	RequestId string `protobuf:"bytes,3,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
}

func (x *SentimentRequest) Reset() {
	*x = SentimentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sentiment_v1_sentiment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SentimentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SentimentRequest) ProtoMessage() {}

func (x *SentimentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sentiment_v1_sentiment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SentimentRequest.ProtoReflect.Descriptor instead.
func (*SentimentRequest) Descriptor() ([]byte, []int) {
	return file_sentiment_v1_sentiment_proto_rawDescGZIP(), []int{0}
}

func (x *SentimentRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *SentimentRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *SentimentRequest) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

// 情感分析响应
type SentimentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId        string             `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Sentiment        string             `protobuf:"bytes,2,opt,name=sentiment,proto3" json:"sentiment,omitempty"`
	Score            float64            `protobuf:"fixed64,3,opt,name=score,proto3" json:"score,omitempty"`
	ConfidenceScores map[string]float64 `protobuf:"bytes,4,rep,name=confidence_scores,json=confidenceScores,proto3" json:"confidence_scores,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	Keywords         []string           `protobuf:"bytes,5,rep,name=keywords,proto3" json:"keywords,omitempty"`
}

func (x *SentimentResponse) Reset() {
	*x = SentimentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sentiment_v1_sentiment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SentimentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SentimentResponse) ProtoMessage() {}

func (x *SentimentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sentiment_v1_sentiment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SentimentResponse.ProtoReflect.Descriptor instead.
func (*SentimentResponse) Descriptor() ([]byte, []int) {
	return file_sentiment_v1_sentiment_proto_rawDescGZIP(), []int{1}
}

func (x *SentimentResponse) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *SentimentResponse) GetSentiment() string {
	if x != nil {
		return x.Sentiment
	}
	return ""
}

func (x *SentimentResponse) GetScore() float64 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *SentimentResponse) GetConfidenceScores() map[string]float64 {
	if x != nil {
		return x.ConfidenceScores
	}
	return nil
}

func (x *SentimentResponse) GetKeywords() []string {
	if x != nil {
		return x.Keywords
	}
	return nil
}

var File_sentiment_v1_sentiment_proto protoreflect.FileDescriptor

var file_sentiment_v1_sentiment_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x22, 0x61, 0x0a, 0x10,
	0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x22,
	0xab, 0x02, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x62, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08,
	0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x1a, 0x43, 0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xca, 0x01,
	0x0a, 0x11, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x6e, 0x61, 0x6c, 0x79,
	0x7a, 0x65, 0x72, 0x12, 0x55, 0x0a, 0x10, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x53, 0x65,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5e, 0x0a, 0x15, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x3d, 0x5a, 0x3b, 0x73, 0x65,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x65,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_sentiment_v1_sentiment_proto_rawDescOnce sync.Once
	file_sentiment_v1_sentiment_proto_rawDescData = file_sentiment_v1_sentiment_proto_rawDesc
)

func file_sentiment_v1_sentiment_proto_rawDescGZIP() []byte {
	file_sentiment_v1_sentiment_proto_rawDescOnce.Do(func() {
		file_sentiment_v1_sentiment_proto_rawDescData = protoimpl.X.CompressGZIP(file_sentiment_v1_sentiment_proto_rawDescData)
	})
	return file_sentiment_v1_sentiment_proto_rawDescData
}

var file_sentiment_v1_sentiment_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_sentiment_v1_sentiment_proto_goTypes = []interface{}{
	(*SentimentRequest)(nil),  // 0: sentiment.v1.SentimentRequest
	(*SentimentResponse)(nil), // 1: sentiment.v1.SentimentResponse
	nil,                       // 2: sentiment.v1.SentimentResponse.ConfidenceScoresEntry
}
var file_sentiment_v1_sentiment_proto_depIdxs = []int32{
	2, // 0: sentiment.v1.SentimentResponse.confidence_scores:type_name -> sentiment.v1.SentimentResponse.ConfidenceScoresEntry
	0, // 1: sentiment.v1.SentimentAnalyzer.AnalyzeSentiment:input_type -> sentiment.v1.SentimentRequest
	0, // 2: sentiment.v1.SentimentAnalyzer.BatchAnalyzeSentiment:input_type -> sentiment.v1.SentimentRequest
	1, // 3: sentiment.v1.SentimentAnalyzer.AnalyzeSentiment:output_type -> sentiment.v1.SentimentResponse
	1, // 4: sentiment.v1.SentimentAnalyzer.BatchAnalyzeSentiment:output_type -> sentiment.v1.SentimentResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_sentiment_v1_sentiment_proto_init() }
func file_sentiment_v1_sentiment_proto_init() {
	if File_sentiment_v1_sentiment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sentiment_v1_sentiment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SentimentRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sentiment_v1_sentiment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SentimentResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sentiment_v1_sentiment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sentiment_v1_sentiment_proto_goTypes,
		DependencyIndexes: file_sentiment_v1_sentiment_proto_depIdxs,
		MessageInfos:      file_sentiment_v1_sentiment_proto_msgTypes,
	}.Build()
	File_sentiment_v1_sentiment_proto = out.File
	file_sentiment_v1_sentiment_proto_rawDesc = nil
	file_sentiment_v1_sentiment_proto_goTypes = nil
	file_sentiment_v1_sentiment_proto_depIdxs = nil
}
