package main

import (
	"log"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/gateway"
)

func main() {
	err := config.InitializeEnv()
	if err != nil {
		log.Fatal("error init env")
	}

	config.InitCookiesEnv()

	resources := database.NewGatewayResources()

	defer func(){
		resources.Close()
	}()
		
	gateway.NewServer(resources)
}