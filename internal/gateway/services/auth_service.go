package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)


type AuthService struct {
	authclient authv1.AuthServiceClient
}

func (s *AuthService) SignUp(ctx context.Context, req *dto.GatewaySignUpRequest) (string, error){
	res, err := s.authclient.SignUpUser(ctx, &authv1.SignUpUserRequest{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		RememberMe: req.RememberMe,
	})
	return res.GetSessionId(), err
}

func (s *AuthService) Login(ctx context.Context, ) {
	// s.authclient.LoginUser()
}

func (s *AuthService) Logout(ctx context.Context) {
	
}


func NewAuthServer(authclient authv1.AuthServiceClient) *AuthService {
	return &AuthService{
		authclient: authclient,
	}
}