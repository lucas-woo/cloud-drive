package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/config"
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



func NewMediaHandler(service *services.MediaService) *MediaHandler{
	return &MediaHandler{
		service: service,
	}
}