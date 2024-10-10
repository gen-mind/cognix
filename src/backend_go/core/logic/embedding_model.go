package logic

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-pg/pg/v10"
	"time"
)

type (

	// EmbeddingModelBL is an interface that represents business logic for embedding models.
	EmbeddingModelBL interface {
		GetAll(ctx context.Context, user *model.User, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error)
		Create(ctx context.Context, user *model.User, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error)
		Update(ctx context.Context, user *model.User, id int64, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error)
		Delete(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error)
		Restore(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error)
	}

	// embeddingModelBL is a struct that represents the business logic for embedding models.
	embeddingModelBL struct {
		emRepo repository.EmbeddingModelRepository
	}
)

// Create creates a new embedding model in the system.
//
// Parameters:
// - ctx: the context.Context object.
// - user: a pointer to the User struct representing the user associated with the operation.
// - em: a pointer to the EmbeddingModelParam struct containing the data for the new embedding model.
//
// Returns:
// - *EmbeddingModel: a pointer to the EmbeddingModel struct representing the created embedding model.
// - error: an error if the operation fails, nil otherwise.
func (b *embeddingModelBL) Create(ctx context.Context, user *model.User, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error) {
	embeddingModel := model.EmbeddingModel{
		TenantID:     user.TenantID,
		ModelID:      em.ModelID,
		ModelName:    em.ModelName,
		ModelDim:     em.ModelDim,
		URL:          em.URL,
		IsActive:     em.IsActive,
		CreationDate: time.Now().UTC(),
	}
	if err := b.emRepo.Create(ctx, &embeddingModel); err != nil {
		return nil, err
	}
	return &embeddingModel, nil
}

// Update updates an existing embedding model in the system.
//
// Parameters:
// - ctx: the context.Context object.
// - user: a pointer to the User struct representing the user associated with the operation.
// - id: the ID of the embedding model to update.
// - em: a pointer to the EmbeddingModelParam struct containing the updated data.
//
// Returns:
// - *EmbeddingModel: a pointer to the EmbeddingModel struct representing the updated embedding model.
// - error: an error if the operation fails, nil otherwise.
func (b *embeddingModelBL) Update(ctx context.Context, user *model.User, id int64, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error) {
	existingEM, err := b.emRepo.GetByID(ctx, user.TenantID, id)
	if err != nil {
		return nil, err
	}
	existingEM.ModelID = em.ModelID
	existingEM.ModelName = em.ModelName
	existingEM.ModelDim = em.ModelDim
	existingEM.URL = em.URL
	existingEM.IsActive = em.IsActive
	existingEM.LastUpdate = pg.NullTime{time.Now().UTC()}
	if err = b.emRepo.Update(ctx, existingEM); err != nil {
		return nil, err
	}
	return existingEM, nil
}

func (b *embeddingModelBL) Delete(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error) {
	existingEM, err := b.emRepo.GetByID(ctx, user.TenantID, id)
	if err != nil {
		return nil, err
	}
	existingEM.DeletedDate = pg.NullTime{time.Now().UTC()}
	if err = b.emRepo.Update(ctx, existingEM); err != nil {
		return nil, err
	}
	return existingEM, nil
}

func (b *embeddingModelBL) Restore(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error) {
	existingEM, err := b.emRepo.GetByID(ctx, user.TenantID, id)
	if err != nil {
		return nil, err
	}
	existingEM.DeletedDate = pg.NullTime{time.Time{}}
	if err = b.emRepo.Update(ctx, existingEM); err != nil {
		return nil, err
	}
	return existingEM, nil
}

// GetAll retrieves all embedding models based on the provided parameters.
// Parameters:
// - ctx: the context.Context object.
// - user: a pointer to the User struct representing the user associated with the operation.
// - param: a pointer to the ArchivedParam struct containing the parameters for the retrieval.
// Returns:
// - []*EmbeddingModel: a slice of pointers to EmbeddingModel structs representing the retrieved embedding models.
// - error: an error if the retrieval fails, nil otherwise.
func (b *embeddingModelBL) GetAll(ctx context.Context, user *model.User, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error) {
	return b.emRepo.GetAll(ctx, user.TenantID, param)
}

// GetByID retrieves an embedding model by its ID.
//
// Parameters:
// - ctx: the context.Context object.
// - user: a pointer to the User struct representing the user associated with the operation.
// - id: the ID of the embedding model to retrieve.
//
// Returns:
// - *EmbeddingModel: a pointer to the EmbeddingModel struct representing the retrieved embedding model.
// - error: an error if the retrieval fails, nil otherwise.
func (b *embeddingModelBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error) {
	return b.emRepo.GetByID(ctx, user.TenantID, id)
}

// NewEmbeddingModelBL creates a new instance of the EmbeddingModelBL interface
// with the provided EmbeddingModelRepository.
//
// Parameters:
//   - emRepo: an instance of the EmbeddingModelRepository interface used to
//     interact with the underlying data storage.
//
// Returns:
// - EmbeddingModelBL: an instance of the EmbeddingModelBL interface.
func NewEmbeddingModelBL(emRepo repository.EmbeddingModelRepository) EmbeddingModelBL {
	return &embeddingModelBL{emRepo: emRepo}
}
