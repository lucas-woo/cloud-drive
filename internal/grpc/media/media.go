package media

import (
	"context"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	mediav1.UnimplementedMediaServiceServer
	service *Service
}

func (s *Server) CreateNewProject(context.Context, *mediav1.CreateNewProjectRequest) (*mediav1.CreateNewProjectResponse, error) {

	
	return nil, status.Error(codes.Unimplemented, "method CreateNewProject not implemented")


}

func NewMediaServer(mediaResources *database.MediaResources) *Server {
	return &Server{
		service: NewMediaService(mediaResources),
	}
}