package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type (
	ConnectorBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error)
		Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Connector, error)
	}
	connectorBL struct {
		connectorRepo  repository.ConnectorRepository
		credentialRepo repository.CredentialRepository
		//messenger      messaging.Client
	}
)

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

func NewConnectorBL(connectorRepo repository.ConnectorRepository,
	credentialRepo repository.CredentialRepository,
	// messenger messaging.Client,
) ConnectorBL {
	return &connectorBL{
		connectorRepo:  connectorRepo,
		credentialRepo: credentialRepo,
		//messenger:      messenger,
	}
}

func (c *connectorBL) Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error) {
	param.CredentialID = decimal.NullDecimal{Valid: false}
	tenantID := uuid.NullUUID{Valid: false}
	if param.Shared {
		tenantID.Valid = true
		tenantID.UUID = user.TenantID
	}
	conn := model.Connector{
		CredentialID:            decimal.NullDecimal{Valid: false},
		Name:                    param.Name,
		Type:                    model.SourceType(param.Source),
		ConnectorSpecificConfig: param.ConnectorSpecificConfig,
		RefreshFreq:             param.RefreshFreq,
		UserID:                  user.ID,
		TenantID:                tenantID,
		Disabled:                param.Disabled,
		LastAttemptStatus:       model.ConnectorStatusReadyToProcessed,
		CreationDate:            time.Now().UTC(),
	}
	//if param.CredentialID.Valid {
	//	cred, err := c.credentialRepo.GetByID(ctx, param.CredentialID.Decimal.IntPart(), user.TenantID, user.ID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if cred.Type != model.SourceType(param.Type) {
	//		return nil, utils.ErrorBadRequest.New("wrong credential source")
	//	}
	//	conn.CredentialID = decimal.NewNullDecimal(cred.ID)
	//}

	if err := c.connectorRepo.Create(ctx, &conn); err != nil {
		return nil, err
	}
	//if err := c.messenger.Publish(ctx, model.TopicUpdateConnector, &proto.Body{Payload: &proto.Body_Trigger{Trigger: &proto.ConnectorRequest{Id: conn.ID.IntPart()}}}); err != nil {
	//	return nil, err
	//}
	return &conn, nil
}

func (c *connectorBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error) {
	conn, err := c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	//if param.CredentialID.Valid {
	//	cred, err := c.credentialRepo.GetByID(ctx, param.CredentialID.Decimal.IntPart(), user.TenantID, user.ID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if cred.Type != conn.Type {
	//		return nil, utils.ErrorBadRequest.New("wrong credential source")
	//	}
	//}
	conn.ConnectorSpecificConfig = param.ConnectorSpecificConfig
	//conn.CredentialID = param.CredentialID
	conn.Name = param.Name
	conn.RefreshFreq = param.RefreshFreq
	tenantID := uuid.NullUUID{Valid: false}
	if param.Shared {
		tenantID.Valid = true
		tenantID.UUID = user.TenantID
	}
	conn.TenantID = tenantID
	conn.Disabled = param.Disabled
	conn.LastUpdate = pg.NullTime{time.Now().UTC()}

	if err = c.connectorRepo.Update(ctx, conn); err != nil {
		return nil, err
	}
	//if err = c.messenger.Publish(ctx, model.TopicUpdateConnector, &proto.Body{Payload: &proto.Body_Trigger{Trigger: &proto.ConnectorRequest{Id: conn.ID.IntPart()}}}); err != nil {
	//	return nil, err
	//}
	return conn, nil
}

func (c *connectorBL) GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error) {
	return c.connectorRepo.GetAllByUser(ctx, user.TenantID, user.ID)
}

func (c *connectorBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error) {
	return c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
}
