package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"sync"
)

// ResponseMessage is a constant variable representing the key for the message in a response.
const (
	ResponseMessage  = "message"
	ResponseError    = "error"
	ResponseDocument = "document"
	ResponseEnd      = "end"
)

// Response represents a response object containing information about a chat message response.
// It includes fields for validity, type, chat message, document response, and an error.
type Response struct {
	IsValid  bool
	Type     string
	Message  *model.ChatMessage
	Document *model.DocumentResponse
	Err      error
}

// ChatResponder is an interface that represents an object capable of sending chat responses.
// It defines a method `Send` that takes in a context, a response channel, a wait group, a user,
// a boolean flag, a parent message, and a persona, and sends the chat response.
type ChatResponder interface {
	Send(cx context.Context, ch chan *Response, wg *sync.WaitGroup, user *model.User, noLLM bool, parentMessage *model.ChatMessage, persona *model.Persona)
}
