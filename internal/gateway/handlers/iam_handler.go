package handlers

import "github.com/lucas-woo/cloud-drive/internal/gateway/services"

type IamHandler struct {
	service *services.IamService
}

func NewIamHandler(iamService *services.IamService) *IamHandler {
	return &IamHandler{
		service: iamService,
	}
}