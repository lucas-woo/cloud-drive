package services

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)


type AuthService struct {
	authClient authv1.AuthServiceClient
	iamClient iamv1.IAMServiceClient
	mediaClient mediav1.MediaServiceClient
}

func (s *AuthService) SignUp(ctx context.Context, req *dto.GatewaySignUpRequest) (string, string, error){
	res, err := s.authClient.SignUpUser(ctx, &authv1.SignUpUserRequest{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		RememberMe: req.RememberMe,
	})
	return res.GetSessionId(), res.GetUserId(), err
}

func (s *AuthService) Login(ctx context.Context, req *dto.GatewayLoginRequest) (string, error) {
	res, err := s.authClient.LoginUser(ctx, &authv1.LoginUserRequest{
		Email: req.Email,
		Password: req.Password,
		RememberMe: req.RememberMe,
	})
	return res.GetSessionId(), err
}

func (s *AuthService) Logout(ctx context.Context, sessionId string) (bool, error){
	res, err := s.authClient.LogoutUser(ctx, &authv1.LogoutUserRequest{
		SessionId: sessionId,
	})
	return res.GetLoggedOut(), err
}

func (s *AuthService) CreateNewProject(ctx context.Context, userId string) (string, error) {
	res, err := s.mediaClient.CreateNewProject(ctx, &mediav1.CreateNewProjectRequest{
		UserId: userId,
	})
	return res.GetProjectId(), err
}

func (s *AuthService) AddAdminRole(ctx context.Context, userId string, projectId string) (error) {
	_, err := s.iamClient.AddUserRolePermission(ctx, &iamv1.AddUserRolePermissionRequest{
		Role: iamv1.AddUserRolePermissionRequest_PERMISSION_ADMIN_ROLE,
		UserId: userId,
		ProjectId: projectId,
	})
	return err
}

func NewAuthServer(authClient authv1.AuthServiceClient, mediaClient mediav1.MediaServiceClient) *AuthService {
	return &AuthService{
		authClient: authClient,
		mediaClient: mediaClient,
	}
}