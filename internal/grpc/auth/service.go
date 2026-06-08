package authgrpc

import (
	"context"
	"errors"

	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/models"
)

type Service struct {	
	authResources *database.AuthResources
}

// returns session id, err
func (s *Service) Register(ctx context.Context, user models.SignUpUserRequest) (string, error) {
	
	taken, err := s.authResources.AuthRepo.EmailTaken(ctx, user.Email)
	if err != nil || taken {
		return "", errors.New("email taken")
	}

	userId, err := s.authResources.AuthRepo.CreateNewUser(ctx, user)
	
	sessionId, err := s.authResources.RedisRepo.SetUserSession(ctx, userId, user.RememberMe)
	if err != nil || taken {
		return "", err
	}

	return sessionId, nil
}

func (s *Service) Login(ctx context.Context, user models.LoginUserRequest) (string, error) {
	userId, err := s.authResources.AuthRepo.LoginUser(ctx, user)
	if err != nil {
		return "", err
	}

	sessionId, err := s.authResources.RedisRepo.SetUserSession(ctx, userId, user.RememberMe)
	if err != nil {
		return "", err
	}	
	
	return sessionId, nil
}

func (s *Service) Logout(ctx context.Context, sessionId string) (bool, error) {

	ok, err := s.authResources.RedisRepo.RemoveUserSession(ctx, sessionId)

	return ok, err
}

func (s *Service) ValidateUserSession(ctx context.Context, sessionId string) (bool, error) {
	exists, err := s.authResources.RedisRepo.FindExist(ctx, sessionId)
	return exists, err
}

func NewService(authResources *database.AuthResources) *Service {
	return &Service{
		authResources: authResources,
	}
}
