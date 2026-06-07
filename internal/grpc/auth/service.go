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

	//need to insert to redis
	return "", nil
}

func NewService(authResources *database.AuthResources) *Service {
	return &Service{
		authResources: authResources,
	}
}