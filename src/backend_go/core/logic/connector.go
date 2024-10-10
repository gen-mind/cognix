package logic

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

const minRefreshFreq = 3600

type (

	// ConnectorBL represents the business logic interface for managing connectors.
	ConnectorBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error)
		Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Connector, error)
	}

	// connectorBL represents the business logic implementation for managing connectors.
	// It contains a reference to the connector repository for data access.
	connectorBL struct {
		connectorRepo repository.ConnectorRepository
		//messenger      messaging.Client
	}
)

// Archive archives or restores a connector.
// If restore is false, the connector will be marked as deleted by setting the DeletedDate field to the current date and time.
// If restore is true, the DeletedDate field will be cleared, indicating that the connector is no longer deleted.
// The LastUpdate field will be updated to the current date and time.
// The user must be the owner of the connector, an admin, or a super admin to perform this operation.
//
// Parameters:
// - ctx: the context.Context object for request cancellation and deadline propagation.
// - user: the user performing the operation.
// - id: the ID of the connector to archive or restore.
// - restore: a boolean value indicating whether to archive or restore the connector.
//
// Returns:
// - *model.Connector: the archive or restored connector.
// - error: an error object if an error occurs, otherwise it will be nil.
func (c *connectorBL) Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Connector, error) {
	connector, err := c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	if !(connector.UserID == user.ID || user.HasRoles(model.RoleAdmin, model.RoleSuperAdmin)) {
		return nil, utils.ErrorPermission.New("permission denied")
	}
	if !restore {
		connector.DeletedDate = pg.NullTime{time.Now().UTC()}
	} else {
		connector.DeletedDate = pg.NullTime{}
	}
	connector.LastUpdate = pg.NullTime{time.Now().UTC()}
	if err = c.connectorRepo.Update(ctx, connector); err != nil {
		return nil, err
	}
	return connector, nil
}

// NewConnectorBL creates a new instance of ConnectorBL implementation.
//
// Parameters:
// - connectorRepo: the ConnectorRepository used to access connector data.
//
// Returns:
// - ConnectorBL: an instance of ConnectorBL interface.
func NewConnectorBL(connectorRepo repository.ConnectorRepository) ConnectorBL {
	return &connectorBL{connectorRepo: connectorRepo}
}

// Create creates a new connector with the provided parameters.
// If the connector is shared, the tenant ID of the user will be set as the tenant ID of the connector.
// The connector's status is set to "ready to processed" and its creation date is set to the current date and time.
// If the specified refresh frequency is less than the minimum refresh frequency,
// the connector's refresh frequency will be set to the minimum refresh frequency.
//
// Parameters:
// - ctx: the context.Context object for request cancellation and deadline propagation.
// - user: the user performing the operation.
// - param: the parameters to create the connector.
//
// Returns:
// - *model.Connector: the created connector.
// - error: an error object if an error occurs, otherwise it will be nil.
func (c *connectorBL) Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error) {

	tenantID := uuid.NullUUID{Valid: false}
	if param.Shared {
		tenantID.Valid = true
		tenantID.UUID = user.TenantID
	}
	conn := model.Connector{
		Name:                    param.Name,
		Type:                    model.SourceType(param.Source),
		ConnectorSpecificConfig: param.ConnectorSpecificConfig,
		RefreshFreq:             param.RefreshFreq,
		UserID:                  user.ID,
		TenantID:                tenantID,
		Status:                  model.ConnectorStatusReadyToProcessed,
		CreationDate:            time.Now().UTC(),
	}

	if conn.RefreshFreq < minRefreshFreq {
		conn.RefreshFreq = minRefreshFreq
	}
	if err := c.connectorRepo.Create(ctx, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

// Update updates the details of a connector.
// It updates the connector's ConnectorSpecificConfig, Name, RefreshFreq, TenantID, and LastUpdate fields.
// If the connector is shared, the TenantID field will be set to the user's TenantID.
// The LastUpdate field will be updated to the current date and time.
// If the Status field is provided, it will be updated.
// If the specified RefreshFreq is less than the minimum refresh frequency, it will be set to the minimum refresh frequency.
//
// Parameters:
// - ctx: the context.Context object for request cancellation and deadline propagation.
// - id: the ID of the connector to update.
// - user: the user performing the operation.
// - param: the parameters to update the connector.
//
// Returns:
// - *model.Connector: the updated connector.
// - error: an error object if an error occurs, otherwise it will be nil.
func (c *connectorBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error) {
	conn, err := c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	conn.ConnectorSpecificConfig = param.ConnectorSpecificConfig
	conn.Name = param.Name
	conn.RefreshFreq = param.RefreshFreq
	tenantID := uuid.NullUUID{Valid: false}
	if param.Shared {
		tenantID.Valid = true
		tenantID.UUID = user.TenantID
	}

	conn.TenantID = tenantID
	conn.LastUpdate = pg.NullTime{time.Now().UTC()}
	if param.Status != "" {
		conn.Status = param.Status
	}
	if conn.RefreshFreq < minRefreshFreq {
		conn.RefreshFreq = minRefreshFreq
	}
	if err = c.connectorRepo.Update(ctx, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

// GetAll returns all connectors associated with the user.
//
// Parameters:
// - ctx: the context.Context object for request cancellation and deadline propagation.
// - user: the user object for which connectors are retrieved.
//
// Returns:
// - []*model.Connector: a slice containing all connectors associated with the user.
// - error: an error object if an error occurs, otherwise it will be nil.
func (c *connectorBL) GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error) {
	return c.connectorRepo.GetAllByUser(ctx, user.TenantID, user.ID)
}

// GetByID retrieves a connector by its ID and checks if the user has access to it.
//
// Parameters:
// - ctx: the context.Context object for request cancellation and deadline propagation.
// - user: the user performing the operation.
// - id: the ID of the connector to retrieve.
//
// Returns:
// - *model.Connector: the retrieved connector.
// - error: an error object if an error occurs, otherwise it will be nil.
func (c *connectorBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error) {
	return c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
}
