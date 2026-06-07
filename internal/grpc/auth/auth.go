package authgrpc

import (
	"context"
	"errors"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	authResources *database.AuthResources
}

func (s *Server) SignUpUser(ctx context.Context, req *authv1.SignUpUserRequest) (*authv1.SignUpUserResponse, error) {
	err := parseSignUpUserRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	
	return nil, status.Error(codes.Unimplemented, "method SignUpUser not implemented")
}
func (s *Server) LoginUser(ctx context.Context, req *authv1.LoginUserRequest) (*authv1.LoginUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method LoginUser not implemented")
}
func (s *Server) LogoutUser(ctx context.Context, req *authv1.LogoutUserRequest) (*authv1.LogoutUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method LogoutUser not implemented")
}

func (s *Server) ValidateUserSession(ctx context.Context,req  *authv1.ValidateUserSessionRequest) (*authv1.ValidateUserSessionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ValidateUserSession not implemented")
}


func parseSignUpUserRequest(signupRequest *authv1.SignUpUserRequest) (error) {

	var err error

	if len(signupRequest.Email) == 0 {
		newErr := errors.New("invalid email")
		err = errors.Join(err, newErr)
	}
	if len(signupRequest.Username) == 0 {
		newErr := errors.New("invalid username")
		err = errors.Join(err, newErr)		
	}
	if len(signupRequest.Password) == 0 {
		newErr := errors.New("invalid password")
		err = errors.Join(err, newErr)
	}

	return err
}

func NewAuthServer(authResources *database.AuthResources) *Server {
	return &Server{
		authResources: authResources,
	}
}