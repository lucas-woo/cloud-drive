package authgrpc

import (
	"context"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	authResources *database.AuthResources
}

func (s *Server) SignUpUser(context.Context, *authv1.SignUpUserRequest) (*authv1.SignUpUserResponse, error) {
	
	return nil, status.Error(codes.Unimplemented, "method SignUpUser not implemented")
}
func (Server) LoginUser(context.Context, *authv1.LoginUserRequest) (*authv1.LoginUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method LoginUser not implemented")
}
func (Server) LogoutUser(context.Context, *authv1.LogoutUserRequest) (*authv1.LogoutUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method LogoutUser not implemented")
}

func NewAuthServer(authResources *database.AuthResources) *Server {
	return &Server{
		authResources: authResources,
	}
}