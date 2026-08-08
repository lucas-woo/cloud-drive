package services

import mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"

type AwsService struct {
	mediaClient mediav1.MediaServiceClient
}

func NewAwsService(mediaClient mediav1.MediaServiceClient) *AwsService {
	return &AwsService{
		mediaClient: mediaClient,
	}
}