package apikeys

import (
	"context"

	apikeysv1 "github.com/lucas-woo/cloud-drive/api/apikeys/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	apikeysv1.UnimplementedApiKeysServiceServer
	service *Service
}


func (s *Server) CreateNewProject(ctx context.Context, req *apikeysv1.CreateNewProjectRequest) (*apikeysv1.CreateNewProjectResponse, error) {
	
	return nil, status.Error(codes.Unimplemented, "method CreateNewProject not implemented")
}



func NewApiKeysServer() *Server {
	return &Server{

	}
}
