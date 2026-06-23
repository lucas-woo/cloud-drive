package mediagrpc

import (
	"context"

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
		SessionId: req.GetSessionId(),
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
	s.service.GetUploadObjectSignedUrl(ctx, &dto.UploadObjectRequest{
		
	})
	return nil, status.Error(codes.Unimplemented, "method UploadObject not implemented")
}


func NewMediaServer(mediaResources *database.MediaResources) *Server {
	return &Server{
		service: NewMediaService(mediaResources),
	}
}