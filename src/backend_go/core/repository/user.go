package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// UserRepository is an interface that provides database operations with User model
// GetByUserName gets a user by username
// Parameters:
// - c: context
// - username: username of the user
// Returns:
// - *model.User: user
// - error: any error that occurred
// GetByIDAndTenantID gets a user by ID and Tenant ID
// Parameters:
// - c: context
// - id: ID of the user
// - tenantID: Tenant ID of the user
// Returns:
// - *model.User: user
// - error: any error that occurred
// IsUserExists checks if a user exists
// Parameters:
// - c: context
// - username: username of the user
// Returns:
// - bool: true if the user exists, false otherwise
// - error: any error that occurred
// RegisterUser registers a user
// Parameters:
// - c: context
// - user: user to register
// Returns:
// - error: any error that occurred
// Create creates a user
// Parameters:
// - c: context
// - user: user to create
// Returns:
// - error: any error that occurred
// Update updates a user
// Parameters:
// - c: context
// - user: user to update
// Returns:
// - error: any error that occurred
// UserRepository provides database operations with User model
type (
	UserRepository interface {
		GetByUserName(c context.Context, username string) (*model.User, error)
		GetByIDAndTenantID(c context.Context, id, tenantID uuid.UUID) (*model.User, error)
		IsUserExists(c context.Context, username string) (bool, error)
		RegisterUser(c context.Context, user *model.User) error
		Create(c context.Context, user *model.User) error
		Update(c context.Context, user *model.User) error
	}
	// UserRepository provides database operations with User model
	userRepository struct {
		db *pg.DB
	}
)

// NewUserRepository creates a new instance of the UserRepository interface, using the provided *pg.DB.
func NewUserRepository(db *pg.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByUserName retrieves a user from the database by their username.
// It takes a context and the username as parameters and returns a pointer to
// the User object and an error. If the user is not found, it returns nil and
// an error with a message indicating that the user could not be found.
//
// Parameters:
// - c: the context.
// - username: the username of the user to retrieve.
//
// Returns:
// - *model.User: a pointer to the User object.
// - error: an error object if there was an error retrieving the user.
func (u *userRepository) GetByUserName(c context.Context, username string) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(c).Model(&user).Where("user_name = ?", username).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find user")
	}
	return &user, nil
}

// GetByIDAndTenantID retrieves a user from the database based on their ID and tenant ID.
//
// Parameters:
// - c: the context.Context
// - id: the ID of the user
// - tenantID: the ID of the tenant
//
// Returns:
// - *model.User: the user matching the ID and tenant ID
// - error: an error if the user cannot be found or an internal error occurs
func (u *userRepository) GetByIDAndTenantID(c context.Context, id, tenantID uuid.UUID) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(c).Model(&user).
		Where("id = ?", id).
		Where("tenant_id = ?", tenantID).
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find user")
	}
	return &user, nil
}

// IsUserExists checks if a user with the given username exists in the database.
//
// Parameters:
// - c: the context.Context object to carry deadlines, cancellations and other request-scoped values across API boundaries.
// - username: the username of the user to check for existence.
//
// Returns:
// - bool: true if the user with the given username exists, false otherwise.
// - error: any error occurred during the execution of the function, or nil if there was no error.
func (u *userRepository) IsUserExists(c context.Context, username string) (bool, error) {
	exists, err := u.db.WithContext(c).Model(&model.User{}).Where("user_name = ?", username).Exists()
	if err != nil {
		return false, utils.NotFound.Wrap(err, "can not find user")
	}
	return exists, err
}

// Create inserts a new user into the database.
//
// Parameters:
// - c: the context.Context used for the database operation.
// - user: a pointer to the model.User object to be created.
//
// Returns:
// - error: an error if the database operation fails, nil otherwise.
func (u *userRepository) Create(c context.Context, user *model.User) error {
	if _, err := u.db.WithContext(c).Model(user).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create user")
	}
	return nil
}

// Update updates a user in the database based on the provided user object.
// It uses the user's ID to identify the record to be updated.
//
// Parameters:
// - c: the context.Context object for the operation.
// - user: the User object representing the user to be updated.
//
// Returns:
// - error: an error if the update operation fails, or nil if the update is successful.
func (u *userRepository) Update(c context.Context, user *model.User) error {
	if _, err := u.db.WithContext(c).Model(user).Where("id = ?", user.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update user")
	}
	return nil
}

// RegisterUser registers a new user in the system.
//
// Parameters:
// - c: context.Context - the context object.
// - user: *model.User - the user to register.
//
// Returns:
// - error - error if any occurred during the registration process.
func (u *userRepository) RegisterUser(c context.Context, user *model.User) error {
	return u.db.RunInTransaction(c, func(tx *pg.Tx) error {
		tenant := model.Tenant{
			ID:            user.TenantID,
			Name:          user.FirstName,
			Configuration: nil,
		}
		if _, err := tx.Model(&tenant).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not create tenant")
		}
		if _, err := tx.Model(user).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not create user")
		}
		if _, err := tx.Model(user.Defaults.EmbeddingModel).Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not create default embedding model")
		}

		return nil
	})
}
