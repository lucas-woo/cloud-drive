package iamgrpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	apikeysmodels "github.com/lucas-woo/cloud-drive/internal/models/apikeys"
)

type Service struct {
	iamResources *database.IamResources
}

func (s *Service) GenerateNewApiKey(ctx context.Context, req *apikeysmodels.GenerateNewApiKeyRequest) (string, string, time.Time, error) {

	res, err := s.iamResources.AuthClient.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: req.SessionId,
	})
	if err != nil {
		return "", "", time.Time{}, err
	}
	uid, err := uuid.Parse(res.UserId)
	if err != nil {
		return "", "", time.Time{}, err
	}
	pid, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return "", "", time.Time{}, err
	}
	
	ok, err := s.iamResources.ProjectRepository.CheckProjectUserRole(uid, pid, config.ADMIN_ROLE)
	if err != nil {
		return "", "", time.Time{}, err
	}
	if !ok {
		return "", "", time.Time{}, errors.New("not allowed")
	}
	//table for userid and project id
	return 
}

func NewIamService(iamResources *database.IamResources) *Service {
	return &Service{
		iamResources: iamResources,
	}
}