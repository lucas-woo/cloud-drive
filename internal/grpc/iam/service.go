package iamgrpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type Service struct {
	iamResources *database.IamResources
}

func (s *Service) GenerateNewApiKey(ctx context.Context, req *dto.GenerateNewApiKeyRequest) (*dto.GenerateNewApiKeyResponse, error) {

	pid, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	createdKeyResponse, err := s.iamResources.ApiKeysRepository.CreateAPIKey(ctx, req, pid)

	if err != nil {
		return nil, err
	}
	
	err = s.iamResources.ApiKeysRepository.AddAPIKeyPermission(ctx, createdKeyResponse.ApiKey, config.UploadPermission)
	if err != nil {
		return nil, err
	}

	err = s.iamResources.ApiKeysRepository.AddAPIKeyPermission(ctx, createdKeyResponse.ApiKey, config.DeletePermission)
	if err != nil {
		return nil, err
	}
	return createdKeyResponse, nil
}


func (s *Service) ValidateApiKeyPermission(ctx context.Context, req *dto.ValidateApiKeyPermissionRequest) (bool, error) {
	ok, err := s.iamResources.ApiKeysRepository.ValidateApiKeyPermission(ctx, req)
	return ok, err
}

func (s *Service) AddUserRolePermission(ctx context.Context, req *dto.AddUserRolePermissionRequest) (bool, error) {
	var role string;

	if req.UserRole == iamv1.AddUserRolePermissionRequest_PERMISSION_ADMIN_ROLE {
		role = config.ADMIN_ROLE
	}

	if len(role) == 0 {
		return false, nil
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return false, err
	}
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return false, err
	}

	err = s.iamResources.ProjectRepository.AddProjectUserRole(ctx, userId, projectId, role)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) ValidatedUserPermission(ctx context.Context, req *dto.ValidatedUserPermissionRequest) (bool, error) {
	var role string;

	if req.UserRole == iamv1.ValidateUserPermissionRequest_PERMISSION_ADMIN_ROLE {
		role = config.ADMIN_ROLE
	}
	if len(role) == 0 {
		return false, nil
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return false, err
	}
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return false, err
	}

	authorized, err := s.iamResources.ProjectRepository.CheckProjectUserRole(ctx, userId, projectId, role)

	return authorized, err
}

func (s *Service) GetAllApiKeys(ctx context.Context, projectIdString string) ([]*dto.ApiKey, error) {
	projectId, err := uuid.Parse(projectIdString)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	s.iamResources.ApiKeysRepository
}

func NewIamService(iamResources *database.IamResources) *Service {
	return &Service{
		iamResources: iamResources,
	}
}