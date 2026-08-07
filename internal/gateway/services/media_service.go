package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/utils"
)

type MediaService struct {
	authClient authv1.AuthServiceClient
	mediaClient mediav1.MediaServiceClient
	iamClient iamv1.IAMServiceClient
	mapper utils.RestMapper
}

func (s *MediaService) GetDashboard(ctx context.Context, userId string) (*api.GetDashboardResponse, error) {
	
	res, err := s.mediaClient.GetDashboard(ctx, &mediav1.GetDashboardRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return &api.GetDashboardResponse{
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

func (s *MediaService) GetAssetsPage(ctx context.Context, projectId string) (*api.GetAssetsPageResponse, error) {
	res, err := s.mediaClient.GetAssetsPage(ctx, &mediav1.GetAssetsPageRequest{
		ProjectId: projectId,
	})

	if err != nil {
		return nil, err
	}

	return &api.GetAssetsPageResponse{
		AssetCursor: s.mapper.ConvertAssetCursor(res.GetNextAssetCursor()),
		
	}, nil
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