package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)



func InitializeRouter(ginEngine *gin.Engine, middlewares *middlewares.AuthMiddleware, authHandler *handlers.AuthHandler) {
	InitializeAuthRoutes(ginEngine, middlewares, authHandler)
}