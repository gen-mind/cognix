package ai

import (
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SearcherConfig is a configuration struct for the searcher module.
type SearcherConfig struct {
	ApiVectorSearch  string `env:"API-VECTOR-SEARCH" envDefault:"INTERNAL"`
	InternalSearcher *EmbeddingConfig
	GRPCSearcher     *VectorSearchConfig
}

// SearcherResponse represents the response structure for a search operation.
type SearcherResponse struct {
	DocumentID int64  `json:"document_id,omitempty"`
	Content    string `json:"content,omitempty"`
}

// Searcher is an interface that defines the method for finding documents based on search criteria.
type Searcher interface {
	FindDocuments(ctx context.Context, userID, tenantID uuid.UUID,
		embeddingModel string,
		message string,
		collectionNames ...string) ([]*SearcherResponse, error)
}

// NewSearcher creates a new Searcher based on the specified searcherType.
//
// It takes in the searcherType as a string, embeddBuilder as an EmbeddingBuilder,
// vectorDB as a VectorDBClient, and embeddGRPCBuilder as a GRPCEmbeddingBuilder.
//
// It returns a Searcher interface and an error.
//
// It checks the value of searcherType and returns an instance of InternalSearcher if
// the searcherType is VectorSearchInternal, or an instance of SearcherGRPC if the
// searcherType is VectorSearchGRPCService. Otherwise, it returns an error indicating
// that the specified searcherType is not implemented.
//
// The InternalSearcher implementation of the Searcher interface uses the embeddBuilder
// and vectorDB to search for documents by performing embedding and loading from the vector
// database.
//
// The SearcherGRPC implementation of the Searcher interface uses the embeddGRPCBuilder to
// search for documents by performing vector search over gRPC.
func NewSearcher(
	searcherType string,
	embeddBuilder *EmbeddingBuilder,
	vectorDB storage.VectorDBClient,
	embeddGRPCBuilder *GRPCEmbeddingBuilder,
) (Searcher, error) {
	zap.S().Debugf("searcher type %s", searcherType)
	switch searcherType {
	case VectorSearchInternal:
		return &InternalSearcher{
			embeddBuilder: embeddBuilder,
			vectorDB:      vectorDB,
		}, nil
	case VectorSearchGRPCService:
		return &SearcherGRPC{
			embeddBuilder: embeddGRPCBuilder,
		}, nil
	}
	return nil, fmt.Errorf("vector searcher %s not implemented", searcherType)
}
