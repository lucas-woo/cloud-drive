package dto

import "time"

type GatewayGetDashboardResponse struct {
	ProjectName     string    `json:"projectName"`
	ProjectId       string    `json:"projectId"`
	AssetsAmount    int32     `json:"assetsAmount"`
	Transformations int32     `json:"transformations"`
	StorageBytes    int64     `json:"storageBytes"`
	Description     string    `json:"description,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type GatewayGetAssetsPageRequest struct {
	ProjectId       string    `json:"projectId" binding:"required"`
}

type GatewayProjectObjectResponse struct {

}

type GatewayGetAssetsPageResponse struct {

}