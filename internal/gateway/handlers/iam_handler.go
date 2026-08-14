package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type IamHandler struct {
	service *services.IamService
}

func (h *IamHandler) GenerateNewApiKey(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.GenerateApiKeyRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	projectId := reqBody.ProjectId
	keyName := reqBody.KeyName

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}		

	res, err := h.service.GenerateNewApiKey(c.Request.Context(), projectId, keyName)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return 
	}

	c.JSON(http.StatusOK, res)	
}


func (h *IamHandler) GetAllApiKeys(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.ApiKeysPageRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	projectId := reqBody.ProjectId

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}			

	

}


func NewIamHandler(iamService *services.IamService) *IamHandler {
	return &IamHandler{
		service: iamService,
	}
}