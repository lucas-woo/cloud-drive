package api

type S3UploadEventRequest struct {
	ObjectPath     string `json:"objectPath" binding:"required"`
	FileSize     uint64 `json:"fileSize" binding:"gte=0"`
}