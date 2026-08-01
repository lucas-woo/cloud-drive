package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)

func InitializeMediaWebRoutes(webGroup *gin.RouterGroup, middlewares *middlewares.AuthMiddleware, mediaHandler *handlers.MediaHandler) {
	//needs validation middleware
	webGroup.POST("")

}