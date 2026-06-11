package apikeys

import (
	"log"
	"os"

	apikeysv1 "github.com/lucas-woo/cloud-drive/api/apikeys/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func NewApiKeysClient() apikeysv1.ApiKeysServiceClient {
	
	port, found := os.LookupEnv("APIKEYS_CLIENT_PORT")

	if !found {
		log.Fatal("error with auth port env");
	}

	conn, err := grpc.NewClient("localhost:" + port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal()
	}

	client := apikeysv1.NewApiKeysServiceClient(conn)

	return client
}