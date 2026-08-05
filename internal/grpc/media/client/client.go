package mediaclient

import (
	"log"
	"os"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func NewMediaServiceClient() (mediav1.MediaServiceClient, *grpc.ClientConn) {
	
	port, found := os.LookupEnv("MEDIA_SERVER_PORT")

	if !found {
		log.Fatal("error with media port env");
	}

	conn, err := grpc.NewClient("localhost:" + port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		conn.Close()
		log.Fatal()
	}

	client := mediav1.NewMediaServiceClient(conn)

	return client, conn
}