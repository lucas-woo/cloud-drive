package mediagrpc

import (
	"context"
	"errors"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"github.com/lucas-woo/cloud-drive/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	mediav1.UnimplementedMediaServiceServer
	service *Service
}

func (s *Server) CreateNewProject(ctx context.Context, req *mediav1.CreateNewProjectRequest) (*mediav1.CreateNewProjectResponse, error) {

	projReq := &dto.CreateNewProjectRequest{
		UserId: req.GetUserId(),
		ProjectName: req.GetProjectName(),
		Description: req.GetDescription(),
	}

	projectName, projectId, err := s.service.CreateNewProject(ctx, projReq)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &mediav1.CreateNewProjectResponse{
		ProjectName: projectName,
		ProjectId: projectId,
	}, nil
}


func (s *Server) UploadObject(ctx context.Context, req *mediav1.UploadObjectRequest) (*mediav1.UploadObjectResponse, error) {
	url, objectId, err := s.service.GetUploadObjectSignedUrl(ctx, &dto.UploadObjectRequest{
		ProjectId: req.GetProjectId(),
		ObjectName: req.GetObjectName(),
		FolderId: req.GetFolderId(),
		IsActive: req.GetIsActive(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &mediav1.UploadObjectResponse{
		SignedUrl: url,
		ObjectId: objectId,
	}, nil
}

func (s *Server) UploadImageApi(stream mediav1.MediaService_UploadImageApiServer) error {
	objectId, err := s.service.UploadImageApiService(stream)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return stream.SendAndClose(&mediav1.UploadImageApiResponse{
		ObjectId: objectId,
	})
}

func (s *Server) UploadFileApi(stream mediav1.MediaService_UploadFileApiServer) error {
	objectId, err := s.service.UploadFileApiService(stream)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return stream.SendAndClose(&mediav1.UploadFileApiResponse{
		ObjectId: objectId,
	})
}

func (s *Server) ObjectUploadConfirmation(ctx context.Context, req *mediav1.ObjectUploadConfirmationRequest) (*mediav1.ObjectUploadConfirmationResponse, error) {

	reqStatus := req.GetStatus()

	var errorStatus error
	if reqStatus != nil {
		errorStatus = errors.New(reqStatus.GetMessage())
	}

	err := s.service.ConfirmObjectUpload(ctx, &dto.ObjectUploadConfirmationRequest{
		ObjectId: req.GetObjectId(),
		FileSize: req.GetFileSize(),
		Format: req.GetFormat(),
		ErrorStatus: errorStatus,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "method LambdaS3UploadConfirmation not implemented")		
	}
	return &mediav1.ObjectUploadConfirmationResponse{}, nil
}

func (s *Server) GetDashboard(ctx context.Context, req *mediav1.GetDashboardRequest) (*mediav1.GetDashboardResponse, error) {
	projectInfo, err := s.service.GetDashboard(ctx, &dto.GetDashboardRequest{
		UserId: req.GetUserId(),
	})

	if err != nil {
		return nil, status.Error(codes.Unimplemented, "method GetDashboard not implemented")
	}
	return &mediav1.GetDashboardResponse{
		ProjectName: projectInfo.ProjectName,
		ProjectId: projectInfo.ProjectId.String(),
		AssetsAmount: int32(projectInfo.AssetsAmount),
		Transformations: int32(projectInfo.Transformations),
		StorageBytes: projectInfo.StorageBytes,
		Description: projectInfo.Description,
	}, nil
}

func (s *Server) GetAssets(ctx context.Context, req *mediav1.GetAssetsRequest) (*mediav1.GetAssetsResponse, error) {
	var assetCursor *dto.AssetCursorRequest
	
	if req.GetAssetCursor() != nil {
		assetCursor = &dto.AssetCursorRequest{
			CreatedAt: req.GetAssetCursor().GetCreatedAt().AsTime(),
			ObjectId: req.GetAssetCursor().GetObjectId(),
		}
	}
	assets, nextCursor, err := s.service.GetAssets(ctx, &dto.GetAssetsRequest{
		ProjectId: req.GetProjectId(),
		AssetCursor: assetCursor,
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}


	return &mediav1.GetAssetsResponse{
		NextAssetCursor: &mediav1.AssetCursor{
			ObjectId: nextCursor.ObjectId.String(),
			CreatedAt: timestamppb.New(nextCursor.CreatedAt),
		},
		ProjectObjects: utils.ConvertProjectObjectsToResponse(assets),
	}, nil
}

func (s *Server) GetAllFolders(ctx context.Context, req *mediav1.GetAllFoldersRequest) (*mediav1.GetAllFoldersResponse, error) {

	allFolders, err := s.service.GetAllFolders(ctx, req.GetProjectId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &mediav1.GetAllFoldersResponse{
		ProjectFolders: utils.ConvertProjectFoldersToResponse(allFolders),
	}, nil
}


func (s *Server) GetAllCollections(ctx context.Context, req *mediav1.GetAllCollectionsRequest) (*mediav1.GetAllCollectionsResponse, error) {
	
	collections, err := s.service.GetAllCollections(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &mediav1.GetAllCollectionsResponse{
		ProjectCollections: utils.ConvertProjectCollectionsToResponse(collections),
	}, nil
}

func (s *Server) GetAssetsInFolder(ctx context.Context, req *mediav1.GetAssetsInFolderRequest) (*mediav1.GetAssetsInFolderResponse, error) {
	var assetCursor *dto.AssetCursorRequest
	
	if req.GetAssetCursor() != nil {
		assetCursor = &dto.AssetCursorRequest{
			CreatedAt: req.GetAssetCursor().GetCreatedAt().AsTime(),
			ObjectId: req.GetAssetCursor().GetObjectId(),
		}
	}

	assets, nextAssetCursor, err := s.service.GetAssetsInFolder(ctx, &dto.GetAssetsInFolderRequest{
		ProjectId: req.GetProjectId(),
		FolderId: req.GetFolderId(),
		AssetCursor: assetCursor,
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &mediav1.GetAssetsInFolderResponse{
		NextAssetCursor: &mediav1.AssetCursor{
			CreatedAt: timestamppb.New(nextAssetCursor.CreatedAt),
			ObjectId: nextAssetCursor.ObjectId.String(),
		},
		ProjectObjects: utils.ConvertProjectObjectsToResponse(assets),
	}, nil
}

func (s *Server) GetAssetsInCollection(ctx context.Context, req *mediav1.GetAssetsInCollectionRequest) (*mediav1.GetAssetsInCollectionResponse, error) {

	var assetCursor *dto.AssetCursorRequest
	
	if req.GetAssetCursor() != nil {
		assetCursor = &dto.AssetCursorRequest{
			CreatedAt: req.GetAssetCursor().GetCreatedAt().AsTime(),
			ObjectId: req.GetAssetCursor().GetObjectId(),
		}
	}
	assets, nextAssetCursor, err := s.service.GetAssetsInCollection(ctx, &dto.GetAssetsInCollectionRequest{
		ProjectId: req.GetProjectId(),
		CollectionId: req.GetCollectionId(),
		AssetCursor: assetCursor,
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &mediav1.GetAssetsInCollectionResponse{
		NextAssetCursor: &mediav1.AssetCursor{
			CreatedAt: timestamppb.New(nextAssetCursor.CreatedAt),
			ObjectId: nextAssetCursor.ObjectId.String(),
		},
		ProjectObjects: utils.ConvertProjectObjectsToResponse(assets),
	}, nil
}

func (s *Server) CreateNewFolder(ctx context.Context,req *mediav1.CreateNewFolderRequest) (*mediav1.CreateNewFolderResponse, error) {
	createdId, err := s.service.CreateNewFolder(ctx, &dto.CreateNewFolderRequest{
		FolderName: req.GetFolderName(),
		ProjectId: req.GetProjectId(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &mediav1.CreateNewFolderResponse{
		FolderId: createdId.FolderId.String(),
	}, nil
}

func (s *Server) CreateNewCollection(ctx context.Context, req *mediav1.CreateNewCollectionRequest) (*mediav1.CreateNewCollectionResponse, error) {
	createdId, err := s.service.CreateNewCollection(ctx, &dto.CreateNewCollectionRequest{
		ProjectId: req.GetProjectId(),
		CreatorId: req.GetCreatorId(),
		Name: req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &mediav1.CreateNewCollectionResponse{
		CollectionId: createdId.CollectionId.String(),
	}, nil

}

func (s *Server) GetAssetsPage(ctx context.Context, req *mediav1.GetAssetsPageRequest) (*mediav1.GetAssetsPageResponse, error) {

	objects, nextCursor, err := s.service.GetAssets(ctx, &dto.GetAssetsRequest{
		ProjectId: req.GetProjectId(),
		AssetCursor: nil,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	folders, err := s.service.GetAllFolders(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	collections, err := s.service.GetAllCollections(ctx, req.GetProjectId())	
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &mediav1.GetAssetsPageResponse{
		ProjectObjects: utils.ConvertProjectObjectsToResponse(objects),
		NextAssetCursor: &mediav1.AssetCursor{
			ObjectId: nextCursor.ObjectId.String(),
			CreatedAt: timestamppb.New(nextCursor.CreatedAt),
		},
		ProjectFolders: utils.ConvertProjectFoldersToResponse(folders),
		ProjectCollections: utils.ConvertProjectCollectionsToResponse(collections),
	}, nil
}

func NewMediaServer(mediaResources *database.MediaResources) *Server {
	return &Server{
		service: NewMediaService(mediaResources),
	}
}