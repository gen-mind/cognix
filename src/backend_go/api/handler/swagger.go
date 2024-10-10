package handler

import (
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"encoding/json"
	"net/http"

	_ "cognix.ch/api/v2/api/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

// SwaggerHandler represents the handler for Swagger related operations.
type SwaggerHandler struct {
}

// NewSwaggerHandler returns a new instance of SwaggerHandler.
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// Mount sets up Swagger routes and health check route.
func (h *SwaggerHandler) Mount(router *gin.Engine) {
	url := ginSwagger.URL("docs/doc.json")
	///url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
	router.GET("/swagger/docs", server.HandlerErrorFunc(h.GetDoc))
	router.GET("/api/health", h.Health)

}

// GetDoc returns the swagger doc in the json format.
func (h *SwaggerHandler) GetDoc(c *gin.Context) error {
	jsonDoc, err := swag.ReadDoc()
	if err != nil {
		return utils.NotFound.Wrap(err, "swagger doc not found")
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonDoc), &result)
	if err != nil {
		return utils.NotFound.Wrap(err, "can not parse swagger doc")
	}
	return server.JsonResult(c, http.StatusOK, result)
}

// Health returns a string response with status code 200 indicating that the server is functioning properly.
func (h *SwaggerHandler) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
