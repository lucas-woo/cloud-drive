package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)



func InitializeRouter(webGroup *gin.RouterGroup, middlewares *middlewares.AuthMiddleware, authHandler *handlers.AuthHandler, mediaHandler *handlers.MediaHandler) {
	InitializeAuthWebRoutes(webGroup, middlewares, authHandler)
}