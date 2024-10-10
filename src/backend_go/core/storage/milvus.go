package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// Package-level constants for column names and vector dimension.
const (
	ColumnNameID         = "id"
	ColumnNameDocumentID = "document_id"
	ColumnNameContent    = "content"
	ColumnNameVector     = "vector"

	VectorDimension = 1536

	IndexStrategyDISKANN   = "DISKANN"
	IndexStrategyAUTOINDEX = "AUTOINDEX"
	IndexStrategyNoIndex   = "NOINDEX"
)

// responseColumns is a variable of type []string that contains three column names: "id", "document_id", and "content".
var responseColumns = []string{ColumnNameID, ColumnNameDocumentID, ColumnNameContent}

// MilvusConfig represents the configuration for connecting to the Milvus service.
type (
	MilvusConfig struct {
		Address       string `env:"MILVUS_URL,required"`
		Username      string `env:"MILVUS_USERNAME,required"`
		Password      string `env:"MILVUS_PASSWORD,required"`
		MetricType    string `env:"MILVUS_METRIC_TYPE" envDefault:"COSINE"`
		IndexStrategy string `env:"MILVUS_INDEX_STRATEGY" envDefault:"DISKANN"`
	}
	MilvusPayload struct {
		ID         int64     `json:"id"`
		DocumentID int64     `json:"document_id"`
		Chunk      int64     `json:"chunk"`
		Content    string    `json:"content"`
		Vector     []float32 `json:"vector"`
	}
	VectorDBClient interface {
		CreateSchema(ctx context.Context, name string) error
		Save(ctx context.Context, collection string, payloads ...*MilvusPayload) error
		Load(ctx context.Context, collection string, vector []float32) ([]*MilvusPayload, error)
		Delete(ctx context.Context, collection string, documentID ...int64) error
	}
	milvusClient struct {
		client     milvus.Client
		cfg        *MilvusConfig
		MetricType entity.MetricType
	}
)

// Delete deletes documents from a collection based on their document IDs.
// It first checks the connection to Milvus. If the connection is not ready,
// an error is returned. Then, it converts the document IDs to strings and
// queries the collection to get the IDs of the documents to be deleted.
// If the query is successful, it converts the IDs to strings and deletes the
// documents using the `Delete` method of the Milvus client. If there are no
// documents to be deleted, the method returns nil.
//
// Parameters:
//   - ctx: The context.Context used for the operation.
//   - collection: The name of the collection to delete documents from.
//   - documentIDs: The ID(s) of the documents to be deleted.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func (c *milvusClient) Delete(ctx context.Context, collection string, documentIDs ...int64) error {
	if err := c.checkConnection(); err != nil {
		return err
	}
	var docsID []string
	for _, id := range documentIDs {
		docsID = append(docsID, strconv.FormatInt(id, 10))
	}
	queryResult, err := c.client.Query(ctx, collection, []string{},
		fmt.Sprintf("document_id in [%s]", strings.Join(docsID, ",")),
		[]string{"id"},
	)
	if err != nil {
		return err
	}
	var ids []string
	for _, result := range queryResult {
		for i := 0; i < result.Len(); i++ {
			if id, err := result.GetAsInt64(i); err == nil {
				ids = append(ids, strconv.FormatInt(id, 10))
			}
		}
	}
	if len(ids) == 0 {
		return c.client.Delete(ctx, collection, "", fmt.Sprintf("id in [%s]", strings.Join(ids, ",")))
	}
	return nil
}

func (v MilvusConfig) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.IndexStrategy, validation.Required,
			validation.In(IndexStrategyDISKANN, IndexStrategyAUTOINDEX, IndexStrategyNoIndex)),
		validation.Field(&v.MetricType, validation.Required,
			validation.In(string(entity.COSINE), string(entity.L2), string(entity.IP))),
	)
}

