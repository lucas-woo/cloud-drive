package services

import (
	"context"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
)

type IamService struct {
	iamClient iamv1.IAMServiceClient
}

func (s *IamService) GenerateNewApiKey(ctx context.Context, projectId, keyName string) (*api.GenerateApiKeyResponse, error) {
	res, err := s.iamClient.GenerateNewApiKey(ctx, &iamv1.GenerateNewApiKeyRequest{
		ProjectId: projectId,
		KeyName: keyName,
	})

	if err != nil {
		return nil, err
	}

	return &api.GenerateApiKeyResponse{
		ApiKey: res.GetApiKey(),
		ApiSecret: res.GetApiSecret(),
		CreatedAt: res.CreatedAt.AsTime(),
	}, nil
}

func (s *IamService) ValidateUserRole(ctx context.Context, userId, projectId string, role iamv1.ValidateUserPermissionRequest_Permission) (bool, error) {

	res, err := s.iamClient.ValidateUserPermission(ctx, &iamv1.ValidateUserPermissionRequest{
		UserId: userId,
		ProjectId: projectId,
		Role: role,
	})

	return res.GetAuthorized(), err
}

func NewIamService(iamClient iamv1.IAMServiceClient) *IamService {
	return &IamService{
		iamClient: iamClient,
	}
}