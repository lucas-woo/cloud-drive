package handlers

import "github.com/lucas-woo/cloud-drive/internal/gateway/services"

type AwsHandler struct {
	service *services.AwsService
}

func NewAwsHandler (service *services.AwsService) *AwsHandler {
	return &AwsHandler{
		service: service,
	}
}