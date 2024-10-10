package main

import (
	"bytes"
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"strings"

	"io"
	"time"
)

// Executor represents a type that executes tasks related to connectors.
type Executor struct {
	cfg            *Config
	connectorRepo  repository.ConnectorRepository
	docRepo        repository.DocumentRepository
	msgClient      messaging.Client
	minioClient    storage.FileStorageClient
	milvusClient   storage.VectorDBClient
	oauthClient    *resty.Client
	downloadClient *resty.Client
}

// run is a method that listens to a specific stream and topic using the messaging.Client provided
// and performs a specific task, which is passed as a MessageHandler function.
// It returns an error if there is any issue with listening to the stream and topic.
func (e *Executor) run(streamName, topic string, task messaging.MessageHandler) {
	if err := e.msgClient.Listen(context.Background(), streamName, topic, task); err != nil {
		zap.S().Errorf("failed to listen[%s]: %v", topic, err)
	}
	return
}

// runConnector is a method that handles the execution of a connector task based on a connector request message.
// It unmarshals the message data to obtain the trigger information, retrieves the connector model from the repository,
// creates a new instance of a connector, executes the connector task, handles the results by saving content and updating
// documents, and publishes messages to either the Voice or Semantic stream based on the file type.
// After processing all the results, it deletes unused files associated with the connector and updates the connector status.
// Finally, it updates the connector in the repository and returns any error that occurred during the process.
func (e *Executor) runConnector(ctx context.Context, msg jetstream.Msg) error {
	startTime := time.Now()
	//ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Header()))
	var trigger proto.ConnectorRequest

	if err := proto2.Unmarshal(msg.Data(), &trigger); err != nil {
		zap.S().Errorf("Error unmarshalling message: %s", err.Error())
		return err
	}
	// read connector model with documents, embedding model
	connectorModel, err := e.connectorRepo.GetByID(ctx, trigger.GetId())
	if err != nil {
		return err
	}
	defer func() {
		zap.S().Infof("connector %s completed. elapsed time: %d ms", connectorModel.Name, time.Since(startTime)/time.Millisecond)
	}()

	zap.S().Infof("receive message : %s [%d]", connectorModel.Name, connectorModel.ID.IntPart())
	// refresh token if needed
	connectorModel.Status = model.ConnectorStatusWorking

	// create new instance of connector by connector model
	connectorWF, err := connector.New(connectorModel, e.connectorRepo, e.cfg.OAuthURL)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	if trigger.Params == nil {
		trigger.Params = make(map[string]string)
	}
	// execute connector
	resultCh := connectorWF.Execute(ctx, trigger.Params)
	// read result from channel
	hasSemanticMessage := false
	for result := range resultCh {
		var loopErr error
		// empty result when channel was closed.
		if result.SourceID == "" {
			break
		}
		hasSemanticMessage = true

		// save content in minio
		if result.Content != nil {
			if err = e.saveContent(ctx, result); err != nil {
				loopErr = err
			}

		}
		// find or create document from result
		doc := e.handleResult(connectorModel, result)
		// create or update document in database
		if doc.ID.IntPart() != 0 {
			loopErr = e.docRepo.Update(ctx, doc)
		} else {
			loopErr = e.docRepo.Create(ctx, doc)
		}

		if loopErr != nil {
			err = loopErr
			zap.S().Errorf("Failed to update document: %v", loopErr)
			continue
		}

		// send message to chunking service

		if _, ok := model.VoiceFileTypes[result.FileType]; ok {
			// send message to Voice
			voiceDate := proto.VoiceData{
				Url:            result.URL,
				DocumentId:     doc.ID.IntPart(),
				ConnectorId:    connectorModel.ID.IntPart(),
				FileType:       result.FileType,
				CollectionName: connectorModel.CollectionName(),
				ModelName:      connectorModel.User.EmbeddingModel.ModelID,
				ModelDimension: int32(connectorModel.User.EmbeddingModel.ModelDim),
			}
			zap.S().Infof("send message to voice service %s - %s", connectorModel.Name, result.URL)
			if loopErr = e.msgClient.Publish(ctx,
				e.msgClient.StreamConfig().VoiceStreamName,
				e.msgClient.StreamConfig().VoiceStreamSubject,
				&voiceDate); loopErr != nil {
				err = loopErr
				zap.S().Errorf("Failed to publish voice service: %v", loopErr)
				continue
			}

		} else {
			// send message to semantic
			semanticData := proto.SemanticData{
				Url:            result.URL,
				DocumentId:     doc.ID.IntPart(),
				ConnectorId:    connectorModel.ID.IntPart(),
				FileType:       result.FileType,
				CollectionName: connectorModel.CollectionName(),
				ModelName:      connectorModel.User.EmbeddingModel.ModelID,
				ModelDimension: int32(connectorModel.User.EmbeddingModel.ModelDim),
			}
			zap.S().Infof("send message to semantic %s - %s", connectorModel.Name, result.URL)
			if loopErr = e.msgClient.Publish(ctx,
				e.msgClient.StreamConfig().SemanticStreamName,
				e.msgClient.StreamConfig().SemanticStreamSubject,
				&semanticData); loopErr != nil {
				err = loopErr
				zap.S().Errorf("Failed to publish semantic: %v", loopErr)
				continue
			}
		}
	}
	if errr := e.deleteUnusedFiles(ctx, connectorModel); err != nil {
		zap.S().Errorf("deleting unused files: %v", errr)
		if err == nil {
			err = errr
		}
	}
	if err != nil {
		zap.S().Errorf("failed to update documents: %v", err)
		connectorModel.Status = model.ConnectorStatusUnableProcess
	} else {
		if !hasSemanticMessage {
			connectorModel.Status = model.ConnectorStatusSuccess
		}
	}
	connectorModel.LastUpdate = pg.NullTime{time.Now().UTC()}

	if err = e.connectorRepo.Update(ctx, connectorModel); err != nil {
		return err
	}
	return nil
}

