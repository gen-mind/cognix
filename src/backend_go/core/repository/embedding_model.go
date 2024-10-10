package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

// EmbeddingModelRepository is an interface that specifies the methods for accessing embedding model data.
type (
	EmbeddingModelRepository interface {
		GetAll(ctx context.Context, tenantID uuid.UUID, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error)
		GetByID(ctx context.Context, tenantID uuid.UUID, id int64) (*model.EmbeddingModel, error)
		GetDefault(ctx context.Context, tenantID uuid.UUID) (*model.EmbeddingModel, error)
		Create(ctx context.Context, em *model.EmbeddingModel) error
		Update(ctx context.Context, em *model.EmbeddingModel) error
		Delete(ctx context.Context, em *model.EmbeddingModel) error
	}
	embeddingModelRepository struct {
		db *pg.DB
	}
)

// GetDefault retrieves the default embedding model for a given tenant.
// It queries the repository for an active embedding model with the specified tenant ID,
// and limits the result to one record. If no record is found,
// it returns a "not found" error with a wrapped error indicating the failure.
// If successful, it returns a pointer to the retrieved embedding model and nil error.
func (r *embeddingModelRepository) GetDefault(ctx context.Context, tenantID uuid.UUID) (*model.EmbeddingModel, error) {
	var em model.EmbeddingModel
	if err := r.db.Model(&em).Where("tenant_id = ?", tenantID).
		Where("is_active = true").
		Limit(1).
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "Cannot get default embedding model")
	}
	return &em, nil
}

// GetAll returns all the embedding models that belong to the given tenant.
// If the archived parameter is false, it only returns the non-archived embedding models.
// If the operation fails, it returns an error with a wrapped error message indicating the failure.
// The GetAll method returns a slice of embedding models and an error.
func (r *embeddingModelRepository) GetAll(ctx context.Context, tenantID uuid.UUID, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error) {
	ems := make([]*model.EmbeddingModel, 0)
	stm := r.db.WithContext(ctx).Model(&ems).Where("tenant_id = ?", tenantID)
	if !param.Archived {
		stm = stm.Where("deleted_date is null")
	}

	if err := stm.Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find embedding models")
	}
	return ems, nil
}

// GetByID retrieves an embedding model from the repository based on its ID and tenant ID.
// If no embedding model is found, it returns a not found error with a wrapped error message.
func (r *embeddingModelRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id int64) (*model.EmbeddingModel, error) {
	var em model.EmbeddingModel
	if err := r.db.WithContext(ctx).Model(&em).Where("id = ?", id).
		Where("tenant_id = ?", tenantID).
		Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find embedding models")
	}
	return &em, nil
}

func (r *embeddingModelRepository) Create(ctx context.Context, em *model.EmbeddingModel) error {
	if _, err := r.db.WithContext(ctx).Model(em).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create embedding models")
	}
	return nil
}

// Update updates the embedding model in the repository.
// It sets the LastUpdate field to the current UTC time.
// If the update operation fails, it returns an error with a wrapped error message indicating the failure.
// The Update method does not return any values.
func (r *embeddingModelRepository) Update(ctx context.Context, em *model.EmbeddingModel) error {
	em.LastUpdate = pg.NullTime{time.Now().UTC()}
	if _, err := r.db.WithContext(ctx).Model(em).
		Where("id = ?", em.ID).
		Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update embedding models")
	}
	return nil
}

func (r *embeddingModelRepository) Delete(ctx context.Context, em *model.EmbeddingModel) error {
	em.DeletedDate = pg.NullTime{time.Now().UTC()}
	if _, err := r.db.WithContext(ctx).Model(em).
		Where("id = ?", em.ID).
		Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update embedding models")
	}
	return nil
}

// NewEmbeddingModelRepository creates a new instance of the EmbeddingModelRepository interface, using the provided *pg.DB.
func NewEmbeddingModelRepository(db *pg.DB) EmbeddingModelRepository {
	return &embeddingModelRepository{
		db: db,
	}
}
