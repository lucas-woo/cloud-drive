package apikeys

import (
	apikeysv1 "github.com/lucas-woo/cloud-drive/api/apikeys/v1"
)



type Server struct {
	apikeysv1.UnimplementedApiKeysServiceServer
	service *Service
}


