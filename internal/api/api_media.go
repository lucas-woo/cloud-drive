package api

type UploadObjectApiRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`	
	OriginalFileName string `json:"originalFileName" binding:"required"`
  FolderId  string `json:"folderId" binding:"required"`
	IsActive bool `json:"isActive" binding:"required"`	
}

