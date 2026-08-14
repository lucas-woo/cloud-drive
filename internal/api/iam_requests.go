package api

import "time"

type GenerateApiKeyRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  KeyName string `json:"keyName" binding:"required"`	
}

type ApiKeysPageRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type GenerateApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`		
	CreatedAt time.Time `json:"CreatedAt"`		
}


