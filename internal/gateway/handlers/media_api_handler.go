package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type MediaApiHandler struct {
	service *services.MediaApiService
}

func (h *MediaApiHandler) UploadFileApi(c *gin.Context) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid multipart request",
		})
		return
	}

	apiKey := c.GetHeader("API-Key")
	apiSecret := c.GetHeader("API-Secret")

	if apiKey == "" || apiSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing API credentials",
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

	var req api.UploadObjectApiMetadataRequest
	if err := json.NewDecoder(metadataPart).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid metadata",
		})
		return
	}

	if req.ProjectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "projectId is required",
		})
		return
	}

	ok, err := h.service.ValidateApiKey(
		c.Request.Context(),
		apiKey,
		apiSecret,
		req.ProjectId,
		iamv1.Permission_PERMISSION_UPLOAD,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "permission denied",
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

	res, err := h.service.UploadFile(
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

func (h *MediaApiHandler) UploadImageApi(c *gin.Context) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid multipart request",
		})
		return
	}

	apiKey := c.GetHeader("API-Key")
	apiSecret := c.GetHeader("API-Secret")

	if apiKey == "" || apiSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing API credentials",
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

	var req api.UploadImageApiMetadataRequst
	if err := json.NewDecoder(metadataPart).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid metadata",
		})
		return
	}

	if req.ProjectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "projectId is required",
		})
		return
	}

	ok, err := h.service.ValidateApiKey(
		c.Request.Context(),
		apiKey,
		apiSecret,
		req.ProjectId,
		iamv1.Permission_PERMISSION_UPLOAD,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "permission denied",
		})
		return
	}

	imagePart, err := reader.NextPart()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image is required",
		})
		return
	}
	defer imagePart.Close()

	if imagePart.FormName() != "image" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image must be the second multipart part",
		})
		return
	}

	contentType := imagePart.Header.Get("Content-Type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image content type is required",
		})
		return
	}

	res, err := h.service.UploadImage(
		c.Request.Context(),
		req.ProjectId,
		req.Name,
		req.OriginalFileName,
		req.FolderId,
		req.IsActive,
		contentType,
		&req.Transformations,
		imagePart,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *MediaApiHandler) GetProjectId(c *gin.Context) {
	apiKey := c.GetHeader("API-Key")

	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing api key in header",
		})
		return
	}	

	projectId, err := h.service.GetProjectId(c.Request.Context(), apiKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid api key",
		})
		return
	}

	c.JSON(http.StatusOK, api.GetProjectIdApiResponse{
		ProjectId: projectId,
	})
}

func (h *MediaApiHandler) GetAllFolders(c *gin.Context) {
	apiKey := c.GetHeader("API-Key")
	apiSecret := c.GetHeader("API-Secret")

	if apiKey == "" || apiSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing API credentials",
		})
		return
	}	

	var req api.GetAllFoldersApiRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing projectId",
		})
		return
	}	

	ok, err := h.service.ValidateApiKey(c.Request.Context(), apiKey, apiSecret, req.ProjectId, iamv1.Permission_PERMISSION_GET)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "permission denied",
		})
		return		
	}	
	
	res, err := h.service.GetAllFolders(c.Request.Context(), req.ProjectId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error getting all folders",
		})
		return
	}
	
	c.JSON(http.StatusOK, res)
}

func (h *MediaApiHandler) CreateNewFolder(c *gin.Context) {
	apiKey := c.GetHeader("API-Key")
	apiSecret := c.GetHeader("API-Secret")

	if apiKey == "" || apiSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing API credentials",
		})
		return
	}	

	var req api.CreateFolderApiRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing projectId",
		})
		return
	}	

	ok, err := h.service.ValidateApiKey(c.Request.Context(), apiKey, apiSecret, req.ProjectId, iamv1.Permission_PERMISSION_CREATE_FOLDER)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "permission denied",
		})
		return		
	}

	res, err := h.service.CreateFolder(c.Request.Context(), req.ProjectId, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to create new folder",
		})
		return		
	}
	c.JSON(http.StatusOK, res)
}

func NewMediaApiHandler(service *services.MediaApiService) *MediaApiHandler{
	return &MediaApiHandler{
		service: service,
	}
}