package gateway

import (
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
	"github.com/lucas-woo/cloud-drive/internal/utils"
)


type Server struct {
	gin *gin.Engine
}

func (s *Server) Run() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("listening on port %s", port)

	if err := s.gin.Run(addr); err != nil {
		log.Fatalf("error starting gin server: %v", err)
	}
}

func NewServer(resources *database.GatewayResources) *Server {

	ginEngine := gin.Default()

	webCors := cors.New(cors.Config{
		AllowOrigins: []string{"*"},//"http://localhost:3000"
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

	awsGroup := ginEngine.Group("/webhooks/aws")
	
	server := &Server{
		gin: ginEngine,
	}

	authMiddlewares := middlewares.NewAuthMiddleware(resources.RedisRepo)
	eventbridgeMiddlewares := middlewares.NewEventbridgeMiddleware()
	// auth routes
	authService := services.NewAuthServer(resources.AuthClient, resources.MediaClient, resources.IamClient)
	authHandler := handlers.NewAuthHandler(authService)

	//media routes
	mediaService := services.NewMediaService(resources.AuthClient, resources.MediaClient, resources.IamClient, utils.RestMapper{})
	mediaHandler := handlers.NewMediaHandler(mediaService)

	//iam routes
	iamService := services.NewIamService(resources.IamClient)
	iamHandler := handlers.NewIamHandler(iamService)

	//event bridge route
	awsService := services.NewAwsService(resources.MediaClient)
	awsHandler := handlers.NewAwsHandler(awsService)

	routes.InitializeRouter(webGroup, awsGroup, authMiddlewares, eventbridgeMiddlewares, authHandler, mediaHandler, iamHandler, awsHandler)

	return server
}