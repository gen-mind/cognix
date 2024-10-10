package ai

import (
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

var EmbeddingModule = fx.Options(
	fx.Provide(func() (*SearcherConfig, error) {
		cfg := SearcherConfig{
			InternalSearcher: &EmbeddingConfig{},
			GRPCSearcher:     &VectorSearchConfig{},
		}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newInternal,
		newGrpcEmbedder,
		newSearcher),
)

func newInternal(cfg *SearcherConfig) *EmbeddingBuilder {
	return NewEmbeddingBuilder(cfg.InternalSearcher)
}
func newGrpcEmbedder(cfg *SearcherConfig) *GRPCEmbeddingBuilder {
	return NewGRPCEmbeddingBuilder(cfg.GRPCSearcher)
}

func newSearcher(cfg *SearcherConfig,
	internalBuilder *EmbeddingBuilder,
	vectorDBClinet storage.VectorDBClient,
	grpcBuilder *GRPCEmbeddingBuilder) (Searcher, error) {
	return NewSearcher(cfg.ApiVectorSearch, internalBuilder, vectorDBClinet, grpcBuilder)
}
