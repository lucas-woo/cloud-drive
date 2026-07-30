package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	mediagrpc "github.com/lucas-woo/cloud-drive/internal/grpc/media"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	err := config.InitializeEnv();
	if err != nil {
		log.Fatal(err)
	}

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()	

	mediaResources := database.NewMediaResources(startupCtx)

	defer func(){
		mediaResources.Close()
	}()

	port, found := os.LookupEnv("MEDIA_SERVER_PORT")

	if !found {
		log.Fatal("error with media port env");
	}	

	lis, err := net.Listen("tcp", ":" + port)

	if err != nil {
		log.Fatalf("err %v", err);
	}

	grpcServer := grpc.NewServer();

	mediav1.RegisterMediaServiceServer(grpcServer, mediagrpc.NewMediaServer(mediaResources))

	healthServer := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthServer);
	fmt.Println("media server running")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("err in starting grpc server: %v",err)
	}	
}

