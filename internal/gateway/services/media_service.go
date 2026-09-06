package services

import (
	"context"
	"io"

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

func (s *MediaService) CreateNewCollection(ctx context.Context, projectId, creatorId, name, description string) (*api.CreateCollectionResponse, error) {
	res, err := s.mediaClient.CreateNewCollection(ctx, &mediav1.CreateNewCollectionRequest{
		ProjectId: projectId,
		CreatorId: creatorId,
		Name: name,
		Description: description,
	})

	if err != nil {
		return nil, err
	}

	return &api.CreateCollectionResponse{
		CollectionId: res.GetCollectionId(),
	}, nil
}

func (s *MediaService) CreateNewFolder(ctx context.Context, projectId, name string) (*api.CreateFolderResponse, error) {
	res, err := s.mediaClient.CreateNewFolder(ctx, &mediav1.CreateNewFolderRequest{
		ProjectId: projectId,
		FolderName: name,
	})

	if err != nil {
		return nil, err
	}

	return &api.CreateFolderResponse{
		FolderId: res.GetFolderId(),
	}, nil
}

func (s *MediaService) CreateNewProject(ctx context.Context, userId, projectName, projectDescription string) (string, error) {
	res, err := s.mediaClient.CreateNewProject(ctx, &mediav1.CreateNewProjectRequest{
		UserId: userId,
		ProjectName: projectName,
		Description: projectDescription,
	})
	return res.GetProjectId(), err
}

func (s *MediaService) AddAdminRole(ctx context.Context, userId string, projectId string) (error) {
	_, err := s.iamClient.AddUserRolePermission(ctx, &iamv1.AddUserRolePermissionRequest{
		Role: iamv1.AddUserRolePermissionRequest_PERMISSION_ADMIN_ROLE,
		UserId: userId,
		ProjectId: projectId,
	})
	return err
}

func (s *MediaService) GetUploadObjectUrl(ctx context.Context, projectId, folderId, objectName, originalFileName string, isActive bool) (*api.CreateObjectResponse, error) {

	format := s.mapper.FindFormat(originalFileName)
	res, err := s.mediaClient.UploadObject(ctx, &mediav1.UploadObjectRequest{
		ProjectId: projectId,
		FolderId: folderId,
		ObjectName: objectName,
		Format: format,
		IsActive: isActive,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateObjectResponse{
		Url: res.GetSignedUrl(),
		ObjectId: res.GetObjectId(),
	}, nil
}


func (s *MediaService) UploadFile(
	ctx context.Context,
	projectId string,
	objectName string,
	originalFileName string,
	folderId string,
	isActive bool,
	contentType string,
	file io.Reader,
) (*api.UploadObjectApiResponse, error) {

	format := s.mapper.FindFormat(originalFileName)

	stream, err := s.mediaClient.UploadFileApi(ctx)
	if err != nil {
		return nil, err
	}

	err = stream.Send(&mediav1.UploadFileApiRequest{
		Payload: &mediav1.UploadFileApiRequest_UploadInfo{
			UploadInfo: &mediav1.FileUploadInfo{
				ProjectId:   projectId,
				ObjectName:  objectName,
				FolderId:    folderId,
				ContentType: contentType,
				Format:      format,
				IsActive:    isActive,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 64*1024)

	for {
		n, err := file.Read(buf)

		if n > 0 {
			if err := stream.Send(&mediav1.UploadFileApiRequest{
				Payload: &mediav1.UploadFileApiRequest_FileChunk{
					FileChunk: buf[:n],
				},
			}); err != nil {
				return nil, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return &api.UploadObjectApiResponse{
		ObjectId: res.GetObjectId(),
	}, nil
}

func (s *MediaService) UploadImage(
	ctx context.Context,
	projectId string,
	objectName string,
	originalFileName string,
	folderId string,
	isActive bool,
	contentType string,
	image io.Reader,
) (*api.UploadImageApiResponse, error) {

	format := s.mapper.FindFormat(originalFileName)

	stream, err := s.mediaClient.UploadImageApi(ctx)
	if err != nil {
		return nil, err
	}

	err = stream.Send(&mediav1.UploadImageApiRequest{
		Payload: &mediav1.UploadImageApiRequest_UploadInfo{
			UploadInfo: &mediav1.ImageUploadInfo{
				ProjectId: projectId,
				ObjectName: objectName,
				FolderId: folderId,
				ContentType: contentType,
				Format: format,
				IsActive: isActive,
				Transformations: nil,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 64*1024)

	for {
		n, err := image.Read(buf)

		if n > 0 {
			if err := stream.Send(&mediav1.UploadImageApiRequest{
				Payload: &mediav1.UploadImageApiRequest_ImageChunk{
					ImageChunk: buf[:n],
				},
			}); err != nil {
				return nil, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return &api.UploadImageApiResponse{
		ObjectId: res.GetObjectId(),
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