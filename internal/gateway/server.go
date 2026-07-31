package gateway

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	webCors := cors.New(cors.Config{
		AllowOrigins: []string{"https://yourwebsite.com", "https://yourwebsite.com"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 5 * time.Minute,
	})

	publicCors := cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "X-API-Key", "X-API-Secret"},
	})	

	webGroup := ginEngine.Group("/app")
	webGroup.Use(webCors)

	apiGroup := ginEngine.Group("/api")
	apiGroup.Use(publicCors)

	server := &Server{
		gin: ginEngine,
	}

	middlewares := middlewares.NewAuthMiddleware(resources.RedisRepo)

	authService := services.NewAuthServer(resources.AuthClient, resources.MediaClient)

	authHandler := handlers.NewAuthHandler(authService)

	routes.InitializeRouter(webGroup, middlewares, authHandler)

	return server
}