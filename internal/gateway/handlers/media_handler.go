package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type MediaHandler struct {
	service *services.MediaService
}


func (h *MediaHandler) GetDashboard(c *gin.Context) {

	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)

	dashboard, err := h.service.GetDashboard(c.Request.Context(), userId)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dashboard)
}


func (h *MediaHandler) GetAssetsPage(c *gin.Context) {

	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.GetAssetsPageRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	projectId := reqBody.ProjectId

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	res, err := h.service.GetAssetsPage(c.Request.Context(), projectId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return		
	}
	
	c.JSON(http.StatusOK, res)
}

func (h *MediaHandler) GetFolderAssets(c *gin.Context) {


	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.GetFolderAssetsRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	projectId := reqBody.ProjectId
	folderId := reqBody.FolderId
	assetCursor := reqBody.AssetCursor

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}	
	
	res, err := h.service.GetAssetsInFolder(c.Request.Context(), projectId, folderId, assetCursor)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return		
	}
	
	c.JSON(http.StatusOK, res)	

}

func (h *MediaHandler) GetCollectionAssets(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.GetCollectionAssetsRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	projectId := reqBody.ProjectId
	collectionId := reqBody.CollectionId
	assetCursor := reqBody.AssetCursor

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}	
	
	res, err := h.service.GetAssetsInFolder(c.Request.Context(), projectId, collectionId, assetCursor)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return		
	}
	
	c.JSON(http.StatusOK, res)	

}

func (h *MediaHandler) CreateNewProject(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.CreateNewProjectRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	name := reqBody.ProjectName
	description := reqBody.Description

	projectId, err := h.service.CreateNewProject(c.Request.Context(), userId, name, description)
	if err != nil {
		//should either retry or delete the user
		log.Printf("error creating new user project: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	res := api.CreateNewProjectResponse{
		ProjectId: projectId,
	}

	c.JSON(http.StatusOK, res)		
}

func (h *MediaHandler) CreateNewCollection(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.CreateCollectionRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	projectId := reqBody.ProjectId
	name := reqBody.Name
	description := reqBody.Description

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}	

	res, err := h.service.CreateNewCollection(c.Request.Context(), projectId, userId, name, description)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return 
	}

	c.JSON(http.StatusOK, res)	
	
}

func (h *MediaHandler) CreateNewFolder(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.CreateFolderRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	projectId := reqBody.ProjectId
	name := reqBody.Name

	//needs better validation later
	if len(name) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return				
	}

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}	

	res, err := h.service.CreateNewFolder(c.Request.Context(), projectId, config.FolderPrefix + name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return 
	}

	c.JSON(http.StatusOK, res)	
	
}

func (h *MediaHandler) CreateObject(c *gin.Context) {
	userIdString, exists := c.Get(config.GinUserId)

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId := userIdString.(string)	

	var reqBody api.CreateObjectRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	

	projectId := reqBody.ProjectId

	authorized, err := h.service.ValidateUserRole(c.Request.Context(), userId, projectId, iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}		

	name := reqBody.Name
	folderId := reqBody.FolderId
	isActive := reqBody.IsActive
	originalFileName := reqBody.FileName

	res, err := h.service.GetUploadObjectUrl(c.Request.Context(), projectId, folderId, name, originalFileName, isActive)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return 
	}

	c.JSON(http.StatusOK, res)	
}


func (h *MediaHandler) UploadObjectApi(c *gin.Context) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid multipart request",
		})
		return
	}

	metadataPart, err := reader.NextPart()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "metadata is required",
		})
		return
	}
	defer metadataPart.Close()

	if metadataPart.FormName() != "metadata" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "metadata must be the first multipart part",
		})
		return
	}

	var req api.UploadObjectApiRequest
	if err := json.NewDecoder(metadataPart).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid metadata",
		})
		return
	}

	filePart, err := reader.NextPart()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}
	defer filePart.Close()

	if filePart.FormName() != "file" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file must be the second multipart part",
		})
		return
	}

	contentType := filePart.Header.Get("Content-Type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file content type is required",
		})
		return
	}

	res, err := h.service.UploadObject(
		c.Request.Context(),
		req.ProjectId,
		req.Name,
		req.OriginalFileName,
		req.FolderId,
		req.IsActive,
		contentType,
		filePart,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}


func (h *MediaHandler) UploadImageApi(c *gin.Context) {

}

func NewMediaHandler(service *services.MediaService) *MediaHandler{
	return &MediaHandler{
		service: service,
	}
}