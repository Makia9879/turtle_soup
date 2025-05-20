package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatCompletionRequest struct {
	Messages    []map[string]string `json:"messages"`
	Model       string              `json:"model"`
	MaxTokens   int                 `json:"max_tokens"`
	Temperature float64             `json:"temperature"`
	Stream      bool                `json:"stream"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Client struct {
	apiKey string
	url    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		url:    "https://api.deepseek.com/v1/chat/completions",
	}
}

func (c *Client) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	req.Stream = false
	payload, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[DeepSeek] json.Marshal error: %v", err)
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(payload))
	if err != nil {
		logx.WithContext(ctx).Errorf("[DeepSeek] NewRequest error: %v", err)
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		logx.WithContext(ctx).Errorf("[DeepSeek] Do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logx.WithContext(ctx).Errorf("[DeepSeek] API error: %s, body: %s", resp.Status, string(body))
		return nil, errors.New("deepseek api error: " + resp.Status)
	}

	var result ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logx.WithContext(ctx).Errorf("[DeepSeek] Decode error: %v", err)
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("no response content received")
	}

	return &result, nil
}
