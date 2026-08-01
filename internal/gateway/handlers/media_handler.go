package handlers

import "github.com/lucas-woo/cloud-drive/internal/gateway/services"

type MediaHandler struct {
	service *services.MediaService
}

func NewMediaHandler(service *services.MediaService) *MediaHandler{
	return &MediaHandler{
		service: service,
	}
}