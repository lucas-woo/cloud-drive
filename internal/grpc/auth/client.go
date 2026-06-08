package authgrpc

import (
	"log"
	"os"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func NewAuthServiceClient() authv1.AuthServiceClient {
	
	port, found := os.LookupEnv("AUTH_CLIENT_PORT")

	if !found {
		log.Fatal("error with auth port env");
	}

	conn, err := grpc.NewClient("localhost:" + port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal()
	}

	client := authv1.NewAuthServiceClient(conn)

	return client
}