package api

type S3UploadEventRequest struct {
	ObjectId     string `json:"objectId" binding:"required"`
	FileSize     uint64 `json:"fileSize" binding:"gte=0"`
}