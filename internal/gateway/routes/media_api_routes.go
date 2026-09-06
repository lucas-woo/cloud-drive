package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
)

func InitializeMediaApiRoutes(apiGroup *gin.RouterGroup, mediaHandler *handlers.MediaHandler) {
	apiGroup.POST("/upload-object", mediaHandler.UploadObjectApi)
	apiGroup.POST("/upload-image")
}