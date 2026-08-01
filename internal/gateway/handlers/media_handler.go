package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type MediaHandler struct {
	service *services.MediaService
}


func (h *MediaHandler) GetDashboard(c *gin.Context) {

	
}



func NewMediaHandler(service *services.MediaService) *MediaHandler{
	return &MediaHandler{
		service: service,
	}
}