package ai

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type VectorSearchConfig struct {
	VectorSearchHost string `env:"SEARCH_GRPC_HOST,required"`
	VectorSearchPort int    `env:"SEARCH_GRPC_PORT,required"`
}

// GRPCEmbeddingBuilder is a struct that represents a builder for gRPC embedding functionality.
// It has a configuration, a client, and a mutex for concurrency safety.
type GRPCEmbeddingBuilder struct {
	cfg    *VectorSearchConfig
	client proto.SearchServiceClient
	mx     sync.Mutex
}

// NewGRPCEmbeddingBuilder returns a new instance of GRPCEmbeddingBuilder with the given VectorSearchConfig
// Parameters:
// - cfg: A pointer to a VectorSearchConfig object
// Returns:
// - A pointer to a GRPCEmbeddingBuilder object
// Example Usage:
//
//	cfg := &VectorSearchConfig{
//	    VectorSearchHost: "localhost",
//	    VectorSearchPort: 8080,
//	}
//
// builder := NewGRPCEmbeddingBuilder(cfg)
func NewGRPCEmbeddingBuilder(cfg *VectorSearchConfig) *GRPCEmbeddingBuilder {
	return &GRPCEmbeddingBuilder{
		cfg: cfg,
		mx:  sync.Mutex{},
	}
}

// Client returns a SearchServiceClient instance for making search requests.
// If the client is not already initialized, it creates a new connection
// to the VectorSearch service using the provided configuration and sets the client.
// It returns the client and an error if any occurred during initialization.
func (e *GRPCEmbeddingBuilder) Client() (proto.SearchServiceClient, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.client == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials())}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", e.cfg.VectorSearchHost, e.cfg.VectorSearchPort), dialOptions...)
		if err != nil {
			return nil, err
		}
		e.client = proto.NewSearchServiceClient(conn)
	}
	return e.client, nil
}
