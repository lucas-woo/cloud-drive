package authclient

import (
	"log"
	"os"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func NewAuthServiceClient() (authv1.AuthServiceClient, *grpc.ClientConn) {
	
	addr, found := os.LookupEnv("AUTH_SERVER_ADDR")
	if !found {
		log.Fatal("AUTH_SERVER_ADDR not set")
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		conn.Close()
		log.Fatal()
	}

	client := authv1.NewAuthServiceClient(conn)

	return client, conn
}