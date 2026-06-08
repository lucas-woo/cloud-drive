package authgrpc

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	service *Service
}

func (s *Server) SignUpUser(ctx context.Context, req *authv1.SignUpUserRequest) (*authv1.SignUpUserResponse, error) {
	
	sessionId, err := s.service.Register(ctx, models.SignUpUserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email: req.GetEmail(),
		RememberMe: req.GetRememberMe(),
	})
	
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}


	return &authv1.SignUpUserResponse{
		SessionId: sessionId,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *authv1.LoginUserRequest) (*authv1.LoginUserResponse, error) {
	
	sessionId, err := s.service.Login(ctx, models.LoginUserRequest{
		Email: req.GetEmail(),
		Password: req.GetPassword(),
		RememberMe: req.GetRememberMe(),
	})

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authv1.LoginUserResponse{
		SessionId: sessionId,
	}, nil
}

func (s *Server) LogoutUser(ctx context.Context, req *authv1.LogoutUserRequest) (*authv1.LogoutUserResponse, error) {
	loggedOut, err := s.service.Logout(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.LogoutUserResponse{
		LoggedOut: loggedOut,
	}, nil
}

func (s *Server) ValidateUserSession(ctx context.Context,req  *authv1.ValidateUserSessionRequest) (*authv1.ValidateUserSessionResponse, error) {
	exists, err := s.service.ValidateUserSession(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.ValidateUserSessionResponse{
		LoggedIn: exists,
	}, nil
}



func NewAuthServer(authResources *database.AuthResources) *Server {
	return &Server{
		service: NewService(authResources),
	}
}