package main

import (
	"log"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
)

func main() {
	err := config.InitializeEnv()
	if err != nil {
		log.Fatal("error init env")
	}

	config.InitCookiesEnv()

	resources := database.NewGatewayResources()
	
}