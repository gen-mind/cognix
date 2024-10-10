package logic

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/google/uuid"
)

type (

	// TenantBL represents the business logic interface for managing tenants and users in the system.
	TenantBL interface {
		GetUsers(ctx context.Context, user *model.User) ([]*model.User, error)
		AddUser(ctx context.Context, user *model.User, email, role string) (*model.User, error)
		UpdateUser(ctx context.Context, user *model.User, id uuid.UUID, role string) (*model.User, error)
	}

	// tenantBL represents the business logic implementation for managing tenants and users.
	tenantBL struct {
		tenantRepo repository.TenantRepository
		userRepo   repository.UserRepository
	}
)

// GetUsers retrieves a list of users based on the provided user in the context.
// If the user has no roles or has the role "user", an access denied error is returned.
// Otherwise, the tenant repository is used to retrieve the list of users associated with the user's tenant ID.
// The function returns a list of users and an error.
func (b *tenantBL) GetUsers(ctx context.Context, user *model.User) ([]*model.User, error) {
	if len(user.Roles) == 0 || user.Roles[0] == model.RoleUser {
		return nil, utils.ErrorPermission.New("access denied")
	}
	return b.tenantRepo.GetUsers(ctx, user.TenantID)
}

// NewTenantBL is a function that creates a new instance of the TenantBL interface.
// It takes a TenantRepository and a UserRepository as input parameters and returns a TenantBL.
// The TenantBL implementation uses the provided TenantRepository and UserRepository to interact with the data layer.
// The TenantBL interface provides methods for managing tenants and users, such as retrieving a list of users, adding a new user, and updating a user.
func NewTenantBL(tenantRepo repository.TenantRepository, userRepo repository.UserRepository) TenantBL {
	return &tenantBL{tenantRepo: tenantRepo, userRepo: userRepo}
}

// AddUser adds a new user to the system with the provided email and role.
// It checks if the user already exists and returns an error if so.
// It creates a new User object with the provided email, role, and user ID generated using uuid.New().
// The User object is then passed to the UserRepository's Create method to persist it in the database.
// The function returns the newly created User object and an error if any.
func (b *tenantBL) AddUser(ctx context.Context, user *model.User, email, role string) (*model.User, error) {
	exists, err := b.userRepo.IsUserExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrorBadRequest.New("user already exists")
	}
	newUser := &model.User{
		ID:         uuid.New(),
		TenantID:   user.TenantID,
		UserName:   email,
		FirstName:  "",
		LastName:   "",
		ExternalID: "",
		Roles:      model.StringSlice{role},
	}
	if err := b.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

// UpdateUser updates the role of a user identified by the provided ID.
// It retrieves the user from the UserRepository based on the provided ID and the tenant ID of the user in the context.
// If the user is not found, an error is returned.
// The user's role is updated with the provided role value.
// The updated user is then passed to the UserRepository's Update method to persist the changes in the database.
// The function returns the updated user and an error if any.
func (b *tenantBL) UpdateUser(ctx context.Context, user *model.User, id uuid.UUID, role string) (*model.User, error) {
	user, err := b.userRepo.GetByIDAndTenantID(ctx, id, user.TenantID)
	if err != nil {
		return nil, err
	}
	user.Roles = model.StringSlice{role}
	if err = b.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
