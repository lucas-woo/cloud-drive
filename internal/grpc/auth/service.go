package authgrpc

import "github.com/lucas-woo/cloud-drive/internal/database"

type Service struct {	
	authResources *database.AuthResources
}

// returns session id, err
func (s *Service) Register() (string, error) {
	return "", nil
}

func NewService(authResources *database.AuthResources) *Service {
	return &Service{
		authResources: authResources,
	}
}