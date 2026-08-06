package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto/gateway"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type MediaHandler struct {
	service *services.MediaService
}


func (h *MediaHandler) GetDashboard(c *gin.Context) {

	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)

	dashboard, err := h.service.GetDashboard(c.Request.Context(), userId)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// TODO: need to make a check if the userid has the role to get that project 
func (h *MediaHandler) GetAssetsPage(c *gin.Context) {

	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody dto.GatewayGetAssetsPageRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, reqBody.ProjectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	h.service.

}

func NewMediaHandler(service *services.MediaService) *MediaHandler{
	return &MediaHandler{
		service: service,
	}
}