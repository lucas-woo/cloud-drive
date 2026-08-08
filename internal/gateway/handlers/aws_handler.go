package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type AwsHandler struct {
	service *services.AwsService
}

func (h *AwsHandler) ConfirmObjectUpload(c *gin.Context) {

	var reqBody api.S3UploadEventRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.Status(http.StatusOK)
		log.Printf("error binding json eventbridge request: %v", err)
		return
	}		
	if h.service.ConfirmObjectUpload(c.Request.Context(), &reqBody) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}	
	c.Status(http.StatusOK)
}

func NewAwsHandler (service *services.AwsService) *AwsHandler {
	return &AwsHandler{
		service: service,
	}
}