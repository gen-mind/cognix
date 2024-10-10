package server

import (
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-contrib/cors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlerFunc represents a function that takes a *gin.Context as input and returns an error.
type HandlerFunc func(c *gin.Context) error

// HandleFuncAuth represents a function type that handles authentication in a web server.
// It takes a Gin context and an identity object as parameters, and returns an error.
//
// Parameters:
// - c: a pointer to a Gin context
// - identity: a pointer to a security.Identity object that represents the user's identity
//
// Returns:
// - error: an error if the authentication fails, nil otherwise
type HandleFuncAuth func(c *gin.Context, identity *security.Identity) error

// JsonErrorResponse represents a JSON response format for reporting errors.
//
// It has the following fields:
// - Status: An integer indicating the HTTP status code.
// - Error: A string describing the error.
// - OriginalError: A string providing additional information about the error.
type JsonErrorResponse struct {
	Status        int    `json:"status,omitempty"`
	Error         string `json:"error,omitempty"`
	OriginalError string `json:"original_error,omitempty"`
}

// JsonResponse represents a JSON response that can be sent by an API.
type JsonResponse struct {
	Status int         `json:"status,omitempty"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// HandlerErrorFuncAuth is a higher-order function that creates a Gin middleware for handling errors during authenticated requests.
// It takes a function f of type HandleFuncAuth, which accepts a Gin context and a *security.Identity as parameters and returns an error.
// The middleware extracts the identity from the context using the GetContextIdentity function, and if it exists, calls the provided function f with the context and identity.
// It then calls the handleError function to handle any errors returned by f.
//
// Parameters:
// - f: a function of type HandleFuncAuth that accepts a Gin context and a *security.Identity and returns an error.
//
// Returns:
// - gin.HandlerFunc: a Gin middleware function
//
// Usage Example:
//
//	router := gin.Default()
//	authMiddleware := CustomAuthMiddleware()
//	router.GET("/api/endpoint", HandlerErrorFuncAuth(HandleFunc))
//
// Note: Replace HandleFunc with the actual function to be called within the middleware.
func HandlerErrorFuncAuth(f HandleFuncAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, err := GetContextIdentity(c)
		if err == nil {
			err = f(c, identity)
		}
		handleError(c, err)
	}
}

// HandlerErrorFunc wraps a HandlerFunc with error handling.
// It takes a HandlerFunc as input and returns a gin.HandlerFunc.
// The returned gin.HandlerFunc calls the input HandlerFunc and passes its result to handleError().
// If an error occurs, handleError() handles the error and sends an appropriate JSON response.
func HandlerErrorFunc(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handleError(c, f(c))
	}
}

// handleError takes a *gin.Context and an error as input. If the error is not nil, it checks if it is of type
// utils.Errors. If it is not, it creates a new utils.Errors object and sets its properties based on the error.
// Then, it creates a JsonErrorResponse struct and sets its properties based on the utils.Errors object.
// If the original error is not nil, it sets the OriginalError property of the JsonErrorResponse struct.
// Finally, it logs the error using zap.S().Errorf and returns a JSON response with the appropriate status code.
func handleError(c *gin.Context, err error) {
	if err != nil {
		ew, ok := err.(utils.Errors)
		if !ok {
			ew.Original = err
			ew.Code = http.StatusInternalServerError
			ew.Message = err.Error()
		}
		errResp := JsonErrorResponse{
			Status: int(ew.Code),
			Error:  ew.Message,
		}
		if ew.Original != nil {
			errResp.OriginalError = ew.Original.Error()
		}
		zap.S().Errorf("[%s] %v", ew.Message, ew.Original)
		c.JSON(int(ew.Code), errResp)
	}
}

// JsonResult returns a JSON response with the given status and data.
// The response format follows the JsonResponse struct.
// The status parameter indicates the HTTP status code of the response.
// The data parameter contains the data to be included in the response body.
func JsonResult(c *gin.Context, status int, data interface{}) error {
	c.JSON(status, JsonResponse{
		Status: status,
		Error:  "",
		Data:   data,
	})
	return nil
}

// StringResult sends the data as the response body with the specified status code.
//
// Parameters:
//   - c: The Gin context object.
//   - status: The status code to be returned in the response.
//   - data: The data to be sent as the response body.
//
// Returns:
//   - An error if the response cannot be sent.
func StringResult(c *gin.Context, status int, data []byte) error {
	c.Data(status, "", data)

	return nil
}

// BindJsonAndValidate extracts JSON data from the request body and validates it.
// It binds the JSON data to the provided 'data' interface. If 'data' implements the 'validation.Validatable' interface,
// it also calls the 'Validate' method on it and returns an error if validation fails.
// If an error occurs while binding or validating the JSON data, the error is returned.
//
// Parameters:
//   - c: gin.Context
//     The Gin context containing the request and response information.
//   - data: interface{}
//     The interface to bind the JSON data to.
//
// Returns:
//   - error
//     If an error occurs while binding or validating the JSON data, the error is returned. Otherwise, nil is returned.
func BindJsonAndValidate(c *gin.Context, data interface{}) error {
	if err := c.BindJSON(data); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	if vl, ok := data.(validation.Validatable); ok {
		if err := vl.Validate(); err != nil {
			return utils.ErrorBadRequest.New(err.Error())
		}
	}
	return nil
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware("service-name"))
	corsConfig := cors.DefaultConfig()

	corsConfig.CustomSchemas = cors.DefaultSchemas
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowWildcard = true
	router.Use(cors.New(corsConfig))
	return router
}
