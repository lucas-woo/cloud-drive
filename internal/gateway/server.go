package gateway

import (
	"fmt"
	"log"
	"net/http"
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
		fmt.Println("no port found")
		port = "3000"
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
	webGroup.Use(RequireEnvTestCookie())

	apiGroup := ginEngine.Group("/api")
	apiGroup.Use(publicCors)

	awsGroup := ginEngine.Group("/webhooks/aws")

	healthGroup := ginEngine.Group("/health")
	
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
	iamService := services.NewIamService(resources.IamClient, utils.RestMapper{})
	iamHandler := handlers.NewIamHandler(iamService)

	//event bridge route
	awsService := services.NewAwsService(resources.MediaClient)
	awsHandler := handlers.NewAwsHandler(awsService)

	//health route
	healthHandler := handlers.NewHealthHandler()

	routes.InitializeRouter(webGroup, awsGroup, healthGroup, authMiddlewares, eventbridgeMiddlewares, authHandler, mediaHandler, iamHandler, awsHandler, healthHandler)

	return server
}



func RequireEnvTestCookie() gin.HandlerFunc {
	expectedCookie := os.Getenv("TEST_COOKIE")
	fmt.Println("cookie password:", expectedCookie)
	
	if expectedCookie == "" {
		println("no test cookie")
	}

	return func(c *gin.Context) {

		cookieVal, err := c.Cookie("TEST_COOKIE")

		fmt.Println("HERE!", cookieVal, cookieVal == expectedCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			return
		}
		if cookieVal == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			return			
		}

		c.Next()
	}
}