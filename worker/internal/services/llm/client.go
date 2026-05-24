package llm

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

var ErrNotConfigured = errors.New("LLMNotConfigured")
var client openai.Client

// NewClient 从配置读取 llm.endpoint / llm.api_key / llm.model 并返回客户端。
// 若 endpoint 或 api_key 未配置则返回 ErrNotConfigured。
func NewClient(endpoint, apiKey string) (openai.Client, error) {
	if endpoint == "" || apiKey == "" {
		return openai.Client{}, ErrNotConfigured
	}
	client = openai.NewClient(option.WithBaseURL(endpoint), option.WithAPIKey(apiKey))
	return client, nil
}

// Chat 发送 system + user 消息对，返回模型回复的文本内容。
func Chat(ctx context.Context, systemPrompt, userPrompt string, model string) (string, error) {
	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(systemPrompt),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(userPrompt),
		},
		Model: model,
	})

	if err != nil {
		return "", err
	}
	return resp.OutputText(), nil
}

// 流式获取
func ChatWithStream(ctx context.Context, systemPrompt, userPrompt string, handler func(event responses.ResponseStreamEventUnion)) error {
	stream := client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(userPrompt)},
		Model: openai.ChatModelGPT5_2,
	})
	for stream.Next() {
		event := stream.Current()
		handler(event)
	}

	if stream.Err() != nil {
		return stream.Err()
	}
	return nil
}
