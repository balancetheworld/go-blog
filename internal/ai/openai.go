package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/utils"
)

type OpenAIProvider struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

func NewOpenAIProvider(
	baseURL string,
	apiKey string,
	model string,
) *OpenAIProvider {
	return &OpenAIProvider{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func NewOpenAIProviderFromEnv() *OpenAIProvider {
	return NewOpenAIProvider(
		utils.Get(constant.EnvKeyAIModerationBaseURL),
		utils.Get(constant.EnvKeyAIModerationAPIKey),
		utils.Get(constant.EnvKeyAIModerationModel),
	)
}

type openAIChatRequest struct {
	Model          string               `json:"model"`
	Messages       []openAIMessage      `json:"messages"`
	Temperature    float64              `json:"temperature"`
	ResponseFormat openAIResponseFormat `json:"response_format"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatResponse struct {
	Choices []openAIChoice `json:"choices"`
}

type openAIChoice struct {
	Message openAIMessage `json:"message"`
}

// OpenAIProvider 实现 ModerationProvider 接口，调用OpenAI兼容大模型做内容审核
func (p *OpenAIProvider) Moderate(
	ctx context.Context,
	input ModerationInput,
) (ModerationResult, error) {
	// 1. 防御校验：接收器p不能为nil，内部HTTPClient不能是nil
	if err := p.ValidateConfig(); err != nil {
		return ModerationResult{}, err
	}
	// 3. 把待审核输入 ModerationInput 结构体序列化为JSON字节
	inputJSON, err := json.Marshal(input)
	if err != nil {
		// %w 包装错误，保留原始错误信息，上层可以用errors.Is解析
		return ModerationResult{}, fmt.Errorf("encode moderation input: %w", err)
	}

	//组装发给大模型的请求体 openAIChatRequest
	requestBody := openAIChatRequest{
		Model: p.Model,
		Messages: []openAIMessage{
			{
				Role:    "system", //system 消息是告诉模型：你现在是什么身份，要做什么，输出要遵守什么格式。
				Content: `Review the comment and return only a JSON object with status, categories, confidence, and reason. The JSON format is {"status":"approved","categories":[],"confidence":0.0,"reason":""}. status must be exactly one of: approved, rejected, manual_review. Approve by default. Reject only clearly identifiable advertising, violence, pornography, insults, or politically extremist content. Meaningless text, casual jokes, and mildly impolite but non-abusive content must be approved. Use manual_review only when the content may be disallowed but is genuinely ambiguous. confidence must be a number from 0 to 1. categories must contain at most 10 strings. reason must not exceed 500 characters.`,
			},
			{
				Role:    "user", //user 是交给模型处理的实际业务数据，这里就是待审核的内容 `inputJSON`。
				Content: string(inputJSON),
			},
		},
		Temperature: 0, // temperature=0，关闭随机性，输出尽量稳定、确定
		ResponseFormat: openAIResponseFormat{
			Type: "json_object", // 强制大模型输出合法JSON格式（OpenAI结构化输出能力）
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return ModerationResult{}, fmt.Errorf("encode moderation request: %w", err)
	}

	// 构建HTTP POST请求，带上上下文ctx
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.BaseURL+"/chat/completions", // openai兼容接口的聊天补全端点
		bytes.NewReader(body),         // body是已经json.Marshal好的请求字节流
	)
	if err != nil {
		return ModerationResult{}, fmt.Errorf("create moderation request: %w", err)
	}
	// 设置鉴权头：Bearer + apiKey，openai协议标准鉴权方式
	request.Header.Set("Authorization", "Bearer "+p.APIKey)
	// 请求体是json格式
	request.Header.Set("Content-Type", "application/json")

	// 使用注入的http客户端发送网络请求
	response, err := p.HTTPClient.Do(request)
	if err != nil {
		return ModerationResult{}, fmt.Errorf("send moderation request: %w", err)
	}
	// 非常重要：response.Body必须关闭，防止文件句柄泄露，defer函数，函数退出自动执行
	defer response.Body.Close()

	// 把http返回的body全部读进字节切片
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return ModerationResult{}, fmt.Errorf("read moderation response: %w", err)
	}

	// 判断HTTP状态码：2xx才认为成功；非2xx（401密钥错误、429限流、500服务报错）直接返回错误，同时把返回的响应体携带出来方便排查日志
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return ModerationResult{}, fmt.Errorf(
			"moderation provider returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	// 解析openai接口外层标准返回结构
	var chatResponse openAIChatResponse
	if err := json.Unmarshal(responseBody, &chatResponse); err != nil {
		return ModerationResult{}, fmt.Errorf("decode moderation response: %w", err)
	}

	// 保护判断：choices数组为空，模型没有返回任何结果，直接报错
	if len(chatResponse.Choices) == 0 {
		return ModerationResult{}, errors.New("moderation response has no choices")
	}

	// Choices[0].Message.Content 里面存放的就是模型输出的JSON字符串
	// 二次json解析，把模型输出json转成我们业务结构体 ModerationResult
	var result ModerationResult
	if err := json.Unmarshal(
		[]byte(chatResponse.Choices[0].Message.Content),
		&result,
	); err != nil {
		return ModerationResult{}, fmt.Errorf("decode moderation result: %w", err)
	}

	return result, nil
}

func (p *OpenAIProvider) ValidateConfig() error {
	if p == nil || p.HTTPClient == nil {
		return errors.New("openai provider is not configured")
	}
	if strings.TrimSpace(p.BaseURL) == "" ||
		strings.TrimSpace(p.APIKey) == "" ||
		strings.TrimSpace(p.Model) == "" {
		return errors.New("openai provider configuration is incomplete")
	}
	return nil
}
