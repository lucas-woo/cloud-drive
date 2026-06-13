package iamgrpc

import (
	"context"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type Server struct {
	iamv1.UnimplementedIAMServiceServer
	service *Service
}

func (s *Server) GenerateNewApiKey(ctx context.Context, req *iamv1.GenerateNewApiKeyRequest) (*iamv1.GenerateNewApiKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GenerateNewApiKey not implemented")
}
func (s *Server) ValidateApiKey(ctx context.Context, req *iamv1.ValidateApiKeyRequest) (*iamv1.ValidateApiKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ValidateApiKey not implemented")
}


func NewIamServer() *Server {
	return &Server{

	}
}