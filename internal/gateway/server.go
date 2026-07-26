package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/database"
)


type Server struct {
	gin *gin.Engine
	resources *database.GatewayResources
}



func NewServer() *Server {
	return &Server{}
}