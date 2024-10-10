package ai

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	_ "github.com/deluan/flowllm/llms/openai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

const (
	VectorSearchInternal    = "INTERNAL"
	VectorSearchGRPCService = "GRPC-SERVICE"
)

// EmbeddingConfig is a configuration struct for embedding module.
//
// It contains the configuration options for connecting to the embedding server
// over gRPC.
type EmbeddingConfig struct {
	EmbedderHost string `env:"EMBEDDER_GRPC_HOST,required"`
	EmbedderPort int    `env:"EMBEDDER_GRPC_PORT,required"`
}

type EmbeddingBuilder struct {
	cfg    *EmbeddingConfig
	client proto.EmbedServiceClient
	mx     sync.Mutex
}

func NewEmbeddingBuilder(cfg *EmbeddingConfig) *EmbeddingBuilder {
	return &EmbeddingBuilder{
		cfg: cfg,
		mx:  sync.Mutex{},
	}
}
func (e *EmbeddingBuilder) Client() (proto.EmbedServiceClient, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.client == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials())}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", e.cfg.EmbedderHost, e.cfg.EmbedderPort), dialOptions...)
		if err != nil {
			return nil, err
		}
		e.client = proto.NewEmbedServiceClient(conn)
	}
	return e.client, nil
}
