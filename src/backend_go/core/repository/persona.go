package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// PersonaRepository is an interface that defines the methods for interacting with
// the personas table in the database.
// - GetAll returns all personas based on the provided context, tenant ID, and archived flag.
// - GetByID returns a persona based on the provided context, ID, tenant ID, and additional relations.
// - Create creates a new persona in the database based on the provided context and persona object.
// - Update updates an existing persona in the database based on the provided context and persona object.
// - Archive marks a persona as archived in the database based on the provided context and persona object.
// - IsExists checks if a persona with the provided ID and tenant ID exists in the database.
//
// personaRepository is an implementation of PersonaRepository.
// It holds the reference to the database connection.
type (
	PersonaRepository interface {
		GetAll(ctx context.Context, tenantID uuid.UUID, archived bool) ([]*model.Persona, error)
		GetByID(ctx context.Context, id int64, tenantID uuid.UUID, relations ...string) (*model.Persona, error)
		Create(ctx context.Context, persona *model.Persona) error
		Update(ctx context.Context, persona *model.Persona) error
		Archive(ctx context.Context, persona *model.Persona) error
		IsExists(ctx context.Context, id int64, tenantID uuid.UUID) (bool, error)
	}
	personaRepository struct {
		db *pg.DB
	}
)

// IsExists checks if a persona with the given ID and tenant ID exists in the database.
// It returns true if the persona exists, otherwise returns false.
// If there is an error while checking the existence of the persona, it returns an error.
func (r *personaRepository) IsExists(ctx context.Context, id int64, tenantID uuid.UUID) (bool, error) {
	exist, err := r.db.WithContext(ctx).Model(&model.Persona{}).
		Where("id = ?", id).Where("tenant_id = ?", tenantID).
		Exists()
	if err != nil {
		return false, utils.NotFound.Wrap(err, "can not find persona")
	}
	return exist, nil
}

// NewPersonaRepository creates a new instance of the PersonaRepository interface
// using the provided *pg.DB.
func NewPersonaRepository(db *pg.DB) PersonaRepository {
	return &personaRepository{db: db}
}

// GetAll retrieves all the personas from the repository based on the provided tenant ID and archived flag.
// If the archived flag is false, only non-archived personas are returned.
// The personas are returned as a slice of pointers to Persona models.
// An error is returned if any issues occur during the retrieval.
func (r *personaRepository) GetAll(ctx context.Context, tenantID uuid.UUID, archived bool) ([]*model.Persona, error) {
	personas := make([]*model.Persona, 0)
	stm := r.db.WithContext(ctx).Model(&personas).
		Relation("LLM.model_id").Relation("LLM.name").Relation("LLM.id").
		Relation("LLM.endpoint").
		Where("persona.tenant_id = ?", tenantID)

	if !archived {
		stm = stm.Where("persona.deleted_date IS NULL")
	}
	if err := stm.Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "personas not found")
	}
	return personas, nil
}

// GetByID retrieves a persona by ID and tenant ID. It returns the persona and an error,
// if any occurred. It allows for specifying additional relations to be eager-loaded.
// Relations "LLM" and "Prompt" are always eagerly loaded and don't need to be specified.
func (r *personaRepository) GetByID(ctx context.Context, id int64, tenantID uuid.UUID, relations ...string) (*model.Persona, error) {
	var persona model.Persona
	stm := r.db.WithContext(ctx).Model(&persona).
		Relation("LLM").
		Relation("Prompt").
		Where("persona.id = ?", id).
		Where("persona.tenant_id = ?", tenantID)
	for _, relation := range relations {
		if relation == "LLM" || relation == "Prompt" {
			continue
		}
		stm = stm.Relation(relation)
	}
	if err := stm.First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "persona not found")
	}
	return &persona, nil
}

// Create inserts a new persona into the database. This method performs the following steps:
// 1. Inserts the associated LLM model into the LLM table.
// 2. Sets the LlmID field of the persona to the ID of the inserted LLM model.
// 3. Inserts the persona into the personas table.
// 4. Sets the PersonaID field of the prompt to the ID of the inserted persona.
// 5. Inserts the prompt into the prompts table.
// If any of these steps fail, an error is returned.
//
// ctx: The context.Context object.
// persona: The persona to be created.
// Returns: An error if the persona creation fails, or nil if it succeeds.
func (r *personaRepository) Create(ctx context.Context, persona *model.Persona) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(persona.LLM).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not insert LLM")
		}
		persona.LlmID = persona.LLM.ID
		if _, err := tx.Model(persona).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not insert persona")
		}
		persona.Prompt.PersonaID = persona.ID
		if _, err := tx.Model(persona.Prompt).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not insert prompt")
		}
		return nil
	})
}

// Update updates the given persona in the repository.
func (r *personaRepository) Update(ctx context.Context, persona *model.Persona) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(persona.LLM).Where("id = ?", persona.LLM.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not update LLM")
		}
		persona.LlmID = persona.LLM.ID
		if _, err := tx.Model(persona).Where("id = ?", persona.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not update persona")
		}
		persona.Prompt.PersonaID = persona.ID
		if _, err := tx.Model(persona.Prompt).Where("id = ?", persona.Prompt.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not update prompt")
		}
		return nil
	})
}

// Archive archives a persona by setting the `deleted_date` field in the `personas` table to the provided date and
// updating the `deleted_date` and `last_update` fields in the `chat_sessions` table for the corresponding chat sessions.
// It runs the database operations in a transaction using the given context and returns any error that occurred.
//
// Parameters:
// - ctx: The context.Context to use for the transaction.
// - persona: The persona to archive.
//
// Returns:
// - error: If an error occurs during the database operations.
//
// Note: This method is implemented by the `personaRepository` type.
func (r *personaRepository) Archive(ctx context.Context, persona *model.Persona) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(&model.ChatSession{}).Where("persona_id = ?", persona.ID).
			Set("deleted_date = ?", persona.DeletedDate).
			Update(); err != nil {
			return utils.Internal.Wrap(err, "can not set deleted_date for chat sessions")
		}
		if _, err := tx.Model(&model.Persona{}).Where("id = ?", persona.ID).
			Set("deleted_date = ?", persona.DeletedDate).
			Set("last_update = ?", persona.LastUpdate).
			Update(); err != nil {
			return utils.Internal.Wrap(err, "can not set deleted_date for chat sessions")
		}
		return nil
	})
}
