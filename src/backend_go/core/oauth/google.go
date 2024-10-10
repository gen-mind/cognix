package oauth

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"time"
)

// Constants for Google OAuth state and code names.
const (
	StateNameGoogle   = "state"
	CodeNameGoogle    = "code"
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

var googleAuthScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}
var googleDriveScopes []string = []string{"https://www.googleapis.com/auth/drive.readonly",
	"https://www.googleapis.com/auth/drive.metadata.readonly",
	"https://www.googleapis.com/auth/drive.activity.readonly",
}

// GoogleConfig represents the configuration for Google OAuth service.
type GoogleConfig struct {
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"`
	GoogleSecret   string `env:"GOOGLE_SECRET"`
	RedirectURL    string `env:"OAUTH_REDIRECT_URL" envDefault:"https://rag.cognix.ch/api/oauth"`
}

// googleProvider represents a provider implementation for Google OAuth client.
type googleProvider struct {
	config     *oauth2.Config
	httpClient *resty.Client
}

// NewGoogleProvider creates a new implementation of the google oAuth client.
// It takes a GoogleConfig pointer and a redirectURL string as input and returns a Proxy interface.
// The returned value is a pointer to a googleProvider struct.
// The googleProvider struct contains an httpClient (a resty.Client) and a config (an oauth2.Config).
// The config is initialized with the client ID, client secret, endpoint, redirect URL, and scopes.
// The config is used for authorization and authentication with the google OAuth service.
// The googleProvider struct also implements the methods defined in the Proxy interface.
// This function is used to create a new Google OAuth provider in the application.
func NewGoogleProvider(cfg *GoogleConfig, redirectURL string) Proxy {
	return &googleProvider{
		httpClient: resty.New().SetTimeout(time.Minute),
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  fmt.Sprintf("%s/google/callback", redirectURL),
			Scopes:       append(googleAuthScopes, googleDriveScopes...),
		},
	}
}

// GetAuthURL returns the authentication URL for the Google provider.
// The redirectURL parameter is used to construct the full callback URL.
// The state parameter is an OAuth state string that will be returned
// along with the user's authorization code.
// This method modifies the RedirectURL field of the googleProvider's config.
// It calls the AuthCodeURL method of the config to generate the URL.
// The ApprovalForce value is passed as the second argument to the AuthCodeURL method.
// The generated URL is returned along with any potential error that may occur.
func (g *googleProvider) GetAuthURL(ctx context.Context, redirectURL, state string) (string, error) {
	g.config.RedirectURL = fmt.Sprintf("%s", redirectURL)

	return g.config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce), nil
}

// ExchangeCode exchanges the authorization code for an access token and retrieves
// the user's identity information from the Google API.
// It returns an IdentityResponse containing the user's identity details and the
// access and refresh tokens.
func (g *googleProvider) ExchangeCode(ctx context.Context, code string) (*IdentityResponse, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, utils.Internal.Wrapf(err, "code exchange wrong: %s", err.Error())
	}
	response, err := g.httpClient.R().Get(oauthGoogleUrlAPI + token.AccessToken)
	if err = utils.WrapRestyError(response, err); err != nil {
		return nil, utils.Internal.Newf("failed getting user info: %s", err.Error())
	}

	contents := response.Body()

	var data IdentityResponse
	if err = json.Unmarshal(contents, &data); err != nil {
		return nil, utils.Internal.Wrapf(err, "can not marshal google response")
	}
	data.AccessToken = token.AccessToken
	data.RefreshToken = token.RefreshToken
	data.Token = token
	return &data, nil
}

// RefreshToken refreshes the provided OAuth2 token and returns the refreshed token
//
// Parameters:
// - token: The OAuth2 token to be refreshed
//
// Returns:
// - *oauth2.Token: The refreshed OAuth2 token
// - error: Any error that occurred during the token refreshing process
func (g *googleProvider) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	return g.config.TokenSource(context.Background(), token).Token()
}
