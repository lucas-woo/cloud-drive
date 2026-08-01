package mediagrpc

import (
	"context"
	"errors"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		UserId: req.GetUserId(),
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
		return nil, status.Error(codes.Unimplemented, "method LambdaS3UploadConfirmation not implemented")		
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
		
	}, nil
}


func NewMediaServer(mediaResources *database.MediaResources) *Server {
	return &Server{
		service: NewMediaService(mediaResources),
	}
}