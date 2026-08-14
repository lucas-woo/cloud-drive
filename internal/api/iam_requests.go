package api

import "time"

type GenerateApiKeyRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
  KeyName string `json:"keyName" binding:"required"`	
}

type ApiKeysPageRequest struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type ApiKey struct {
	ApiKey string `json:"apiKey"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	IsActive bool `json:"isActive"`
}

type ApiKeysPageResponse struct {
	ApiKeys []*ApiKey `json:"apiKeys"`
}

type GenerateApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`		
	CreatedAt time.Time `json:"CreatedAt"`		
}


