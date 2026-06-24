package iamgrpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
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
	
	err = s.iamResources.ApiKeysRepository.AddAPIKeyPermission(ctx, createdKeyResponse.ApiId, config.UploadPermission)
	if err != nil {
		return nil, err
	}

	err = s.iamResources.ApiKeysRepository.AddAPIKeyPermission(ctx, createdKeyResponse.ApiId, config.DeletePermission)
	if err != nil {
		return nil, err
	}
	return createdKeyResponse, nil
}


func (s *Service) ValidateApiKeyPermission(ctx context.Context, req *dto.ValidateApiKeyPermissionRequest) (bool, error) {
	ok, err := s.iamResources.ApiKeysRepository.ValidateApiKeyPermission(ctx, req)
	return ok, err
}

func NewIamService(iamResources *database.IamResources) *Service {
	return &Service{
		iamResources: iamResources,
	}
}