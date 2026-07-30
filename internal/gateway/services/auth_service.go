package services

import (
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
)


type AuthService struct {
	authclient authv1.AuthServiceClient
}

func NewAuthServer(authclient authv1.AuthServiceClient) *AuthService {
	return &AuthService{
		authclient: authclient,
	}
}