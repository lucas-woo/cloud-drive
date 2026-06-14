package iamgrpc

import (
	"context"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	apikeysmodels "github.com/lucas-woo/cloud-drive/internal/models/apikeys"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type Server struct {
	iamv1.UnimplementedIAMServiceServer
	service *Service
}

func (s *Server) GenerateNewApiKey(ctx context.Context, req *iamv1.GenerateNewApiKeyRequest) (*iamv1.GenerateNewApiKeyResponse, error) {

	s.service.GenerateNewApiKey(ctx, &apikeysmodels.GenerateNewApiKeyRequest{
		SessionId: req.GetSessionId(),
		KeyName: req.GetKeyName(),
		ProjectId: req.GetProjectId(),
	})

	return nil, status.Error(codes.Unimplemented, "method GenerateNewApiKey not implemented")
}
func (s *Server) ValidateApiKey(ctx context.Context, req *iamv1.ValidateApiKeyRequest) (*iamv1.ValidateApiKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ValidateApiKey not implemented")
}


func NewIamServer(iamResources *database.IamResources) *Server {
	return &Server{
		service: NewIamService(iamResources),
	}
}