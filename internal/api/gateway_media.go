package api

import "time"

type GetDashboardResponse struct {
	ProjectName     string    `json:"projectName"`
	ProjectId       string    `json:"projectId"`
	AssetsAmount    int32     `json:"assetsAmount"`
	Transformations int32     `json:"transformations"`
	StorageBytes    int64     `json:"storageBytes"`
	Description     string    `json:"description,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type GetAssetsPageRequest struct {
	ProjectId       string    `json:"projectId" binding:"required"`
}

type ProjectObject struct {
	ProjectId    string    `json:"project_id"`
	CollectionId string    `json:"collection_id"`
	FolderId    string    `json:"folder_id"`
	ObjectId     string    `json:"object_id"`

	FileSize int64  `json:"file_size"`
	Format   string `json:"format"`

	IsActive bool `json:"is_active"`

	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Folder struct {
	FolderId   string    `json:"folder_id"`
	FolderName string    `json:"folder_name"`

	FolderSize int64 `json:"folder_size"`
	AssetCount int32 `json:"asset_count"`

	LastUpload time.Time `json:"last_upload"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Collection struct {
	CollectionId string `json:"collection_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`

	CreatedAt   time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`

	IsPublic bool `json:"is_public"`
}

type AssetCursor struct {
  CreatedAt time.Time `json:"created_at"`
	ObjectId string `json:"object_id"`
}

