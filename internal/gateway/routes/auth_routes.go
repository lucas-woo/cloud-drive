package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)

func InitializeAuthRoutes(ginEngine *gin.Engine, middlewares *middlewares.AuthMiddleware, authHandler *handlers.AuthHandler) {
	//needs validation middleware
	ginEngine.POST("/signup", middlewares.RedirectIfAuthenticated(), authHandler.SignUp)
	ginEngine.POST("/login", middlewares.RedirectIfAuthenticated(), authHandler.Login)
	ginEngine.POST("/logout", )

}