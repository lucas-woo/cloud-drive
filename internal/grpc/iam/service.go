package iamgrpc

import "github.com/lucas-woo/cloud-drive/internal/database"

type Service struct {
	iamResources *database.IamResources
}

func (s *Service) GenerateNewApiKey() {


	//table for userid and project id
}

func NewIamService(iamResources *database.IamResources) *Service {
	return &Service{
		iamResources: iamResources,
	}
}