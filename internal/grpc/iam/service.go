package iamgrpc

import "github.com/lucas-woo/cloud-drive/internal/database"

type Service struct {

}

func NewIamService(iamResources *database.IamResources) *Service {
	return &Service{}
}