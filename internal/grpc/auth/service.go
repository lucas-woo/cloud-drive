package authgrpc

import (
	"context"

	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {	
	authResources *database.AuthResources
}

// returns session id, err
func (s *Service) Register(ctx context.Context, user models.SignUpUserRequest) (string, error) {
	
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return "", nil
}

func NewService(authResources *database.AuthResources) *Service {
	return &Service{
		authResources: authResources,
	}
}