func (c *milvusClient) Load(ctx context.Context, collection string, vector []float32) ([]*MilvusPayload, error) {
	if err := c.checkConnection(); err != nil {
		return nil, err
	}
	vs := []entity.Vector{entity.FloatVector(vector)}
	sp, _ := entity.NewIndexFlatSearchParam()
	result, err := c.client.Search(ctx, collection, []string{}, "", responseColumns, vs, ColumnNameVector, c.MetricType, 10, sp)
	if err != nil {
		return nil, err
	}
	var payload []*MilvusPayload
	for _, row := range result {
		for i := 0; i < row.ResultCount; i++ {
			var pr MilvusPayload
			if err = pr.FromResult(i, row); err != nil {
				return nil, err
			}
			payload = append(payload, &pr)
		}
	}
	return payload, nil
}

// MilvusModule is a variable of type fx.Option. It provides dependencies for Milvus configuration
// and client initialization.
//
// The provided dependencies include:
// - A function that reads the MilvusConfig from environment variables and validates it.
// - A function that creates a new Milvus client based on the provided config.
//
// This module can be used to configure and initialize a VectorDBClient.
var MilvusModule = fx.Options(
	fx.Provide(func() (*MilvusConfig, error) {
		cfg := MilvusConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		NewMilvusClient,
	),
)

// NewMilvusClient creates a new instance of VectorDBClient
func NewMilvusClient(cfg *MilvusConfig) (VectorDBClient, error) {
	client, err := connect(cfg)
	if err != nil {
		zap.S().Errorf("connect to milvus error %s ", err.Error())
	}
	return &milvusClient{
		client:     client,
		cfg:        cfg,
		MetricType: entity.MetricType(cfg.MetricType),
	}, nil
}

// checks if the connection to Milvus is ready
func (c *milvusClient) Save(ctx context.Context, collection string, payloads ...*MilvusPayload) error {
	var ids, documentIDs, chunks []int64
	var contents [][]byte
	var vectors [][]float32
	if err := c.checkConnection(); err != nil {
		return err
	}
	for _, payload := range payloads {
		ids = append(ids, payload.ID)
		documentIDs = append(documentIDs, payload.DocumentID)
		chunks = append(chunks, payload.Chunk)
		contents = append(contents, []byte(fmt.Sprintf(`{"content":"%s"}`, payload.Content)))
		vectors = append(vectors, payload.Vector)
	}
	if _, err := c.client.Insert(ctx, collection, "",
		entity.NewColumnInt64(ColumnNameID, ids),
		entity.NewColumnInt64(ColumnNameDocumentID, documentIDs),
		entity.NewColumnJSONBytes(ColumnNameContent, contents),
		entity.NewColumnFloatVector(ColumnNameVector, VectorDimension, vectors),
	); err != nil {
		return err
	}
	return nil
}

// indexStrategy returns the appropriate index strategy based on the configuration
// If the index strategy is AUTOINDEX, it returns a new instance of entity.IndexAUTOINDEX
// If the index strategy is DISKANN, it returns a new instance of entity.IndexDISKANN
// If the index strategy is not supported, it returns an error with a corresponding message
func (c *milvusClient) indexStrategy() (entity.Index, error) {
	switch c.cfg.IndexStrategy {
	case IndexStrategyAUTOINDEX:
		return entity.NewIndexAUTOINDEX(c.MetricType)
	case IndexStrategyDISKANN:
		return entity.NewIndexDISKANN(c.MetricType)
	}
	return nil, fmt.Errorf("index strategy %s not supported yet", c.cfg.IndexStrategy)
}

// CreateSchema creates a new schema in Milvus.
func (c *milvusClient) CreateSchema(ctx context.Context, name string) error {

	collExists, err := c.client.HasCollection(ctx, name)
	if err != nil {
		return err
	}
	if collExists {
		if err = c.client.DropCollection(ctx, name); err != nil {
			return err
		}
		collExists = false
	}
	schema := entity.NewSchema().WithName(name).
		WithField(entity.NewField().WithName(ColumnNameID).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
		WithField(entity.NewField().WithName(ColumnNameDocumentID).WithDataType(entity.FieldTypeInt64)).
		WithField(entity.NewField().WithName(ColumnNameContent).WithDataType(entity.FieldTypeJSON)).
		WithField(entity.NewField().WithName(ColumnNameVector).WithDataType(entity.FieldTypeFloatVector).WithDim(1536))
	if err = c.client.CreateCollection(ctx, schema, 2, milvus.WithAutoID(true)); err != nil {
		return err
	}

	if c.cfg.IndexStrategy != IndexStrategyNoIndex {
		indexStrategy, err := c.indexStrategy()
		if err != nil {
			return err
		}
		if err = c.client.CreateIndex(ctx, name, ColumnNameVector, indexStrategy, true); err != nil {
			return err
		}
	}
	return nil
}

// FromResult translates the SearchResult to MilvusPayload.
// It iterates through each field in the SearchResult and assigns the corresponding
// value to the fields in MilvusPayload.
// If the field is ColumnNameID, it assigns the value to the ID field in MilvusPayload.
// If the field is ColumnNameDocumentID, it assigns the value to the DocumentID field in MilvusPayload.
// If the field is ColumnNameContent, it attempts to unmarshal the value into a string map.
// If successful, it assigns the value of the "content" key to the Content field in MilvusPayload.
// Otherwise, it assigns the value directly to the Content field in MilvusPayload.
// If any error occurs during the process, it returns the error.
func (p *MilvusPayload) FromResult(i int, res milvus.SearchResult) error {
	var err error

	for _, field := range res.Fields {
		switch field.Name() {
		case ColumnNameID:
			p.ID, err = field.GetAsInt64(i)
		case ColumnNameDocumentID:
			p.DocumentID, err = field.GetAsInt64(i)
		case ColumnNameContent:
			row, err := field.GetAsString(i)
			if err != nil {
				continue
			}
			contentS := ""
			if err = json.Unmarshal([]byte(row), &contentS); err == nil {
				contentS = strings.ReplaceAll(contentS, "\n", "")
				content := make(map[string]string)
				if err = json.Unmarshal([]byte(contentS), &content); err == nil {
					p.Content = content[ColumnNameContent]
				}
			} else {
				content := make(map[string]string)
				if err = json.Unmarshal([]byte(row), &content); err == nil {
					p.Content = content[ColumnNameContent]
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// checkConnection checks if the Milvus connection is ready
func (c *milvusClient) checkConnection() error {
	// creates connection if not exists
	if c.client == nil {
		client, err := connect(c.cfg)
		if err != nil {
			zap.S().Error(err.Error())
			return fmt.Errorf("milvus is not initialized")
		}
		c.client = client
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// check connection status
	state, err := c.client.CheckHealth(ctx)
	if err != nil {
		return fmt.Errorf("client.CheckHealth error %s", err.Error())
	}
	if !state.IsHealthy {
		return fmt.Errorf("milvus is not ready  %s", strings.Join(state.Reasons, " "))
	}

	return nil
}

// connect creates a connection to the Milvus server with the provided configuration.
// It returns a Milvus client and an error if the connection fails.
// The connection is established using the provided MilvusConfig, which contains the server address, username, password, and retry rate limit.
// The function uses a context with a timeout of 2 seconds.
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//	client, err := connect(&MilvusConfig{
//		Address:        "localhost:19530",
//		Username:       "admin",
//		Password:       "password123",
//		RetryRateLimit: &RetryRateLimitOption{MaxRetry: 2, MaxBackoff: 2 * time.Second},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Note: The connect function is used internally by the NewMilvusClient function to establish a connection to the Milvus server.
func connect(cfg *MilvusConfig) (milvus.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return milvus.NewClient(ctx, milvus.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		RetryRateLimit: &milvus.RetryRateLimitOption{
			MaxRetry:   2,
			MaxBackoff: 2 * time.Second,
		},
	})
}
