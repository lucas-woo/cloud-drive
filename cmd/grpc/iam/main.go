package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	iamgrpc "github.com/lucas-woo/cloud-drive/internal/grpc/iam"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	_ = config.InitializeEnv();


	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()	

	iamResources := database.NewIamResources(startupCtx)

	defer func(){
		iamResources.Close()
	}()	

	port, found := os.LookupEnv("IAM_SERVER_PORT")

	if !found {
		log.Fatal("error with iam port env");
	}	

	lis, err := net.Listen("tcp", ":" + port)

	if err != nil {
		log.Fatalf("err %v", err);
	}

	grpcServer := grpc.NewServer();

	iamv1.RegisterIAMServiceServer(grpcServer, iamgrpc.NewIamServer(iamResources))


	healthServer := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthServer);
	fmt.Println("iam server running")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("err in starting grpc server: %v",err)
	}	
}