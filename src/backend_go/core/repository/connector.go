package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

// ConnectorRepository is an interface for accessing and manipulating connector data in the database.
type (
	ConnectorRepository interface {
		GetActive(ctx context.Context) ([]*model.Connector, error)
		GetAllByUser(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Connector, error)
		GetByIDAndUser(ctx context.Context, tenantID, userID uuid.UUID, id int64) (*model.Connector, error)
		GetByID(ctx context.Context, id int64) (*model.Connector, error)
		GetBySource(ctx context.Context, tenantID, userID uuid.UUID, source model.SourceType) (*model.Connector, error)
		Create(ctx context.Context, connector *model.Connector) error
		Update(ctx context.Context, connector *model.Connector) error
	}
	connectorRepository struct {
		db *pg.DB
	}
)

// GetBySource retrieves a connector from the database based on its source, tenant ID, and user ID.
// It returns the connector if found, otherwise returns an error.
// The method performs the following steps:
// - Prepare a Connector model
// - Use the provided context and database connection to execute a query that retrieves the connector based on the specified conditions
// - If the query is successful, return the connector
// - If the query fails, return a NotFound error wrapped with the appropriate message
// Note: The method uses the `pg` package for database operations and the `utils` package for error handling.
// The `model.SourceType` type is a string-based custom type used for representing a source type.
// The `model.Connector` type is a struct that represents a database table connector and has various properties.
// The `utils.NotFound` constant is an `ErrorWrap` type that represents a "not found" error with the corresponding HTTP status code.
// The `JSONMap` type is a custom type used for handling JSON data in the database.
func (r *connectorRepository) GetBySource(ctx context.Context, tenantID, userID uuid.UUID, source model.SourceType) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(ctx).Model(&connector).
		Where("source = ?", source).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("tenant_id = ?", tenantID), nil
		}).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "ca not find connector")
	}
	return &connector, nil
}

// NewConnectorRepository creates a new instance of the ConnectorRepository interface, using the provided *pg.DB.
func NewConnectorRepository(db *pg.DB) ConnectorRepository {
	return &connectorRepository{db: db}
}

// GetAllByUser retrieves all connectors associated with the given tenantID and userID.
// It returns a slice of model.Connector or an error if the retrieval failed.
func (r *connectorRepository) GetAllByUser(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Connector, error) {
	connectors := make([]*model.Connector, 0)
	if err := r.db.WithContext(ctx).Model(&connectors).
		Where("deleted_date is null").
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("tenant_id = ?", tenantID), nil
		}).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connectors")
	}
	return connectors, nil
}

// GetByIDAndUser retrieves a connector by its ID and user ID from the database.
// It returns the connector if found, otherwise it returns nil with an error.
func (r *connectorRepository) GetByIDAndUser(ctx context.Context, tenantID, userID uuid.UUID, id int64) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(ctx).Model(&connector).
		Where("id = ?", id).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("user_id = ?", userID).
				WhereOr("tenant_id = ?", tenantID), nil
		}).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connector")
	}
	return &connector, nil
}

// GetByID retrieves a connector from the database based on its ID.
// It queries the `connectors` table and joins the related documents and user's embedding model.
// The retrieved connector includes a map of its documents, where the key is the `SourceID` of each document.
// If the connector is not found, it returns an error with a message indicating the failure reason.
// Otherwise, it returns the connector object.
// The provided `ctx` serves as the cancellation context for the database operations.
// The `id` parameter specifies the ID of the connector to retrieve.
func (r *connectorRepository) GetByID(ctx context.Context, id int64) (*model.Connector, error) {
	var connector model.Connector
	if err := r.db.WithContext(ctx).Model(&connector).
		Relation("Docs").
		Relation("User.EmbeddingModel").
		Where("connector.id = ?", id).
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not load connector")
	}
	connector.DocsMap = make(map[string]*model.Document)
	for _, doc := range connector.Docs {
		connector.DocsMap[doc.SourceID] = doc
	}
	return &connector, nil
}

// Create creates a new connector in the database.
// It takes a context and a connector model as input.
// If the insertion operation fails, it returns an error with an appropriate message.
// Otherwise, it returns nil.
func (r *connectorRepository) Create(ctx context.Context, connector *model.Connector) error {
	stm := r.db.WithContext(ctx).Model(connector)
	if _, err := stm.Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create connector")
	}
	return nil
}

// Update updates the given connector in the database.
func (r *connectorRepository) Update(ctx context.Context, connector *model.Connector) error {
	if _, err := r.db.WithContext(ctx).Model(connector).Where("id = ?", connector.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update connector")
	}
	return nil
}

// GetActive retrieves active connectors with specific statuses and without deleted date.
//
// The method loads connectors with the following statuses:
// - READY_TO_PROCESS
// - COMPLETED_SUCCESSFULLY
// - COMPLETED_WITH_ERRORS
//
// It returns an array of model.Connector and an error, if any.
func (r *connectorRepository) GetActive(ctx context.Context) ([]*model.Connector, error) {
	connectors := make([]*model.Connector, 0)
	//todo ask Gian how to do this

	// load connectors with status that might be resending
	enabledStatuses := []string{model.ConnectorStatusReadyToProcessed, model.ConnectorStatusSuccess, model.ConnectorStatusError}

	if err := r.db.WithContext(ctx).
		Model(&connectors).
		Relation("Docs").
		Relation("User.EmbeddingModel").
		Where("status = any(?)", pg.Array(enabledStatuses)).
		Where("connector.deleted_date is null").
		Select(); err != nil {
		return nil, utils.Internal.Wrapf(err, "can not load connectors: %s ", err.Error())
	}
	return connectors, nil

}
