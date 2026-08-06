package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type MediaService struct {
	authClient authv1.AuthServiceClient
	mediaClient mediav1.MediaServiceClient
}

func (s *MediaService) GetDashboard(ctx context.Context, userId string) (*dto.GatewayGetDashboardResponse, error) {
	
	res, err := s.mediaClient.GetDashboard(ctx, &mediav1.GetDashboardRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return &dto.GatewayGetDashboardResponse{
		ProjectId: res.GetProjectId(),
		ProjectName: res.GetProjectName(),
		Description: res.GetDescription(),
		AssetsAmount: res.GetAssetsAmount(),
		CreatedAt: res.GetCreatedAt().AsTime(),
		StorageBytes: res.GetStorageBytes(),
		Transformations: res.GetTransformations(),
	}, err
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