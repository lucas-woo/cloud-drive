package iamclient

import (
	"log"
	"os"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func NewIamServiceClient() (iamv1.IAMServiceClient, *grpc.ClientConn) {
	
	addr, found := os.LookupEnv("IAM_SERVER_ADDR")

	if !found {
		log.Fatal("error with iam address env")
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		conn.Close()
		log.Fatal()
	}

	client := iamv1.NewIAMServiceClient(conn)

	return client, conn
}