package routes

import "github.com/gin-gonic/gin"



func InitializeRouter(ginEngine *gin.Engine) {
	ginEngine.GET("/") 
}