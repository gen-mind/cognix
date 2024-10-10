package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

// Youtube is a struct representing a Youtube connector, which is a type of connector that inherits from the Base struct.
//
// The struct contains the following fields:
// - Base: the base properties and methods needed for various connectors.
// - param: a pointer to a Youtube Parameters struct that contains URL parameter.
//
// The Youtube struct does not have any additional methods or functionalities.
type (
	Youtube struct {
		Base
		param *YoutubeParameters
	}
	YoutubeParameters struct {
		URL string `url:"url"`
	}
)

// Validate validates the YouTube parameters.
func (p YoutubeParameters) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.URL, validation.Required,
			is.URL),
	)
}

// Validate validates the Youtube connector by checking if the file parameter is nil.
// If the file parameter is nil, it returns an error indicating that the file parameter is required.
// Otherwise, it calls the Validate method of the file parameter and returns its error (if any).
func (c *Youtube) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

// PrepareTask prepares a task by adding a document to the model if there are no existing documents.
// If the model's status is either ConnectorStatusError or ConnectorStatusSuccess, it returns nil.
// Otherwise, it calls the RunSemantic method of the task with the appropriate data.
func (c *Youtube) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {
	if len(c.model.Docs) == 0 {
		doc, ok := c.model.DocsMap[c.param.URL]
		if !ok {
			doc = &model.Document{
				SourceID:        c.param.URL,
				ConnectorID:     c.Base.model.ID,
				URL:             c.param.URL,
				Signature:       "",
				ChunkingSession: uuid.NullUUID{sessionID, true},
				OriginalURL:     c.param.URL,
			}
			c.model.Docs = append(c.model.Docs, doc)
		}
	}
	// ignore  file that was analyzed
	if c.model.Status == model.ConnectorStatusError || c.model.Status == model.ConnectorStatusSuccess {
		return nil
	}
	return task.RunSemantic(ctx, &proto.SemanticData{
		Url:            c.param.URL,
		ConnectorId:    c.model.ID.IntPart(),
		FileType:       proto.FileType_YT,
		CollectionName: c.model.CollectionName(),
		ModelName:      c.model.User.EmbeddingModel.ModelID,
		ModelDimension: int32(c.model.User.EmbeddingModel.ModelDim),
	})
}

// Execute executes the YouTube connector task.
// This method takes a context and a map of parameters as input and returns a channel of
// *Response objects as output. It runs the task asynchronously by launching a goroutine.
// Once the goroutine starts, it closes the resultCh channel and returns it.
func (c *Youtube) Execute(ctx context.Context, param map[string]string) chan *Response {
	go func() {
		close(c.resultCh)
	}()
	return c.resultCh
}

// NewYoutube creates a new instance of the Youtube connector by configuring the connector and validating its parameters.
// It initializes the Youtube struct and sets the connector configuration using the provided Connector instance.
// If the connector-specific configuration cannot be converted into a YoutubeParameters instance or validation fails,
// an error is returned. Otherwise, the Youtube connector instance is returned with no error.
func NewYoutube(connector *model.Connector) (Connector, error) {
	youtube := Youtube{}
	youtube.Base.Config(connector)
	youtube.param = &YoutubeParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(youtube.param); err != nil {
		return nil, err
	}
	if err := youtube.Validate(); err != nil {
		return nil, err
	}
	return &youtube, nil
}
