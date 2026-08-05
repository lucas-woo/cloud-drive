package dto

import (
	"time"

	"github.com/google/uuid"
)


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

type ProjectObject struct {
	ProjectId    uuid.UUID
	CollectionId uuid.UUID
	FolderId     uuid.UUID
	ObjectId     uuid.UUID

	FileSize int64
	Format   string

	IsActive bool

	CreatedAt  time.Time
	ModifiedAt time.Time
}

type AssetCursor struct {
    CreatedAt time.Time
    ObjectId  uuid.UUID
}

type ProjectFolder struct {
  FolderId uuid.UUID

  FolderName string

	FolderSize int64
  AssetCount int32 

  LastUpload time.Time
  CreatedAt time.Time
  ModifiedAt time.Time
}

type GetAssetsRequest struct {
	ProjectId string
	AssetCursor *AssetCursorRequest
}

type AssetCursorRequest struct {
	CreatedAt time.Time
	ObjectId string
}

type GetAssetsInFolderRequest struct {
	ProjectId string
	FolderId string 
	AssetCursor *AssetCursorRequest
}

type GetAssetsInCollectionRequest struct {
	ProjectId string
	CollectionId string 
	AssetCursor *AssetCursorRequest	
}

type CreateNewFolderRequest struct {
	ProjectId string
	FolderName string
}

type CreateNewFolderResponse struct {
	FolderId uuid.UUID
}