package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)

func InitializeAwsRoutes(awsGroup *gin.RouterGroup, middlewares *middlewares.EventBridgeMiddleware, awsHandler *handlers.AwsHandler) {
	
	awsGroup.POST("/object-upload", middlewares.EventBridgeAuthMiddleware() ,awsHandler.ConfirmObjectUpload)

}