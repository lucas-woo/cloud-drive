package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
)


type AuthService struct {
	authclient authv1.AuthServiceClient
}

func (s *AuthService) Login(ctx context.Context) {
	// s.authclient.LoginUser()
}
func (s *AuthService) SignUp(ctx context.Context) {

}
func (s *AuthService) Logout(ctx context.Context) {
	
}


func NewAuthServer(authclient authv1.AuthServiceClient) *AuthService {
	return &AuthService{
		authclient: authclient,
	}
}