package api

type UploadObjectApiRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`	
	OriginalFileName string `json:"originalFileName" binding:"required"`
  FolderId  string `json:"folderId" binding:"required"`
	IsActive bool `json:"isActive" binding:"required"`	
}

type UploadObjectApiResponse struct {
	ObjectId string `json:"objectId"`
}

type Transformations struct {
  
}

type UploadImageApiRequst struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`	
  FolderId  string `json:"folderId" binding:"required"`
	OriginalFileName string `json:"originalFileName" binding:"required"`
	IsActive bool `json:"isActive" binding:"required"`	
  Transfromations Transformations `json:"transformations" binding:"required"`	
}

type UploadImageApiResponse struct {

}