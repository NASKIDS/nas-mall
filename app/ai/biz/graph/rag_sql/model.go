package rag_sql

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel2' in graph 'text2sql'
func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	config := &ark.ChatModelConfig{
		Timeout:          nil,
		HTTPClient:       nil,
		RetryTimes:       nil,
		BaseURL:          "",
		Region:           "",
		APIKey:           "",
		AccessKey:        "",
		SecretKey:        "",
		Model:            "",
		MaxTokens:        nil,
		Temperature:      nil,
		TopP:             nil,
		Stop:             nil,
		FrequencyPenalty: nil,
		LogitBias:        nil,
		PresencePenalty:  nil,
		CustomHeader:     nil,
	}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
