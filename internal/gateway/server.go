package gateway

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/database"
)


type Server struct {
	gin *gin.Engine
	resources *database.GatewayResources
}

func (s *Server) Run() {
	
	port := os.Getenv("GATEWAY_PORT")
	pp := fmt.Sprintf(":%s", port)

	err := s.gin.Run(pp)
	if err != nil {
		log.Fatal("error starting gin server")
	}
	fmt.Println("listening to port: " + port)
}

func NewServer(resources *database.GatewayResources) *Server {

	server := &Server{
		gin: gin.Default(),
		resources: resources,
	}

	server.InitializeRouter()

	return server
}