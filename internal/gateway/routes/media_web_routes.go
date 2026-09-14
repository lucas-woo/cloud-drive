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
	// gets all folders and collections as well
	webGroup.POST("/assets-page", middlewares.IsAuthenticated(), mediaHandler.GetAssetsPage)


	//returns all folders and for each: how many assets, total size, last upload, location
	webGroup.POST("/folder-assets", middlewares.IsAuthenticated(), mediaHandler.GetFolderAssets) //done

	webGroup.POST("/collection-assets", middlewares.IsAuthenticated(), mediaHandler.GetCollectionAssets) // done

	webGroup.POST("/create-project", middlewares.IsAuthenticated(), mediaHandler.CreateNewProject)

	webGroup.POST("/create-collection", middlewares.IsAuthenticated(), mediaHandler.CreateNewCollection) //done

	webGroup.POST("/create-folder", middlewares.IsAuthenticated(), mediaHandler.CreateNewFolder) // done

	webGroup.POST("/create-object", middlewares.IsAuthenticated(), mediaHandler.CreateObject) 


	// TODO 
	// GetProjects
	// GetDashboard asks for project id

	// proto todo:
	webGroup.POST("/add-assets-to-collection")
	webGroup.POST("/add-assets-to-folder")

	webGroup.POST("update-asset-access-control")
	//update isPublic, create a new url link instead of using collectionId? 
	webGroup.POST("collection-link")
}