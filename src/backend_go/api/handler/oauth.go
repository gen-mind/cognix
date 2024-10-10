package handler

import (
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

// OAuthHandler  provide oauth authentication for  third part services
type OAuthHandler struct {
	oauthConfig *oauth.Config
}

// NewOAuthHandler creates a new OAuthHandler with the given oauthConfig.
func NewOAuthHandler(oauthConfig *oauth.Config) *OAuthHandler {
	return &OAuthHandler{
		oauthConfig: oauthConfig,
	}
}

// Mount sets up the routes for OAuth authentication on the provided Gin router.
// It creates a new route group "/api/oauth" and registers the following routes:
// - GET "/:provider/auth_url" to handle getting the authentication URL for the given provider
// - GET "/:provider/callback" to handle the callback after the authentication process is completed
// - POST "/:provider/refresh_token" to handle refreshing the authentication token for the given provider.
func (h *OAuthHandler) Mount(route *gin.Engine) {
	handler := route.Group("/api/oauth")
	handler.GET("/:provider/auth_url", server.HandlerErrorFunc(h.GetUrl))
	//handler.GET("/google/signup", server.HandlerErrorFunc(h.SignUp))
	handler.GET("/:provider/callback", server.HandlerErrorFunc(h.Callback))
	handler.POST("/:provider/refresh_token", server.HandlerErrorFunc(h.Refresh))
}

// GetUrl retrieves the authentication URL for the given provider.
// It binds the query parameters from the context to the LoginParam struct.
// Then, it creates a new OAuth client for the specified provider and calls GetAuthURL
// to get the authentication URL with the provided redirect URL.
// Finally, it returns the URL as a response with the status code 200.
func (h *OAuthHandler) GetUrl(c *gin.Context) error {
	provider := c.Param("provider")
	var param parameters.LoginParam
	if err := c.ShouldBindQuery(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong redirect url")
	}

	oauthClient, err := oauth.NewProvider(provider, h.oauthConfig, param.RedirectURL)
	if err != nil {
		return utils.Internal.Wrap(err, "unknown provider")
	}
	url, err := oauthClient.GetAuthURL(c.Request.Context(), fmt.Sprintf("%s/api/oauth/%s/callback", param.RedirectURL, provider), "")
	if err != nil {
		return err
	}
	return server.StringResult(c, http.StatusOK, []byte(url))
}

// Callback handles the callback after the authentication process is completed for the specified provider.
// It retrieves the provider from the context parameters and binds the query parameters from the context
// to a map[string]string. Then, it creates a new OAuth client for the specified provider and exchanges
// the received authorization code for an access token. Finally, it returns the result as a JSON response
// with the status code 200.
func (h *OAuthHandler) Callback(c *gin.Context) error {
	provider := c.Param("provider")
	query := make(map[string]string)
	if err := c.BindQuery(&query); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}

	oauthClient, err := oauth.NewProvider(provider, h.oauthConfig, query["redirect_url"])
	if err != nil {
		return utils.Internal.Wrap(err, "unknown provider")
	}
	result, err := oauthClient.ExchangeCode(c.Request.Context(), query["code"])
	if err != nil {
		return utils.ErrorPermission.New(err.Error())
	}
	return server.JsonResult(c, http.StatusOK, result)
}

// Refresh handles the refreshing of the authentication token for the specified provider.
// It retrieves the provider from the context parameters and binds the JSON payload from the context
// to an oauth2.Token struct. Then, it creates a new OAuth client for the specified provider and calls
// the RefreshToken method with the provided token. Finally, it returns the refreshed token as a JSON
// response with the status code 200.
func (h *OAuthHandler) Refresh(c *gin.Context) error {
	provider := c.Param("provider")
	var token oauth2.Token
	if err := c.BindJSON(&token); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	oauthClient, err := oauth.NewProvider(provider, h.oauthConfig, "")
	if err != nil {
		return utils.Internal.Wrap(err, "unknown provider")
	}

	result, err := oauthClient.RefreshToken(&token)
	if err != nil {
		return utils.ErrorPermission.New(err.Error())
	}
	_ = provider
	return server.JsonResult(c, http.StatusOK, result)
}
