package dto


type CreateNewProjectRequest struct {
	SessionId string
	ProjectName string
	Description string
}

type UploadObjectRequest struct {
	SessionId string
	ProjectId string
	ObjectName string
	Folder string
}
