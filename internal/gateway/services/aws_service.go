package services

import (
	"context"
	"log"
	"strings"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
)

type AwsService struct {
	mediaClient mediav1.MediaServiceClient
}

func (s *AwsService) ConfirmObjectUpload(ctx context.Context, req *api.S3UploadEventRequest) error {
    objectPath := req.ObjectPath

    index := strings.LastIndex(objectPath, "/")
    if index == -1 || index == len(objectPath)-1 {
        log.Printf("invalid object path in object confirmation: %s", objectPath)
				return nil
    }

    objectID := objectPath[index+1:]

    _, err := s.mediaClient.ObjectUploadConfirmation(
        ctx,
        &mediav1.ObjectUploadConfirmationRequest{
            FileSize: req.FileSize,
            ObjectId: objectID,
        },
    )

    return err
}

func NewAwsService(mediaClient mediav1.MediaServiceClient) *AwsService {
	return &AwsService{
		mediaClient: mediaClient,
	}
}