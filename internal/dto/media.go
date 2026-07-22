package dto


type CreateNewProjectRequest struct {
	UserId string
	ProjectName string
	Description string
}

type UploadObjectRequest struct {
	UserId string
	ProjectId string
	ObjectName string
	Folder string
	IsActive bool
}

type LambdaS3UploadConfirmationRequest struct {
	ObjectId string
	FileSize uint64
	Format string
	ErrorStatus error
}