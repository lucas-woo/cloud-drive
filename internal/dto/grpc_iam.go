package dto

import (
	"time"

	"github.com/google/uuid"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
)

type GenerateNewApiKeyRequest struct {
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

type AddUserRolePermissionRequest struct {
	UserRole iamv1.AddUserRolePermissionRequest_Permission
	UserId string
	ProjectId string
}

type ValidatedUserPermissionRequest struct {
	UserRole iamv1.ValidateUserPermissionRequest_Permission
	UserId string
	ProjectId string
}

type ApiKey struct {
	ApiKey string
	Name string
	CreatedAt time.Time
	IsActive bool
}

type GetProjectIdRequest struct {
	ApiKey string
	ApiSecret string
}