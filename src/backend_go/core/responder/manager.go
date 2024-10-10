package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"sync"
)

// Manager is a type that represents a manager responsible for managing chat responses.
// It contains a channel to receive responses, a wait group to manage goroutines,
// and a list of chat responders.
type Manager struct {
	ch         chan *Response
	wg         *sync.WaitGroup
	responders []ChatResponder
}

// Send is a method of the Manager struct that sends chat messages to all responders in parallel. It takes in a context, user,
// a boolean flag to indicate if the message should be processed by the language model, a parent message (if any), and a persona (if any).
// It loops through all the responders and calls their individual Send method in separate goroutines. It uses a WaitGroup to wait
// for all goroutines to finish before closing the channel m.ch.
func (m *Manager) Send(cx context.Context,
	user *model.User,
	noLLM bool,
	parentMessage *model.ChatMessage,
	persona *model.Persona) {
	for _, responder := range m.responders {
		m.wg.Add(1)
		go responder.Send(cx, m.ch, m.wg, user, noLLM, parentMessage, persona)
	}
	m.wg.Wait()
	close(m.ch)
}

// Receive receives a response from the Manager's channel.
//
// If there is a response available in the channel, it returns the response and true.
// If there are no more responses in the channel, it returns a default end response and false.
//
// The response is of type *Response, which contains information about the received message.
// The boolean value indicates if a response is available or not.
//
// Example usage:
//
//	response, ok := manager.Receive()
//	if ok {
//	    // Process the response
//	} else {
//	    // No more responses available
//	}
func (m *Manager) Receive() (*Response, bool) {
	for response := range m.ch {
		return response, true
	}
	return &Response{
		IsValid: true,
		Type:    ResponseEnd,
	}, false
}

// NewManager creates a new Manager object with the given ChatResponder(s).
//
// Parameters:
//   - responders: The variable number of ChatResponder objects to be used by the Manager.
//
// Returns:
//   - *Manager: The newly created Manager object.
func NewManager(responders ...ChatResponder) *Manager {
	return &Manager{
		ch:         make(chan *Response, 1),
		wg:         &sync.WaitGroup{},
		responders: append([]ChatResponder{}, responders...),
	}
}
