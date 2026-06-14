package apikeysmodels

type GenerateNewApiKeyRequest struct {
	SessionId string
	ProjectId string
	KeyName string
}
