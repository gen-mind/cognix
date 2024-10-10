package oauth

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"time"
)

const (
	microsoftLoginURL = `https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&scope=%s&response_type=code&redirect_uri=%s`
	microsoftToken    = `https://login.microsoftonline.com/organizations/oauth2/v2.0/token`
)

const microsoftScope = "offline_access Files.Read.All Sites.ReadWrite.All"

const teamsScope = "ChannelMessage.Read.All Chat.Read Chat.ReadBasic Team.ReadBasic.All TeamSettings.Read.All ChannelSettings.Read.All Channel.ReadBasic.All Group.Read.All Directory.Read.All"

type (
	Config struct {
		Microsoft *MicrosoftConfig
		Google    *GoogleConfig
	}

	// MicrosoftConfig represents the configuration for Microsoft OAuth service.
	//
	MicrosoftConfig struct {
		ClientID     string `env:"MICROSOFT_CLIENT_ID,required"`
		ClientSecret string `env:"MICROSOFT_CLIENT_SECRET,required"`
		RedirectUL   string
	}
	microsoftExchangeCodeRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code,omitempty"`
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	tokenResponse struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Microsoft implement OAuth authorization for microsoft`s services
	Microsoft struct {
		cfg        *MicrosoftConfig
		httpClient *resty.Client
	}
)

// GetAuthURL generates the authentication URL for Microsoft OAuth.
// It takes in the context, redirectUrl, and state as parameters.
// It sets the RedirectURL for the MicrosoftConfig and returns
func (m *Microsoft) GetAuthURL(ctx context.Context, redirectUrl, state string) (string, error) {
	m.cfg.RedirectUL = fmt.Sprintf("%s", redirectUrl)
	return fmt.Sprintf(microsoftLoginURL, m.cfg.ClientID, microsoftScope+" "+teamsScope, m.cfg.RedirectUL), nil
}

// ExchangeCode exchanges the authorization code for an access token
// and returns the related IdentityResponse.
// It takes in the context and code as parameters.
// It constructs the payload with the necessary parameters, sends a POST request
// with the payload to the Microsoft token endpoint,
// and handles the response, parsing it into a tokenResponse.
// It then constructs an IdentityResponse with the access token, token type, refresh token,
// and expiry, and returns it.
// If any errors occur during the exchange, it returns a nil IdentityResponse and the error.
// This method requires a valid httpClient for making the API call.
func (m *Microsoft) ExchangeCode(ctx context.Context, code string) (*IdentityResponse, error) {

	payload := map[string]string{
		"client_id":     m.cfg.ClientID,
		"client_secret": m.cfg.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  m.cfg.RedirectUL,
	}
	var response tokenResponse
	resp, err := m.httpClient.R().SetFormData(payload).
		Post(microsoftToken)
	if err = utils.WrapRestyError(resp, err); err != nil {
		return nil, utils.ErrorPermission.Newf("exchange code error: %s", err.Error())
	}
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}
	return &IdentityResponse{
		Token: &oauth2.Token{
			AccessToken:  response.AccessToken,
			TokenType:    response.TokenType,
			RefreshToken: response.RefreshToken,
			Expiry:       time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
		},
	}, nil
}

// RefreshToken refreshes the OAuth2 token using the provided refresh token.
// It takes in a token *oauth2.Token and uses the client ID, client secret, refresh token,
// grant type, and redirect URI to construct the request payload.
// The payload is sent as a POST request to the Microsoft token endpoint.
// The response is parsed into a tokenResponse struct, and a new OAuth2 token is constructed
// using the response's access token, token type, refresh token, and expiry.
// The new token is then returned, along with any error that occurred during the process.
func (m *Microsoft) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	payload := map[string]string{
		"client_id":     m.cfg.ClientID,
		"client_secret": m.cfg.ClientSecret,
		"refresh_token": token.RefreshToken,
		"grant_type":    "refresh_token",
		"redirect_uri":  m.cfg.RedirectUL,
	}
	var response tokenResponse
	resp, err := m.httpClient.R().SetFormData(payload).
		Post(microsoftToken)
	if err = utils.WrapRestyError(resp, err); err != nil {
		return nil, utils.ErrorPermission.Newf("exchange code error: %s", err.Error())
	}
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken:  response.AccessToken,
		TokenType:    response.TokenType,
		RefreshToken: response.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
	}, nil
}

// NewMicrosoft creates a new instance of the Microsoft OAuth provider.
// It takes in a MicrosoftConfig object as a parameter and returns a Proxy interface.
// The returned object can be used to interact with the Microsoft OAuth service.
func NewMicrosoft(cfg *MicrosoftConfig) Proxy {
	return &Microsoft{
		cfg:        cfg,
		httpClient: resty.New().SetTimeout(time.Minute),
	}
}
