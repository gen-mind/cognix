package ai

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	StrategyStatic = "STATIC"
	StrategyLLM    = "LLM"
	MaxChunkSize   = 65536
)

type (
	ChunkingConfig struct {
		Strategy           string `env:"CHUNKING_STRATEGY" envDefault:"static"`
		StaticChunkSize    int    `env:"CHUNKING_STATIC_CHARS" envDefault:"16384"`
		StaticChunkOverlap int    `env:"CHUNKING_STATIC_CHARS_OVERLAP" envDefault:"20"`
	}

	Chunking interface {
		Split(ctx context.Context, text string) ([]string, error)
	}

	EmbeddingParser interface {
		Parse(ctx context.Context, payload *proto.EmbeddingRequest) (*proto.EmbeddingResponse, error)
	}
)

func (v ChunkingConfig) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Strategy, validation.Required, validation.In(StrategyStatic, StrategyLLM)),
		validation.Field(&v.StaticChunkSize, validation.Max(MaxChunkSize)),
		validation.Field(&v.StaticChunkOverlap, validation.Min(0),
			validation.Max(v.StaticChunkSize/3)),
	)
}