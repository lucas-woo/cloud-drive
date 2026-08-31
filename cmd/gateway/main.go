package main

import (
	"context"
	"time"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/gateway"
)

func main() {
	_ = config.InitializeEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resources := database.NewGatewayResources(ctx)

	defer func(){
		resources.Close()
	}()
		
	server := gateway.NewServer(resources)
	server.Run()
}