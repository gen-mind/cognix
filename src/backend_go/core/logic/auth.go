package logic

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/google/uuid"
	"time"
)

type (

	// AuthBL represents the business logic for user authentication and sign-up.
	AuthBL interface {
		Login(ctx context.Context, userName string) (*model.User, error)
		SignUp(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error)
		QuickLogin(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error)
	}

	// authBL represents the business logic for user authentication and sign-up.
	authBL struct {
		userRepo repository.UserRepository
		cfg      *Config
	}
)

// NewAuthBL creates a new instance of AuthBL.
//
// Parameters:
// - userRepo: an implementation of the UserRepository interface.
// - cfg: a pointer to the Config struct.
//
// Returns:
// - AuthBL: a new instance of AuthBL.
func NewAuthBL(userRepo repository.UserRepository,

	cfg *Config) AuthBL {
	return &authBL{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Login performs the authentication process for a user.
//
// Parameters:
// - ctx: the context for the operation.
// - userName: the username of the user to authenticate.
//
// Returns:
// - *model.User: the authenticated user.
// - error: an error if any occurred during the authentication process.
func (a *authBL) Login(ctx context.Context, userName string) (*model.User, error) {
	user, err := a.userRepo.GetByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignUp creates a new user and registers them in the system.
//
// Parameters:
// - ctx: the context for the operation.
// - identity: the identity information of the user obtained from the OAuth provider.
//
// Returns:
// - *model.User: the newly created and registered user.
// - error: an error if any occurred during the sign-up process.
func (a *authBL) SignUp(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error) {
	exists, err := a.userRepo.IsUserExists(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrorBadRequest.New("user already exists")
	}
	userID := uuid.New()
	tenantID := uuid.New()

	// create user  and default connector and embedding model
	user := model.User{
		ID:         userID,
		TenantID:   tenantID,
		UserName:   identity.Email,
		FirstName:  identity.GivenName,
		LastName:   identity.FamilyName,
		ExternalID: identity.ID,
		Roles:      model.StringSlice{model.RoleSuperAdmin},
		Defaults: &model.Defaults{
			EmbeddingModel: &model.EmbeddingModel{
				TenantID:     tenantID,
				ModelID:      a.cfg.DefaultEmbeddingModel,
				ModelName:    a.cfg.DefaultEmbeddingModel,
				ModelDim:     a.cfg.DefaultEmbeddingVectorSize,
				IsActive:     true,
				CreationDate: time.Now().UTC(),
			},
		},
	}
	if user.FirstName == "" {
		user.FirstName = identity.Name
	}
	if err = a.userRepo.RegisterUser(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// QuickLogin performs a quick login process for a user using their identity information from an OAuth provider.
//
// Parameters:
// - ctx: the context for the operation.
// - identity: the identity information of the user obtained from the OAuth provider.
//
// Returns:
// - *model.User: the authenticated user.
// - error: an error if any occurred during the quick login process.
func (a *authBL) QuickLogin(ctx context.Context, identity *oauth.IdentityResponse) (*model.User, error) {
	exists, err := a.userRepo.IsUserExists(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return a.SignUp(ctx, identity)
	}
	user, err := a.userRepo.GetByUserName(ctx, identity.Email)
	if err != nil {
		return nil, err
	}
	if user.ExternalID == "" {
		if identity.GivenName == "" {
			user.FirstName = identity.Name
		} else {
			user.FirstName = identity.GivenName
		}
		user.LastName = identity.FamilyName
		user.ExternalID = identity.ID
	}
	if err = a.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
