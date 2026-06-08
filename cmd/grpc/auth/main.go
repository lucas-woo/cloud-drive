package main

import (
	"fmt"
	"log"
	"net"
	"os"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	authenticationgrpc "github.com/lucas-woo/cloud-drive/internal/grpc/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	err := config.InitializeEnv();
	if err != nil {
		log.Fatal(err)
	}
	config.InitCookiesEnv()

	authResources := database.NewAuthResources()

	port, found := os.LookupEnv("AUTH_CLIENT_PORT")

	if !found {
		log.Fatal("error with auth port env");
	}	

	lis, err := net.Listen("tcp", ":" + port)

	if err != nil {
		log.Fatalf("err %v", err);
	}

	grpcServer := grpc.NewServer();

	authv1.RegisterAuthServiceServer(grpcServer, authenticationgrpc.NewAuthServer(authResources))

	healthServer := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthServer);
	fmt.Println("auth server running")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("err in starting grpc server: %v",err)
	}	
}

