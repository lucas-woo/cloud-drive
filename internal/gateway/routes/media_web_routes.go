package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/handlers"
	"github.com/lucas-woo/cloud-drive/internal/gateway/middlewares"
)

func InitializeMediaWebRoutes(webGroup *gin.RouterGroup, middlewares *middlewares.AuthMiddleware, mediaHandler *handlers.MediaHandler) {
	//needs validation middleware

	//should return: username, project name, project id, total images, total transformations, storage space; later on : credit usage
	webGroup.GET("/dashboard", middlewares.IsAuthenticated(), mediaHandler.GetDashboard)

	//first 40 assets, with info for each asset on folder location, format, file size, dimentions, placed by api/web (date), id, access control
	webGroup.GET("/assets-page")

	//returns all folders and for each: how many assets, total size, last upload, location
	webGroup.GET("/folders-page") //done
	webGroup.GET("/collections-page") // done

	//returns first 40 assets in that folder
	webGroup.GET("/folder") //done
	webGroup.GET("/collection") //done
	webGroup.POST("/upload-asset") // done?

	webGroup.POST("/create-collection") // done

	webGroup.POST("/create-folder") // done



	// proto todo:
	webGroup.POST("/add-assets-to-collection")

	webGroup.POST("update-asset-access-control")
	//update isPublic, create a new url link instead of using collectionId? 
	webGroup.POST("collection-link")
}