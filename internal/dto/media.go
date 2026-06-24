package dto


type CreateNewProjectRequest struct {
	UserId string
	ProjectName string
	Description string
}

type UploadObjectRequest struct {
	SessionId string
	ProjectId string
	ObjectName string
	Folder string
}

type LambdaS3UploadConfirmationRequest struct {
	ObjectId string
	FileSize uint64
	Format string
}