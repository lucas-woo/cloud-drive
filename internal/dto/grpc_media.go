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
	FolderId string
	IsActive bool
}

type ObjectUploadConfirmationRequest struct {
	ObjectId string
	FileSize uint64
	Format string
	ErrorStatus error
}

type GetDashboardRequest struct {
	UserId string
}

type 