// deleteUnusedFiles is a method that deletes unused files associated with a connector.
// It iterates through the documents in the connector's DocsMap, checks if the document
// is not marked as exists and has a non-zero ID. If the document's URL starts with "minio:",
// it uses the e.minioClient to delete the corresponding object from the MinIO storage.
// It also collects the non-zero document IDs and stores them in the "ids" slice.
// After iterating through all the documents, if the "ids" slice is not empty,
// it uses the e.milvusClient to delete the documents from the Milvus storage and
// the e.docRepo to delete the documents from the repository.
// If there are any errors during the process, it returns the error. Otherwise, it returns nil.
func (e *Executor) deleteUnusedFiles(ctx context.Context, connector *model.Connector) error {
	var ids []int64
	for _, doc := range connector.DocsMap {
		if doc.IsExists || doc.ID.IntPart() == 0 {
			continue
		}
		filepath := strings.Split(doc.URL, ":")
		if len(filepath) == 3 && filepath[0] == "minio" {
			if err := e.minioClient.DeleteObject(ctx, filepath[1], filepath[2]); err != nil {
				return err
			}
		}
		ids = append(ids, doc.ID.IntPart())
	}
	if len(ids) > 0 {
		if err := e.milvusClient.Delete(ctx, connector.CollectionName(), ids...); err != nil {
			return err
		}
		return e.docRepo.DeleteByIDS(ctx, ids...)
	}
	return nil
}

// saveContent is a method that saves the content of a response to a storage system. If the
// response contains a URL, it will download the file and save it. Otherwise, it will create a
// reader from the raw content and save it. The method uses the `minioClient` to upload the file
// to the specified bucket using the provided name, MIME type, and reader. Upon successful upload,
// it sets the URL of the response to the corresponding MinIO URL and returns nil. If any error
// occurs during the process, it returns the error.
//
// The method expects a context and a pointer to a `connector.Response` as input parameters.
func (e *Executor) saveContent(ctx context.Context, response *connector.Response) error {

	var reader io.Reader
	//  download file if url presented.
	if response.Content.URL != "" {
		fileResponse, err := e.downloadClient.R().
			SetDoNotParseResponse(true).
			Get(response.Content.URL)
		defer fileResponse.RawBody().Close()
		if err = utils.WrapRestyError(fileResponse, err); err != nil {
			return err
		}
		reader = fileResponse.RawBody()
	} else {
		if response.Content.Reader != nil {
			reader = response.Content.Reader
			defer response.Content.Reader.Close()
		} else {
			// create reader from raw content
			reader = bytes.NewReader(response.Content.Body)
		}
	}

	fileName, _, err := e.minioClient.Upload(ctx, response.Content.Bucket, response.Name, response.MimeType, reader)
	if err != nil {
		return err
	}
	response.URL = fmt.Sprintf("minio:%s:%s", response.Content.Bucket, fileName)
	return nil
}

// handleResult is a method that handles the result of a connector task.
// It takes a pointer to a Connector model and a pointer to a Response model as input.
// It checks if the specified SourceID exists in the Connector's DocsMap.
// If the SourceID does not exist, it creates a new Document model with the necessary fields
// and adds it to the Connector's DocsMap. The Document's URL is set to the Response's URL,
// the CreationDate is set to the current UTC time, and other fields are populated from the Response.
// If the SourceID already exists, it updates the existing Document's URL to the Response's URL
// and updates the LastUpdate field to the current UTC time.
// The method returns the updated or newly created Document.
func (e *Executor) handleResult(connectorModel *model.Connector, result *connector.Response) *model.Document {
	doc, ok := connectorModel.DocsMap[result.SourceID]
	if !ok {
		doc = &model.Document{
			SourceID:     result.SourceID,
			ConnectorID:  connectorModel.ID,
			URL:          result.URL,
			Signature:    result.Signature,
			CreationDate: time.Now().UTC(),
		}
		connectorModel.DocsMap[result.SourceID] = doc
	} else {
		doc.URL = result.URL
		doc.LastUpdate = pg.NullTime{time.Now().UTC()}
	}

	return doc
}

// NewExecutor is a constructor function that creates a new instance of the Executor struct.
// It takes in various dependencies including a *Config for configuration, a ConnectorRepository for accessing connectors,
// a DocumentRepository for accessing documents, a messaging.Client for handling messaging,
// a storage.FileStorageClient for working with MinIO storage, and a storage.VectorDBClient for working with Milvus storage.
// It returns a pointer to the newly created Executor instance.
func NewExecutor(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	streamClient messaging.Client,
	minioClient storage.FileStorageClient,
	milvusClient storage.VectorDBClient,
) *Executor {
	return &Executor{
		cfg:           cfg,
		connectorRepo: connectorRepo,
		docRepo:       docRepo,
		msgClient:     streamClient,
		minioClient:   minioClient,
		milvusClient:  milvusClient,
		oauthClient: resty.New().
			SetTimeout(time.Minute).
			SetBaseURL(cfg.OAuthURL),
		downloadClient: resty.New().
			SetTimeout(time.Minute).
			SetDoNotParseResponse(true),
	}
}
