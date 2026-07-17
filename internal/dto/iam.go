package dto

import (
	"time"

	"github.com/google/uuid"
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
