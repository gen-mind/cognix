package logic

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

// Config is a struct that holds the configuration settings.
//
// The Config struct contains the following fields:
//   - RedirectURL:                A string representing the redirect URL.
//     It is tagged with `env:"REDIRECT_URL"`.
//   - DefaultEmbeddingModel:      A string representing the default embedding model.
//     It is tagged with `env:"DEFAULT_EMBEDDING_MODEL"` and has a default value of "paraphrase-multilingual-mpnet-base-v2".
//   - DefaultEmbeddingVectorSize: An integer representing the default embedding vector size.
//     It is tagged with `env:"DEFAULT_EMBEDDING_VECTOR_SIZE"` and has a default value of 768.
type Config struct {
	RedirectURL                string `env:"REDIRECT_URL"`
	DefaultEmbeddingModel      string `env:"DEFAULT_EMBEDDING_MODEL" envDefault:"paraphrase-multilingual-mpnet-base-v2"`
	DefaultEmbeddingVectorSize int    `env:"DEFAULT_EMBEDDING_VECTOR_SIZE" envDefault:"768"`
}

var BLLModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	}),
	fx.Provide(
		NewConnectorBL,
		NewAuthBL,
		NewPersonaBL,
		NewChatBL,
		NewDocumentBL,
		NewEmbeddingModelBL,
		NewTenantBL,
	),
)
