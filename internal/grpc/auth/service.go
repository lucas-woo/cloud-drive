package authgrpc

import "github.com/lucas-woo/cloud-drive/internal/database"

type Service struct {	
	authResources *database.AuthResources
}

func NewService(authResources *database.AuthResources) *Service {
	return &Service{
		authResources: authResources,
	}
}