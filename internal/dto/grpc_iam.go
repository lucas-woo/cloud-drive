package dto

import (
	"time"

	"github.com/google/uuid"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
)

type GenerateNewApiKeyRequest struct {
	UserId string
	ProjectId string
	KeyName string
}

type GenerateNewApiKeyResponse struct {
	ApiKey uuid.UUID
	ApiSecret string
	CreatedAt time.Time
}

type ValidateApiKeyPermissionRequest struct {
	ApiKey string
	ApiSecret string
	PermissionRequest string
}

type UpdateUserRoleRequest struct {
	UserRole iamv1.UpdateUserRoleRequest_Permission
	UserId string
	ProjectId string
}