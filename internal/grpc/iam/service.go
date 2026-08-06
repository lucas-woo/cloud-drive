package iamgrpc

import (
	"context"
	"errors"

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


	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}
	
	ok, err := s.iamResources.ProjectRepository.CheckProjectUserRole(ctx, uid, pid, config.ADMIN_ROLE)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not allowed")
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

func (s *Service) AddUserRole(ctx context.Context, req *dto.UpdateUserRoleRequest) (bool, error) {
	var role string;

	if req.UserRole == iamv1.UpdateUserRoleRequest_PERMISSION_ADMIN_ROLE {
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

func NewIamService(iamResources *database.IamResources) *Service {
	return &Service{
		iamResources: iamResources,
	}
}