package services

import (
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
)

type MediaService struct {
	authClient authv1.AuthServiceClient
	mediaClient mediav1.MediaServiceClient
}

func NewMediaService(	
	authClient authv1.AuthServiceClient, 
	mediaClient mediav1.MediaServiceClient,
	) *MediaService {

	return &MediaService{
		authClient: authClient,
		mediaClient: mediaClient,
	}
}