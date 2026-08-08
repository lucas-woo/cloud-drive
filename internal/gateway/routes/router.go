package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)



func InitializeRouter(webGroup *gin.RouterGroup, awsGroup *gin.RouterGroup, middlewares *middlewares.AuthMiddleware, authHandler *handlers.AuthHandler, mediaHandler *handlers.MediaHandler, iamHandler *handlers.IamHandler, awsHandler *handlers.AwsHandler) {

	InitializeAuthWebRoutes(webGroup, middlewares, authHandler)

	InitializeMediaWebRoutes(webGroup, middlewares, mediaHandler)

	InitializeAwsRoutes(awsGroup, middlewares, awsHandler)
}