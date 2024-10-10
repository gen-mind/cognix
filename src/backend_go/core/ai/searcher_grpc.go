package ai

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SearcherGRPC is a type that represents a searcher implementation using gRPC.
// It contains an embeddBuilder field of type *GRPCEmbeddingBuilder, which is
// responsible for constructing the gRPC embedding functionality.
type SearcherGRPC struct {
	embeddBuilder *GRPCEmbeddingBuilder
}

// FindDocuments searches for documents based on the provided criteria.
// It uses gRPC embedding functionality to perform the search.
//
// Parameters:
//   - ctx: The context.Context for the search operation.
//   - userID: The ID of the user.
//   - tenantID: The ID of the tenant.
//   - embeddingModel: The embedding model to use for the search.
//   - message: The message to search for.
//   - collectionNames: The optional collection names to search within.
//
// Returns:
//   - []*SearcherResponse: A slice of SearcherResponse objects representing the search results.
//   - error: An error if any occurred during the search.
func (i *SearcherGRPC) FindDocuments(ctx context.Context, userID, tenantID uuid.UUID,
	embeddingModel string,
	message string, collectionNames ...string) ([]*SearcherResponse, error) {
	embedding, err := i.embeddBuilder.Client()
	if err != nil {
		return nil, err
	}
	response, err := embedding.VectorSearch(ctx, &proto.SearchRequest{
		Content:         message,
		UserId:          userID.String(),
		TenantId:        tenantID.String(),
		CollectionNames: collectionNames,
		ModelName: embeddingModel,
	})
	if err != nil {
		zap.S().Errorf("embeding service %s ", err.Error())
		return nil, err
	}
	var result []*SearcherResponse

	for _, doc := range response.GetDocuments() {
		resDocument := &SearcherResponse{
			DocumentID: doc.GetDocumentId(),
			Content:    doc.GetContent(),
		}
		result = append(result, resDocument)
	}
	return result, nil
}
