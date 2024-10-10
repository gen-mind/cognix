package ai

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

const ChatMessageRoleUser = "user"
const chatCompletionURL = "/chat/completions"

type ChatAI struct {
	client    *resty.Client
	apiKey    string
	modelName string
}
type ChatRequest struct {
	Model    string         `json:"model"`
	Messages []*ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Id                string        `json:"id"`
	Object            string        `json:"object"`
	Created           int           `json:"created"`
	Model             string        `json:"model"`
	SystemFingerprint string        `json:"system_fingerprint"`
	Choices           []*ChatChoice `json:"choices"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

func (c *ChatAI) Request(ctx context.Context, message string) (*Response, error) {
	request := &ChatRequest{
		Model: c.modelName,
		Messages: []*ChatMessage{
			{Role: ChatMessageRoleUser,
				Content: message,
			}},
	}
	response, err := c.client.R().SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.apiKey)).SetBody(request).Post(chatCompletionURL)
	if err = utils.WrapRestyError(response, err); err != nil {
		return nil, err
	}
	var chatResponse ChatResponse
	if err = json.Unmarshal(response.Body(), &chatResponse); err != nil {
		return nil, err
	}

	if len(chatResponse.Choices) == 0 {
		return nil, fmt.Errorf("no response from ai")
	}
	return &Response{Message: chatResponse.Choices[0].Message.Content}, nil
}

func NewChatAI(baseUrl, apiKey, modelName string) Client {
	return &ChatAI{
		client:    resty.New().SetTimeout(time.Minute).SetBaseURL(baseUrl),
		apiKey:    apiKey,
		modelName: modelName,
	}
}
