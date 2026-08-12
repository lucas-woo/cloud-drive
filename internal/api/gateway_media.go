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
	ProjectId    string    `json:"projectId"`
	CollectionId string    `json:"collectionId"`
	FolderId    string    `json:"folderId"`
	ObjectId     string    `json:"objectId"`

	FileSize int64  `json:"fileSize"`
	Format   string `json:"format"`

	IsActive bool `json:"isActive"`

	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type Folder struct {
	FolderId   string    `json:"folderId"`
	FolderName string    `json:"folderName"`

	FolderSize int64 `json:"folderSize"`
	AssetCount int32 `json:"assetCount"`

	LastUpload time.Time `json:"lastUpload"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type Collection struct {
	CollectionId string `json:"collectionId"`
	Name         string `json:"name"`
	Description  string `json:"description"`

	CreatedAt   time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`

	IsPublic bool `json:"isPublic"`
}

type AssetCursor struct {
  CreatedAt time.Time `json:"createdAt"`
	ObjectId string `json:"objectId"`
}

type GetAssetsPageResponse struct {
  AssetCursor *AssetCursor `json:"assetCursor"`
  ProjectObjects []ProjectObject `json:"projectObjects"`
  ProjectFolders []Folder `json:"folders"`
  ProjectCollections []Collection `json:"collections"`
}


type GetFolderAssetsRequest struct {
    AssetCursor *AssetCursor `json:"assetCursor"`

    ProjectId string `json:"projectId" binding:"required"`
    FolderId  string `json:"folderId" binding:"required"`
}

type GetFolderAssetsResponse struct {
  AssetCursor *AssetCursor `json:"assetCursor"`
  ProjectObjects []ProjectObject `json:"projectObjects"`
}

type GetCollectionAssetsRequest struct {
    AssetCursor *AssetCursor `json:"assetCursor"`

    ProjectId string `json:"projectId" binding:"required"`
    CollectionId  string `json:"collectionId" binding:"required"`
}

type GetCollectionAssetsResponse struct {
  AssetCursor *AssetCursor `json:"assetCursor"`
  ProjectObjects []ProjectObject `json:"projectObjects"`
}

type CreateCollectionRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`
  Description string `json:"description" binding:"required"`
}

type CreateCollectionResponse struct {
	CollectionId string `json:"collectionId"`
}

type CreateFolderRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`
}

type CreateFolderResponse struct {
	FolderId string `json:"folderId"`
}

type CreateNewProjectRequest struct {
  ProjectName string `json:"projectName" binding:"required"`
  Description string `json:"description" binding:"required"`	
}

type CreateNewProjectResponse struct {
	ProjectId string `json:"projectId"`
}

type UploadObjectRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  Name string `json:"name" binding:"required"`	
  FolderId  string `json:"folderId" binding:"required"`	
	IsActive bool `json:"isActive" binding:"required"`	
}

type UploadObjectResponse struct {
	Url string `json:"url"`
	ObjectId string `json:"objectId"`
}