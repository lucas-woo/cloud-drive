package iamgrpc

import iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"


type Server struct {
	iamv1.UnimplementedIAMServiceServer

}

func (UnimplementedIAMServiceServer) GenerateNewApiKey(context.Context, *GenerateNewApiKeyRequest) (*GenerateNewApiKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GenerateNewApiKey not implemented")
}
func (UnimplementedIAMServiceServer) ValidateApiKey(context.Context, *ValidateApiKeyRequest) (*ValidateApiKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ValidateApiKey not implemented")
}

func NewIamServer() *Server {
	return &Server{
		
	}
}