package iamgrpc

import (
	"context"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type Server struct {
	iamv1.UnimplementedIAMServiceServer
	service *Service
}

func (s *Server) GenerateNewApiKey(ctx context.Context, req *iamv1.GenerateNewApiKeyRequest) (*iamv1.GenerateNewApiKeyResponse, error) {

	res, err := s.service.GenerateNewApiKey(ctx, &dto.GenerateNewApiKeyRequest{
		UserId: req.GetUserId(),
		KeyName: req.GetKeyName(),
		ProjectId: req.GetProjectId(),
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return &iamv1.GenerateNewApiKeyResponse{
		ApiKey: res.ApiKey.String(),
		ApiSecret: res.ApiSecret,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}


func (s *Server) ValidateApiKeyPermission(ctx context.Context, req *iamv1.ValidateApiKeyPermissionRequest) (*iamv1.ValidateApiKeyPermissionResponse, error) {

	var permissionRequest string
	if req.GetPermission() == iamv1.ValidateApiKeyPermissionRequest_PERMISSION_CREATE {
		permissionRequest = config.UploadPermission
	} else if req.GetPermission() == iamv1.ValidateApiKeyPermissionRequest_PERMISSION_DELETE {
		permissionRequest = config.DeletePermission
	} else {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	
	exist, err := s.service.ValidateApiKeyPermission(ctx, &dto.ValidateApiKeyPermissionRequest{
		PermissionRequest: permissionRequest,
		ApiKey: req.GetApiKey(),
		ApiSecret: req.GetApiSecret(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &iamv1.ValidateApiKeyPermissionResponse{
		Authorized: exist,
	}, nil
}

func (s *Server) AddUserRolePermission(ctx context.Context, req *iamv1.AddUserRolePermissionRequest) (*iamv1.AddUserRolePermissionResponse, error) {
	ok, err := s.service.AddUserRolePermission(ctx, &dto.AddUserRolePermissionRequest{
		UserRole: req.GetRole(),
		ProjectId: req.GetProjectId(),
		UserId: req.GetUserId(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())	
	}
	return &iamv1.AddUserRolePermissionResponse{
		Updated: ok,
	}, nil
}

func (s *Server) ValidateUserPermission(ctx context.Context, req *iamv1.ValidateUserPermissionRequest) (*iamv1.ValidateUserPermissionResponse, error) {

	authorized, err := s.service.ValidatedUserPermission(ctx, &dto.ValidatedUserPermissionRequest{
		UserRole: req.GetRole(),
		UserId: req.GetUserId(),
		ProjectId: req.GetProjectId(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())	
	}
	return &iamv1.ValidateUserPermissionResponse{
		Authorized: authorized,
	}, nil
}

func NewIamServer(iamResources *database.IamResources) *Server {
	return &Server{
		service: NewIamService(iamResources),
	}
}