package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
)

func InitializeMediaApiRoutes(apiGroup *gin.RouterGroup, mediaHandler *handlers.MediaApiHandler) {
	apiGroup.POST("/upload-object", mediaHandler.UploadFileApi)
	apiGroup.POST("/upload-image", mediaHandler.UploadImageApi)

	apiGroup.GET("/find-project-id", )
	apiGroup.GET("/all-folders")
	apiGroup.POST("/create-folder")
}