package dto

import (
	"time"

	"github.com/google/uuid"
)

type GenerateNewApiKeyRequest struct {
	SessionId string
	ProjectId string
	KeyName string
}

type GenerateNewApiKeyResponse struct {
	ApiKey string
	ApiSecret string
	CreatedAt time.Time
	ApiId uuid.UUID
}