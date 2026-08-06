package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type MediaService struct {
	authClient authv1.AuthServiceClient
	mediaClient mediav1.MediaServiceClient
	iamClient iamv1.IAMServiceClient
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

func (s *MediaService) ValidateUserRole(ctx context.Context, userId, projectId string, role iamv1.ValidateUserPermissionRequest_Permission) (bool, error) {

	res, err := s.iamClient.ValidateUserPermission(ctx, &iamv1.ValidateUserPermissionRequest{
		UserId: userId,
		ProjectId: projectId,
		Role: role,
	})

	return res.GetAuthorized(), err
}

func (s *MediaService) GetAssetsPage(ctx context.Context, projectId string) 

func NewMediaService(	
	authClient authv1.AuthServiceClient, 
	mediaClient mediav1.MediaServiceClient,
	) *MediaService {
	return &MediaService{
		authClient: authClient,
		mediaClient: mediaClient,
	}
}