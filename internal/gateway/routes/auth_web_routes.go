package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)

func InitializeAuthWebRoutes(webGroup *gin.RouterGroup, middlewares *middlewares.AuthMiddleware, authHandler *handlers.AuthHandler) {
	//needs validation middleware
	webGroup.POST("/signup", middlewares.RedirectIfAuthenticated(), authHandler.SignUp)
	webGroup.POST("/login", middlewares.RedirectIfAuthenticated(), authHandler.Login)
	webGroup.POST("/logout", middlewares.CheckIfSessionExists(), authHandler.Logout)

}