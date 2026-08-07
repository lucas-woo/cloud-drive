package services

import iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"

type IamService struct {
	iamClient *iamv1.IAMServiceClient
}

func NewIamService(iamClient *iamv1.IAMServiceClient) *IamService {
	return &IamService{
		iamClient: iamClient,
	}
}