package model

import "github.com/google/uuid"

const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
)

type Defaults struct {
	EmbeddingModel *EmbeddingModel
}

// User represents a user in the system.
type User struct {
	ID             uuid.UUID       `json:"id,omitempty"`
	TenantID       uuid.UUID       `json:"tenant_id,omitempty"`
	UserName       string          `json:"user_name,omitempty"`
	FirstName      string          `json:"first_name,omitempty"`
	LastName       string          `json:"last_name,omitempty"`
	ExternalID     string          `json:"-"`
	Roles          StringSlice     `json:"roles,omitempty" pg:",array"`
	Tenant         *Tenant         `json:"tenant,omitempty" pg:"rel:has-one"`
	EmbeddingModel *EmbeddingModel `json:"-" pg:"rel:has-one,fk:tenant_id,join_fk:tenant_id"`
	Defaults       *Defaults       `json:"-" pg:"-"`
}

// HasRoles checks if the user has any of the specified roles.
// It returns true if the user has at least one of the specified roles,
// otherwise it returns false.
//
// Parameters:
// - role: variadic parameter that represents the roles to check.
//
// Returns:
// - bool: true if the user has at least one of the specified roles, false otherwise.
func (u *User) HasRoles(role ...string) bool {
	for _, r := range role {
		for _, ur := range u.Roles {
			if ur == r {
				return true
			}
		}
	}
	return false
}
