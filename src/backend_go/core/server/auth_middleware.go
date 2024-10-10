package server

import (
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// ContextParamUser is a constant representing the key used to store user identity
// information in the request context.
const ContextParamUser = "CONTEXT_USER"

// AuthMiddleware is a middleware that performs authentication for incoming requests.
type AuthMiddleware struct {
	jwtService security.JWTService
	userRepo   repository.UserRepository
}

// NewAuthMiddleware creates a new AuthMiddleware with the given JWTService and UserRepository.
// Parameters:
// - jwtService: service for JWT operations
// - userRepo: repository for user operations
// Returns:
// - *AuthMiddleware: the new AuthMiddleware instance
func NewAuthMiddleware(jwtService security.JWTService,
	userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService,
		userRepo: userRepo}
}

// RequireAuth requires authentication for the given request.
func (m *AuthMiddleware) RequireAuth(c *gin.Context) {

	//Get the  bearer Token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		handleError(c, utils.ErrorUnauthorized.New("Authorization token is required"))
		c.Abort()
		return
	}

	extractedToken := strings.Split(tokenString, "Bearer ")

	if len(extractedToken) != 2 {
		handleError(c, utils.ErrorBadRequest.New("Incorrect format of authorization token"))
		c.Abort()
		return
	}

	claims, err := m.jwtService.ParseAndValidate(strings.TrimSpace(extractedToken[1]))
	if err != nil {
		handleError(c, utils.ErrorBadRequest.New("Token is not valid"))
		c.Abort()
		return
	}

	if claims.ExpiresAt != 0 && time.Now().Unix() > claims.ExpiresAt {
		handleError(c, utils.ErrorUnauthorized.New("token expired"))
		c.Abort()
		return
	}

	if claims.User, err = m.userRepo.GetByIDAndTenantID(c.Request.Context(), claims.User.ID, claims.User.TenantID); err != nil {
		handleError(c, utils.ErrorUnauthorized.Wrap(err, "wrong user"))
		c.Abort()
		return
	}
	c.Request = c.Request.WithContext(context.WithValue(
		c.Request.Context(), ContextParamUser, claims))
	c.Next()
}

// GetContextIdentity retrieves the identity of a user from the Gin context.
// It returns the identity if it exists, otherwise it returns an error.
// In case of a broken session, it returns an error with the message "broken session".
//
// Parameters:
// - c: a pointer to the Gin context
//
// Returns:
// - *security.Identity: the identity of the user
// - error: an error if the identity does not exist or the session is broken
func GetContextIdentity(c *gin.Context) (*security.Identity, error) {
	claims, ok := c.Request.Context().Value(ContextParamUser).(*security.Identity)
	if !ok {
		return nil, utils.ErrorPermission.New("broken session")
	}
	return claims, nil
}
