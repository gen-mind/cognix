package handler

import (
	"cognix.ch/api/v2/core/configmap"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigMapHandler struct {
	configMapClientBuilder *configmap.ClientBuilder
}

func NewConfigMapHandler(configMapClientBuilder *configmap.ClientBuilder) *ConfigMapHandler {
	return &ConfigMapHandler{configMapClientBuilder: configMapClientBuilder}
}

// Mount sets up the config maps routes by adding them to the specified gin.Engine instance.
func (h *ConfigMapHandler) Mount(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := router.Group("/api/config-map").Use(authMiddleware)
	handler.GET("/:name", server.HandlerErrorFuncAuth(h.GetConfigMap))
	handler.POST("/:name", server.HandlerErrorFuncAuth(h.Save))
	handler.DELETE("/:name/:key", server.HandlerErrorFuncAuth(h.Delete))

}

// GetConfigMap returns list of values from given config map name
// @Summary returns list of values from given config map name
// @Description returns list of values from given config map name
// @Tags ConfigMap
// @ID configmap_get_config_map
// @Produce  json
// @Param name path string true "name of config map"
// @Security ApiKeyAuth
// @Success 200 {array} proto.ConfigMapRecord
// @Router /api/config-map/{name} [get]
func (h *ConfigMapHandler) GetConfigMap(c *gin.Context, identity *security.Identity) error {
	if !identity.User.HasRoles(model.RoleSuperAdmin) {
		return utils.ErrorPermission.New("do not have super admin permission")
	}
	configMapName := c.Param("name")
	if configMapName == "" {
		return utils.ErrorBadRequest.New("config map name should be presented")
	}
	client, err := h.configMapClientBuilder.Client()
	if err != nil {
		return utils.Internal.Wrap(err, err.Error())
	}
	result, err := client.GetList(c.Request.Context(), &proto.ConfigMapList{Name: configMapName})
	if err != nil {
		return utils.Internal.Wrap(err, err.Error())
	}

	return server.JsonResult(c, http.StatusOK, result.GetValues())
}

// Save create or update value in config map
// @Summary create or update value in config map
// @Description create or update value in config map
// @Tags ConfigMap
// @ID configmap_save_config_map
// @Produce  json
// @Param name path string true "name of config map"
// @Param params body parameters.ConfigMapValue true "payload"
// @Security ApiKeyAuth
// @Success 200 {object} server.JsonResponse
// @Router /api/config-map/{name} [post]
func (h *ConfigMapHandler) Save(c *gin.Context, identity *security.Identity) error {
	if !identity.User.HasRoles(model.RoleSuperAdmin) {
		return utils.ErrorPermission.New("do not have super admin permission")
	}
	configMapName := c.Param("name")
	if configMapName == "" {
		return utils.ErrorBadRequest.New("config map name should be presented")
	}
	var param parameters.ConfigMapValue
	if err := c.BindJSON(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	if err := param.Validate(); err != nil {
		return utils.ErrorBadRequest.Wrap(err, err.Error())
	}

	client, err := h.configMapClientBuilder.Client()
	if err != nil {
		return utils.Internal.Wrap(err, err.Error())
	}
	if _, err = client.Save(c.Request.Context(), &proto.ConfigMapSave{
		Name: configMapName,
		Value: &proto.ConfigMapRecord{
			Key:   param.Key,
			Value: param.Value,
		},
	}); err != nil {
		return utils.Internal.Wrap(err, err.Error())
	}
	return server.JsonResult(c, http.StatusOK, "ok")
}

// Delete deletes value from config map
// @Summary deletes value from config map
// @Description deletes value from config map
// @Tags ConfigMap
// @ID configmap_delete_config_map
// @Produce  json
// @Param name path string true "name of config map"
// @Param key path string true "key of config map"
// @Security ApiKeyAuth
// @Success 200 {object} server.JsonResponse
// @Router /api/config-map/{name}/{key} [delete]
func (h *ConfigMapHandler) Delete(c *gin.Context, identity *security.Identity) error {
	if !identity.User.HasRoles(model.RoleSuperAdmin) {
		return utils.ErrorPermission.New("do not have super admin permission")
	}
	configMapName := c.Param("name")
	if configMapName == "" {
		return utils.ErrorBadRequest.New("config map name should be presented")
	}

	configMapKey := c.Param("key")
	if configMapKey == "" {
		return utils.ErrorBadRequest.New("config map key should be presented")
	}

	client, err := h.configMapClientBuilder.Client()
	if err != nil {
		return utils.Internal.Wrap(err, err.Error())
	}
	if _, err = client.Delete(c.Request.Context(), &proto.ConfigMapDelete{Name: configMapName, Key: configMapKey}); err != nil {
		return utils.Internal.Wrap(err, err.Error())
	}
	return server.JsonResult(c, http.StatusOK, "ok")
}
