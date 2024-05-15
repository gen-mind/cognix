package responder

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

type embedding struct {
	embedding      proto.EmbeddServiceClient
	milvusClinet   storage.MilvusClient
	docRepo        repository.DocumentRepository
	embeddingModel string
}

func (r *embedding) Send(ctx context.Context, ch chan *Response, wg *sync.WaitGroup, user *model.User, parentMessage *model.ChatMessage) {

	for i := 0; i < 4; i++ {
		ch <- &Response{
			IsValid: true,
			Type:    ResponseDocument,
			Message: nil,
			Document: &model.DocumentResponse{
				ID:          decimal.NewFromInt(int64(i)),
				DocumentID:  "11",
				Link:        fmt.Sprintf("link for document %d", i),
				Content:     fmt.Sprintf("content of document %d", i),
				UpdatedDate: time.Now().UTC().Add(-48 * time.Hour),
				MessageID:   parentMessage.ID,
			},
		}
	}
	wg.Done()
}

func (r *embedding) FindDocuments(ctx context.Context,
	ch chan *Response,
	message *model.ChatMessage,
	collectionNames ...string) ([]*model.DocumentResponse, error) {
	response, err := r.embedding.GetEmbedd(ctx, &proto.EmbeddRequest{
		Content: message.ParentMessage.Message,
		Model:   r.embeddingModel,
	})
	if err != nil {
		ch <- &Response{
			IsValid: false,
			Type:    ResponseError,
			Err:     err,
		}
		return nil, err
	}
	var result []*model.DocumentResponse
	for _, collectionName := range collectionNames {
		docs, err := r.milvusClinet.Load(ctx, collectionName, response.GetVector())
		if err != nil {
			ch <- &Response{
				IsValid: false,
				Type:    ResponseError,
				Err:     err,
			}
			continue
		}
		for _, doc := range docs {
			resDocument := &model.DocumentResponse{
				ID:        decimal.NewFromInt(doc.DocumentID),
				MessageID: message.ID,
				Content:   doc.Content,
			}
			if dbDoc, err := r.docRepo.FindByID(ctx, doc.DocumentID); err == nil {
				resDocument.Link = dbDoc.Link
				resDocument.DocumentID = dbDoc.DocumentID
			}
			result = append(result, resDocument)
			ch <- &Response{
				IsValid:  true,
				Type:     ResponseDocument,
				Document: resDocument,
			}
		}
	}
	return result, nil
}

func NewEmbeddingResponder(embeddProto proto.EmbeddServiceClient,
	milvusClinet storage.MilvusClient,
	docRepo repository.DocumentRepository,
	embeddingModel string) *embedding {
	return &embedding{
		embedding:      embeddProto,
		milvusClinet:   milvusClinet,
		embeddingModel: embeddingModel,
		docRepo:        docRepo,
	}
}