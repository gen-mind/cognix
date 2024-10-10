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

type (
	// Web struct represents a web connector with base properties and WebParameters.
	Web struct {
		Base
		param *WebParameters
	}

	// WebParameters struct represents the parameters for a Web connector.
	//
	// The struct contains the following fields:
	// - URL: a string representing the URL. (url:"url")
	// - SiteMap: a string representing the site map. (json:"site_map")
	// - SearchForSitemap: a boolean indicating whether to search for a sitemap. (json:"search_for_sitemap")
	// - URLRecursive: a boolean indicating whether the URL should be fetched recursively. (json:"url_recursive")
	WebParameters struct {
		URL              string `url:"url"`
		SiteMap          string `json:"site_map"`
		SearchForSitemap bool   `json:"search_for_sitemap"`
		URLRecursive     bool   `json:"url_recursive"`
	}
)

// Validate method validates the WebParameters struct by checking if the URL field is required and a valid URL using the is.URL validator.
// Returns an error if validation fails.
func (p WebParameters) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.URL, validation.Required,
			is.URL),
	)
}

// Validate method validates the Web struct by checking if the param field is nil.
// If the param field is nil, it returns an error "file parameter is required".
// Otherwise, it calls the Validate method of the param field, which is a WebParameters struct,
// and returns the result of the validation.
func (c *Web) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

// PrepareTask prepares a task for execution by performing necessary setup and validation.
// If the Web connector is new, a connectorTask for preparing the document table is run.
// The task is then executed by calling the RunSemantic method of the task object, passing the necessary data.
//
// Parameters:
//   - ctx: the context.Context object for managing the task execution
//   - sessionID: the UUID representing the session ID
//   - task: the Task object that represents the task to be executed
//
// Returns:
//   - error: an error object if any error occurs during the preparation or execution of the task
//     or nil if the preparation and execution are successful.
func (c *Web) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {

	// if this connector new we need to run connectorTask for prepare document table
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
	var rootDoc *model.Document
	for _, doc := range c.model.Docs {
		if !doc.ParentID.Valid {
			rootDoc = doc
			rootDoc.ChunkingSession = uuid.NullUUID{sessionID, true}
			break
		}
	}
	c.model.Docs = append([]*model.Document{}, rootDoc)

	if rootDoc == nil {
		return fmt.Errorf("root document not found")
	}

	return task.RunSemantic(ctx, &proto.SemanticData{
		Url:              c.param.URL,
		SiteMap:          c.param.SiteMap,
		SearchForSitemap: c.param.SearchForSitemap,
		UrlRecursive:     c.param.URLRecursive,
		DocumentId:       rootDoc.ID.IntPart(),
		ConnectorId:      c.model.ID.IntPart(),
		FileType:         proto.FileType_URL,
		CollectionName:   c.model.CollectionName(),
		ModelName:        c.model.User.EmbeddingModel.ModelID,
		ModelDimension:   int32(c.model.User.EmbeddingModel.ModelDim),
	})
}

// Execute runs the web connector to perform a web scraping operation.
//
// The function takes in a context and a param map and returns a channel of Response objects.
// The function launches a goroutine to execute the operation asynchronously.
// The function fetches the document from the model's DocsMap using the URL from the param map.
// If the document is not found, a new Document object is created and added to the DocsMap.
// The function then sends a Response object to the resultCh channel with the relevant information.
// Finally, the function closes the resultCh channel and returns it.
// The Response object contains the following fields:
// - URL: a string representing the URL.
// - SourceID: a string representing the source ID.
// - SiteMap: a string representing the site map.
// - SearchForSitemap: a boolean indicating whether to search for a sitemap.
// - DocumentID: an int64 representing the document ID.
// - MimeType: a string representing the mime type.
// - FileType: a proto.FileType representing the file type (16 for URL).
// Note: The function does not provide error handling and does not validate the param map.
func (c *Web) Execute(ctx context.Context, param map[string]string) chan *Response {
	go func() {
		doc, ok := c.model.DocsMap[c.param.URL]
		if !ok {
			doc = &model.Document{
				SourceID:    c.param.URL,
				ConnectorID: c.Base.model.ID,
				URL:         c.param.URL,
				Signature:   "",
			}
			c.Base.model.DocsMap[c.param.URL] = doc
		}
		c.resultCh <- &Response{
			URL:              c.param.URL,
			SourceID:         c.param.URL,
			SiteMap:          c.param.SiteMap,
			SearchForSitemap: c.param.SearchForSitemap,
			DocumentID:       doc.ID.IntPart(),
			MimeType:         model.MIMEURL,
			FileType:         proto.FileType_URL,
		}
		close(c.resultCh)
	}()
	return c.resultCh
}

// NewWeb is a function that initializes and configures a new instance of Web connector.
// It takes a pointer to a model.Connector struct as input and returns a Connector interface and an error.
//
// The function follows the following steps:
// 1. Create a new empty instance of Web.
// 2. Call the Config function of web.Base to set the connector model.
// 3. Initialize the param field of web with a new instance of WebParameters.
// 4. Call the ToStruct function of connector.ConnectorSpecificConfig to populate web.param.
// 5. If there is an error during the conversion, return nil and the error.
// 6. Call the Validate function of web to validate the parameters.
// 7. If there is an error during the validation, return nil and the error.
// 8. Return a pointer to web as a Connector interface and nil as the error.
//
// The returned Connector interface implements the Execute, PrepareTask, and Validate methods.
//
// The function uses the following types:
// - model.Connector: a struct that represents a table connector.
// - Connector: an interface that represents a connector.
// - Web: a struct that represents a Web connector and contains the base properties and methods needed.
// - Base: a struct that represents the base properties and methods needed for various connectors.
// - WebParameters: a struct that represents the parameters for a Web connector.
// - connector.ConnectorSpecificConfig: a map that stores the connector-specific configuration.
//
// The function uses the following methods:
// - web.Base.Config: a method that sets the connector model and initializes the result channel.
// - connector.ConnectorSpecificConfig.ToStruct: a method that converts the connector-specific configuration to a struct.
// - web.Validate: a method that validates the WebParameters.
//
// The function uses the following package:
// - "github.com/clinkerhq/config-server/repository": a repository package that provides the interface for interacting with the connector repository.
//
// The function returns the following values
// - Connector: an interface that represents a connector.
// - error: an error, if any, encountered during the execution of the function.
func NewWeb(connector *model.Connector) (Connector, error) {
	web := Web{}
	web.Base.Config(connector)
	web.param = &WebParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(web.param); err != nil {
		return nil, err
	}
	if err := web.Validate(); err != nil {
		return nil, err
	}
	return &web, nil
}
