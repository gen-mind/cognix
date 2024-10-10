package logic

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

type (

	// DocumentBL is an interface that defines the methods for managing documents.
	// UploadDocument uploads a document to the storage and creates a corresponding entry in the repository.
	// Parameters:
	// - ctx: The context.Context for the request.
	// - user: The user who is uploading the document.
	// - fileName: The name of the document file.
	// - contentType: The content type of the document file.
	// - file: The io.Reader containing the document file to be uploaded.
	// Returns:
	// - *model.Document: The newly created Document object.
	// - error: Any error that occurred during the upload process.
	DocumentBL interface {
		UploadDocument(ctx context.Context, user *model.User, fileName, contentType string, file io.Reader) (*model.Document, error)
	}

	// documentBL is a struct that represents a business logic layer for managing documents.
	// It contains dependencies for document repository, MinIO client, and connector repository.
	// The struct implements the DocumentBL interface.
	//
	// Fields:
	// - documentRepo: The repository for managing documents.
	// - minioClient: The MinIO client for uploading documents to storage.
	// - connectorRepo: The repository for managing connectors.
	//
	// Example usage:
	// db := repository.NewDocumentRepository()
	// minio := storage.NewMinIOClient()
	// connector := repository.NewConnectorRepository()
	// bl := documentBL{
	//   documentRepo:  db,
	//   minioClient:   minio,
	//   connectorRepo: connector,
	// }
	documentBL struct {
		documentRepo  repository.DocumentRepository
		minioClient   storage.FileStorageClient
		connectorRepo repository.ConnectorRepository
	}
)

// UploadDocument uploads a document to storage, creates a document record in the repository,
// and returns the created document.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - user: The user who is uploading the document.
// - fileName: The name of the file being uploaded.
// - contentType: The content type of the file being uploaded.
// - file: The file reader providing the content of the file to be uploaded.
//
// Returns:
// - *model.Document: The created document.
// - error: The error encountered during the upload process or document creation.
func (b *documentBL) UploadDocument(ctx context.Context, user *model.User, fileName, contentType string, file io.Reader) (*model.Document, error) {

	fileURL, signature, err := b.minioClient.Upload(ctx, model.BucketName(user.TenantID),
		fmt.Sprintf("user-%s/%s-%s", user.ID.String(), uuid.New().String(), fileName), contentType, file)
	if err != nil {
		return nil, err
	}
	connector, err := b.connectorRepo.GetBySource(ctx, user.TenantID, user.ID, model.SourceTypeFile)
	if err != nil {
		return nil, err
	}
	document := &model.Document{
		SourceID:     fileURL,
		ConnectorID:  connector.ID,
		URL:          fileURL,
		Signature:    signature,
		CreationDate: time.Now().UTC(),
	}
	if err = b.documentRepo.Create(ctx, document); err != nil {
		return nil, err
	}
	return document, nil
}

// NewDocumentBL is a function that creates a new instance of the DocumentBL interface.
// The function takes the following parameters:
// - documentRepo: The repository for managing documents.
// - connectorRepo: The repository for managing connectors.
// - minioClient: The MinIO client for uploading documents to storage.
// The function returns a new instance of the DocumentBL interface.
// Example usage:
// documentRepo := repository.NewDocumentRepository()
// connectorRepo := repository.NewConnectorRepository()
// minioClient := storage.NewMinIOClient()
// documentBL := NewDocumentBL(documentRepo, connectorRepo, minioClient)
func NewDocumentBL(documentRepo repository.DocumentRepository,
	connectorRepo repository.ConnectorRepository,
	minioClient storage.FileStorageClient) DocumentBL {
	return &documentBL{documentRepo: documentRepo,
		connectorRepo: connectorRepo,
		minioClient:   minioClient,
	}
}
