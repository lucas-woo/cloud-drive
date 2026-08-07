package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		ProjectObjects: s.mapper.ConvertProjectObjectSlice(res.GetProjectObjects()),
		ProjectFolders: s.mapper.ConvertFolderSlice(res.GetProjectFolders()),
		ProjectCollections: s.mapper.ConvertCollectionSlice(res.GetProjectCollections()),
	}, nil
}

func (s *MediaService) GetAssetsInFolder(ctx context.Context, projectId, folderId string, assetCursor *api.AssetCursor) (*api.GetFolderAssetsResponse, error) {
	var err error
	var res *mediav1.GetAssetsInFolderResponse
	if assetCursor != nil {
		res, err = s.mediaClient.GetAssetsInFolder(ctx, &mediav1.GetAssetsInFolderRequest{
			ProjectId: projectId,
			FolderId: folderId,
			AssetCursor: &mediav1.AssetCursor{
				CreatedAt: timestamppb.New(assetCursor.CreatedAt),
				ObjectId: assetCursor.ObjectId,
			},
		})
	} else {
		res, err = s.mediaClient.GetAssetsInFolder(ctx, &mediav1.GetAssetsInFolderRequest{
			ProjectId: projectId,
			FolderId: folderId,
		})
	}

	if err != nil {
		return nil, err
	}

	return &api.GetFolderAssetsResponse{
		ProjectObjects: s.mapper.ConvertProjectObjectSlice(res.ProjectObjects),
		AssetCursor: s.mapper.ConvertAssetCursor(res.NextAssetCursor),
	}, nil
}

func (s *MediaService) GetAssetsInCollection(ctx context.Context, projectId, collectionId string, assetCursor *api.AssetCursor) (*api.GetCollectionAssetsResponse, error) {
	var err error
	var res *mediav1.GetAssetsInCollectionResponse
	if assetCursor != nil {
		res, err = s.mediaClient.GetAssetsInCollection(ctx, &mediav1.GetAssetsInCollectionRequest{
			ProjectId: projectId,
			CollectionId: collectionId,
			AssetCursor: &mediav1.AssetCursor{
				CreatedAt: timestamppb.New(assetCursor.CreatedAt),
				ObjectId: assetCursor.ObjectId,
			},
		})
	} else {
		res, err = s.mediaClient.GetAssetsInCollection(ctx, &mediav1.GetAssetsInCollectionRequest{
			ProjectId: projectId,
			CollectionId: collectionId,
		})
	}

	if err != nil {
		return nil, err
	}

	return &api.GetCollectionAssetsResponse{
		ProjectObjects: s.mapper.ConvertProjectObjectSlice(res.ProjectObjects),
		AssetCursor: s.mapper.ConvertAssetCursor(res.NextAssetCursor),
	}, nil
}

func NewMediaService(	
	authClient authv1.AuthServiceClient, 
	mediaClient mediav1.MediaServiceClient,
	iamClient iamv1.IAMServiceClient,
	mapper utils.RestMapper,
	) *MediaService {
	return &MediaService{
		authClient: authClient,
		mediaClient: mediaClient,
		iamClient: iamClient,
		mapper: mapper,
	}
}