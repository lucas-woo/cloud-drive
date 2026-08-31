package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)

func InitializeIamWebRoutes(webGroup *gin.RouterGroup, middlewares *middlewares.AuthMiddleware, iamHandler *handlers.IamHandler) {

	//needs validation middleware

	// add email validation before generating later?
	webGroup.POST("/generate-api-key", middlewares.IsAuthenticated(), iamHandler.GenerateNewApiKey)
	
	webGroup.POST("/api-keys-page", middlewares.IsAuthenticated(), iamHandler.GetAllApiKeys)
	// TODO:
	// get api keys
}