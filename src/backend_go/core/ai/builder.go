package ai

import (
	"cognix.ch/api/v2/core/model"
	"sync"
)

// Builder is a type that manages the creation and caching of Client instances.
// Type Declaration:
//
//	type Builder struct {
//			clients map[int64]Client
//			mx      sync.Mutex
//	}
//
// Usage Example:
// NewBuilder returns a new instance of Builder.
// go doc Builder
//
// Returns:
//
//	An instance of Builder.
//
// Example:
//
//	builder := NewBuilder()
type Builder struct {
	clients map[int64]Client
	mx      sync.Mutex
}

// NewBuilder returns a new instance of Builder.
//
// Returns:
//
//	An instance of Builder.
func NewBuilder() *Builder {
	return &Builder{clients: make(map[int64]Client)}
}

// New returns a new instance of Client. If the client for the given LLM ID already exists in the Builder's cache,
// that client is returned; otherwise, a new client is created using the NewOpenAIClient function. The client is then
// added to the Builder's cache and returned.
//
// Parameters:
//
//	llm - The LLM model used to create the Client.
//
// Returns:
//
//	An instance of Client.
func (b *Builder) New(llm *model.LLM) Client {
	b.mx.Lock()
	defer b.mx.Unlock()
	if client, ok := b.clients[llm.ID.IntPart()]; ok {
		return client
	}
	client := NewChatAI(llm.Endpoint, llm.ApiKey, llm.ModelID)

	b.clients[llm.ID.IntPart()] = client
	return client
}

// Invalidate removes the Client instance for the given LLM ID from the Builder's cache.
func (b *Builder) Invalidate(llm *model.LLM) {
	b.mx.Lock()
	delete(b.clients, llm.ID.IntPart())
	b.mx.Unlock()
}
