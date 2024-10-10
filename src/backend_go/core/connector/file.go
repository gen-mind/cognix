package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// File is a type that represents a file connector.
type (

	// File is a type that represents a file connector.
	File struct {
		Base
		param *FileParameters
		ctx   context.Context
	}

	// FileParameters is a struct that represents the parameters for a file connector.
	//
	// The struct contains the following fields:
	// - FileName: a string that specifies the name of the file.
	// - MIMEType: a string that specifies the MIME type of the file.
	FileParameters struct {
		FileName string `json:"file_name"`
		MIMEType string `json:"mime_type"`
	}
)

// Validate validates the FileParameters struct.
//
// It checks if the FileName and MIMEType fields are required and not empty.
// If any of the fields are empty, it returns an error.
// Otherwise, it returns nil.
func (p FileParameters) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.FileName, validation.Required),
		validation.Field(&p.MIMEType, validation.Required),
	)
}

// Validate validates the File struct.
//
// It checks if the file parameter is not nil.
// If the file parameter is nil, it returns an error.
// Otherwise, it calls the Validate method of the file parameter and returns its result.
func (c *File) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

// PrepareTask prepares the task before execution.
//
// It checks if the length of c.model.Docs is 0. If it is 0, it appends a new Document
// to c.model.Docs with the given link, connector ID, URL, creation date, and sets IsExists to true.
//
// If c.model.Status is either model.ConnectorStatusError or model.ConnectorStatusSuccess,
// it returns nil to ignore the file that was analyzed.
//
// Assigns sessionID to c.model.Docs[0].ChunkingSession.
// Calls task.RunSemantic method passing context and semantic data to run the task.
// Returns the result of task.RunSemantic method.
func (c *File) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {

	link := fmt.Sprintf("minio:tenant-%s:%s", c.model.User.EmbeddingModel.TenantID.String(), c.param.FileName)
	if len(c.model.Docs) == 0 {
		c.model.Docs = append(c.model.Docs, &model.Document{
			SourceID:     link,
			ConnectorID:  c.model.ID,
			URL:          link,
			CreationDate: time.Now().UTC(),
			IsExists:     true,
		})
	}
	// ignore  file that was analyzed
	if c.model.Status == model.ConnectorStatusError || c.model.Status == model.ConnectorStatusSuccess {
		return nil
	}
	c.model.Docs[0].ChunkingSession = uuid.NullUUID{sessionID, true}
	return task.RunSemantic(ctx, &proto.SemanticData{
		Url:            link,
		DocumentId:     c.model.Docs[0].ID.IntPart(),
		ConnectorId:    c.model.ID.IntPart(),
		FileType:       0,
		CollectionName: c.model.CollectionName(),
		ModelName:      c.model.User.EmbeddingModel.ModelID,
		ModelDimension: int32(c.model.User.EmbeddingModel.ModelDim),
	})
}

// Execute executes the File connector and returns a channel of Response pointers.
//
// It checks if the `param` map is nil or if the `FileName` field in the `param` map is empty.
// If either of these conditions is true, it returns an empty channel.
// Otherwise, it checks if the document with the given `FileName` already exists in the `DocsMap` of the File's base model.
// If it doesn't exist, a new Document struct is created with the given `FileName` and other necessary information,
// and it is added to the `DocsMap`.
// The `IsExists` field of the document is set to true.
// Then, it checks if the `MIMEType` field in the `param` map is a supported MIME type.
// If it is, a Response struct is created with the necessary fields and sent to the `resultCh` channel.
// If it is not a supported MIME type, an error message is logged.
// Finally, the `resultCh` channel is closed and returned.
func (c *File) Execute(ctx context.Context, param map[string]string) chan *Response {
	// do no used for this source
	c.ctx = ctx
	go func() {
		defer close(c.resultCh)
		if c.param == nil || c.param.FileName == "" {
			return
		}
		// check id document  already exists
		doc, ok := c.Base.model.DocsMap[c.param.FileName]
		url := fmt.Sprintf("minio:tenant-%s:%s", c.model.User.EmbeddingModel.TenantID, c.param.FileName)
		if !ok {
			doc = &model.Document{
				SourceID:    url,
				ConnectorID: c.model.ID,
				URL:         url,
				Signature:   "",
			}
			c.model.DocsMap[url] = doc
		}
		doc.IsExists = true
		if fileType, ok := model.SupportedMimeTypes[c.param.MIMEType]; ok {
			c.resultCh <- &Response{
				URL:      url,
				SourceID: url,
				FileType: fileType,
			}
		} else {
			zap.S().Errorf("Upsupported file type : %s ", c.param.MIMEType)
		}

	}()
	return c.resultCh
}

// NewFile creates new instance of file connector.
func NewFile(connector *model.Connector) (Connector, error) {
	fileConn := File{}
	fileConn.Base.Config(connector)
	fileConn.param = &FileParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(fileConn.param); err != nil {
		return nil, err
	}
	if err := fileConn.Validate(); err != nil {
		return nil, err
	}
	return &fileConn, nil
}
