package ai

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
)

type (

	// Response is a struct that represents the response from the OpenAI chat API.
	Response struct {
		Message string
	}

	// Client is an interface for making requests to the OpenAI chat API.
	// The Request method takes a context and a message as input and returns a Response or an error.
	Client interface {
		Request(ctx context.Context, message string) (*Response, error)
	}

	// openAIClient is a struct that represents the client for making requests to the OpenAI chat API.
	// It consists of a client of type *openai.Client and a modelID of type string.
	openAIClient struct {
		client  *openai.Client
		modelID string
	}
)

// Request is a method of the openAIClient struct that makes a request to the OpenAI chat API.
// It takes a context.Context parameter and a string message parameter.
// It returns a *Response and an error.
//
// The method first creates a ChatCompletionMessage using the user's message.
// Then it calls the client's CreateChatCompletion method to make the API request.
// If there is an error, it returns nil and the error.
// If the API request is successful, it creates a Response with the content of the first message choice
// and returns it along with nil for the error.
func (o *openAIClient) Request(ctx context.Context, message string) (*Response, error) {

	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	}
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    o.modelID,
			Messages: []openai.ChatCompletionMessage{userMessage},
		},
	)
	if err != nil {
		return nil, err
	}
	response := &Response{Message: resp.Choices[0].Message.Content}
	return response, nil
}

// NewOpenAIClient is a function that creates a new instance of the Client.
// It takes the modelID and apiKey as input parameters and returns an instance of Client.
// The function creates a new openaIClient struct with the provided modelID and apiKey.
// It then initializes the client field with the openai.NewClient function using the apiKey.
// Finally, it sets the modelID field with the provided modelID and returns the created struct as an Client.
func NewOpenAIClient(modelID, apiKey string) Client {

	return &openAIClient{
		client:  openai.NewClient(apiKey),
		modelID: modelID,
	}
}
