package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)



func InitializeRouter(
	webGroup *gin.RouterGroup, 
	apiGroup *gin.RouterGroup, 
	awsGroup *gin.RouterGroup, 
	healthGroup *gin.RouterGroup, 
	authMiddlewares *middlewares.AuthMiddleware, 
	eventbridgeMiddlewares *middlewares.EventBridgeMiddleware, 
	authHandler *handlers.AuthHandler, 
	mediaHandler *handlers.MediaHandler, 
	iamHandler *handlers.IamHandler, 
	awsHandler *handlers.AwsHandler, 
	healthHandler *handlers.HealthHandler, 
	mediaApiHandler *handlers.MediaApiHandler,
	) {

	InitializeAuthWebRoutes(webGroup, authMiddlewares, authHandler)

	InitializeMediaWebRoutes(webGroup, authMiddlewares, mediaHandler)

	InitializeIamWebRoutes(webGroup, authMiddlewares, iamHandler)

	InitializeAwsRoutes(awsGroup, eventbridgeMiddlewares, awsHandler)

	InitializeHealthRoutes(healthGroup, healthHandler)

	InitializeMediaApiRoutes(apiGroup, mediaApiHandler)
}