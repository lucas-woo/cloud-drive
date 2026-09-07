package api

type UploadObjectApiMetadataRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`	
	OriginalFileName string `json:"originalFileName" binding:"required"`
  FolderId  string `json:"folderId" binding:"required"`
	IsActive bool `json:"isActive" binding:"required"`	
}

type UploadObjectApiResponse struct {
	ObjectId string `json:"objectId"`
}

type ImageTransformations struct {
  Crop Crop `json:"crop"`
  Scale Scale `json:"scale"`
  Compression Compression `json:"compression"`
  Conversion Conversion `json:"conversion"`
}

type Crop struct {
  Width uint32 `json:"width" binding:"required"`
  Height uint32 `json:"height" binding:"required"`
}

type Scale struct {
  Width uint32 `json:"width" binding:"required"`
  Height uint32 `json:"height" binding:"required"`
}

type Compression struct {
  Compress bool `json:"compress" binding:"required"`
}

type Conversion struct {
  Format string `json:"format" binding:"required"`
}

type UploadImageApiMetadataRequst struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`	
  FolderId  string `json:"folderId" binding:"required"`
	OriginalFileName string `json:"originalFileName" binding:"required"`
	IsActive bool `json:"isActive" binding:"required"`	
  Transformations ImageTransformations `json:"transformations" binding:"required"`	
}

type UploadImageApiResponse struct {
  ObjectId string `json:"objectId"`
}

