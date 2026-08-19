package main

import (

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/gateway"
)

func main() {
	_ = config.InitializeEnv()



	resources := database.NewGatewayResources()

	defer func(){
		resources.Close()
	}()
		
	server := gateway.NewServer(resources)
	server.Run()
}