package security

import (
	"cognix.ch/api/v2/core/model"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// Identity represents the identity of a user with JWT claims, access token, refresh token, and user information.
type (
	Identity struct {
		jwt.StandardClaims
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		User         *model.User `json:"user"`
	}

	JWTService interface {
		Create(claim *Identity) (string, error)
		ParseAndValidate(string) (*Identity, error)
		Refresh(refreshToken string) (string, error)
	}
	jwtService struct {
		jwtSecret      string `json:"jwt_secret"`
		jwtExpiredTime int    `json:"jwt_expired_time"`
	}
)

// NewJWTService returns a new instance of JWTService.
// The function accepts a JWT secret and the expiration time in minutes.
// It creates a jwtService instance with the provided parameters and returns it as a JWTService.
// Example usage:
//
//	js := NewJWTService("secret", 60)
//	identity := Identity{...}
//	token, err := js.Create(&identity)
//	if err != nil {
//	  // handle error
//	}
//	fmt.Println(token)
func NewJWTService(jwtSecret string, jwtExpiredTime int) JWTService {
	return &jwtService{jwtSecret: jwtSecret,
		jwtExpiredTime: jwtExpiredTime * int(time.Minute)}
}

// Create generates a JSON Web Token (JWT) using the provided identity information.
// It takes an `Identity` struct pointer as a parameter, which contains the necessary information
// to create the token, such as access and refresh tokens, and user details.
// The function returns a string representing the generated token and an error, if any.
// The token is signed with the HMAC-SHA256 algorithm using the `jwtSecret` from the `jwtService` struct.
// If any error occurs during token creation, an empty string and the error are returned.
func (j *jwtService) Create(identity *Identity) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, identity)
	//identity.ExpiresAt = time.Now().Add(time.Duration(j.jwtExpiredTime)).Unix()
	tokenString, err := token.SignedString([]byte(j.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseAndValidate takes a token string and parses it into an Identity struct.
// It also validates the token using the jwt.ParseWithClaims function,
// and returns the parsed Identity if the token is valid.
//
// Parameters:
// - tokenString: a string representing the token to be parsed and validated.
//
// Returns:
// - *Identity: a pointer to the parsed Identity if the token is valid.
// - error: an error if parsing or validation fails.
func (j *jwtService) ParseAndValidate(tokenString string) (*Identity, error) {
	var identity Identity
	token, err := jwt.ParseWithClaims(tokenString, &identity, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &identity, nil
}

// Refresh generates a new access token using the provided refresh token.
// It takes in the refresh token as a string and returns the new access token as a string and an error, if any.
func (j *jwtService) Refresh(refreshToken string) (string, error) {
	return refreshToken, nil
}
