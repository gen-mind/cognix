package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// TenantRepository is an interface that defines methods for retrieving users related to a tenant.
type (
	TenantRepository interface {
		GetUsers(ctx context.Context, tenantID uuid.UUID) ([]*model.User, error)
	}
	tenantRepository struct {
		db *pg.DB
	}
)

// GetUsers retrieves users for a specific tenant.
//
// Parameters:
// - ctx: the context.Context object.
// - tenantID: the UUID of the tenant.
//
// Returns:
// - []*model.User: a slice of User objects.
// - error: an error, if any.
func (r *tenantRepository) GetUsers(ctx context.Context, tenantID uuid.UUID) ([]*model.User, error) {
	users := make([]*model.User, 0)
	if err := r.db.Model(&users).Where("tenant_id = ?", tenantID).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "cannot get users")
	}
	return users, nil
}

// NewTenantRepository creates a new instance of the TenantRepository interface, using the provided *pg.DB.
func NewTenantRepository(db *pg.DB) TenantRepository {
	return &tenantRepository{db: db}
}
