package gateway

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
	"github.com/lucas-woo/cloud-drive/internal/gateway/routes"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)


type Server struct {
	gin *gin.Engine
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

	ginEngine := gin.Default()

	server := &Server{
		gin: ginEngine,
	}

	middlewares := middlewares.NewAuthMiddleware(resources.RedisRepo)

	authService := services.NewAuthServer(resources.AuthClient)

	authHandler := handlers.NewAuthHandler(authService)

	routes.InitializeRouter(ginEngine, middlewares, authHandler)

	return server
}