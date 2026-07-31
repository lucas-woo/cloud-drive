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

func (s *AuthService) Login(ctx context.Context, req *dto.GatewayLoginRequest) (string, error) {
	res, err := s.authclient.LoginUser(ctx, &authv1.LoginUserRequest{
		Email: req.Email,
		Password: req.Password,
		RememberMe: req.RememberMe,
	})
	return res.GetSessionId(), err
}

func (s *AuthService) Logout(ctx context.Context, sessionId string) (bool, error){
	res, err := s.authclient.LogoutUser(ctx, &authv1.LogoutUserRequest{
		SessionId: sessionId,
	})
	return res.GetLoggedOut(), err
}


func NewAuthServer(authclient authv1.AuthServiceClient) *AuthService {
	return &AuthService{
		authclient: authclient,
	}
}