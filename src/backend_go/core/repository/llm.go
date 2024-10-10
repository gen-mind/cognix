package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// LLMRepository interface defines the methods for accessing data from the llms table.
type (
	LLMRepository interface {
		GetAll(ctx context.Context) ([]*model.LLM, error)
		GetByUserID(ctx context.Context, userID uuid.UUID) (*model.LLM, error)
	}
	llmRepository struct {
		db *pg.DB
	}
)

// GetByUserID retrieves an LLM record from the database based on the user ID.
func (r *llmRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.LLM, error) {
	var llm model.LLM
	if err := r.db.WithContext(ctx).Model(&llm).
		Join("inner join users on llm.tenant_id = users.tenant_id").
		Where("users.id = ?", userID).
		Limit(1).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find llm ")
	}
	return &llm, nil
}

// NewLLMRepository creates a new instance of the LLMRepository interface using the provided *pg.DB.
func NewLLMRepository(db *pg.DB) LLMRepository {
	return &llmRepository{db: db}
}

// GetAll returns all LLM records stored in the llms table.
// If there is an error retrieving the records, a NotFound error is returned.
// The NotFound error is wrapped with additional details.
// The method returns a slice of LLM records and an error.
// The slice may be empty if no records are found.
// The method uses the provided context to perform the database operation.
// If the provided context is canceled or times out, the method returns an error.
func (r *llmRepository) GetAll(ctx context.Context) ([]*model.LLM, error) {
	llms := make([]*model.LLM, 0)
	if err := r.db.WithContext(ctx).Model(&llms).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find llm")
	}
	return llms, nil
}
