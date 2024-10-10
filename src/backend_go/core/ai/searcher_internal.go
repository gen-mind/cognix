package ai

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/storage"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InternalSearcher struct {
	embeddBuilder *EmbeddingBuilder
	vectorDB      storage.VectorDBClient
}

func (i *InternalSearcher) FindDocuments(ctx context.Context, userID, tenantID uuid.UUID,
	embeddingModel string,
	message string, collectionNames ...string) ([]*SearcherResponse, error) {
	embedding, err := i.embeddBuilder.Client()
	if err != nil {
		return nil, err
	}
	response, err := embedding.GetEmbedding(ctx, &proto.EmbedRequest{
		Contents: []string{message},
		Model:    embeddingModel,
	})
	if err != nil {
		zap.S().Errorf("embeding service %s ", err.Error())
		return nil, err
	}
	var result []*SearcherResponse

	for _, collectionName := range collectionNames {
		if response.GetEmbeddings() == nil || len(response.GetEmbeddings()) == 0 {
			continue
		}
		docs, err := i.vectorDB.Load(ctx, collectionName, response.GetEmbeddings()[0].GetVector())
		if err != nil {
			zap.S().Errorf("error loading document from vector database :%s", err.Error())
			continue
		}
		for _, doc := range docs {
			resDocument := &SearcherResponse{
				DocumentID: doc.DocumentID,
				Content:    doc.Content,
			}
			result = append(result, resDocument)
		}
	}
	return result, nil
}
