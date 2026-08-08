package services

import (
	"context"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
)

type AwsService struct {
	mediaClient mediav1.MediaServiceClient
}

func (s *AwsService) ConfirmObjectUpload(ctx context.Context, req *api.S3UploadEventRequest) (error) {
	_, err := s.mediaClient.ObjectUploadConfirmation(ctx, &mediav1.ObjectUploadConfirmationRequest{
		FileSize: req.FileSize,
		ObjectId: req.ObjectId,
	})
	return err
}

func NewAwsService(mediaClient mediav1.MediaServiceClient) *AwsService {
	return &AwsService{
		mediaClient: mediaClient,
	}
}