package main

import (
	"context"
	"fmt"
	"log"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"

	iamclient "github.com/lucas-woo/cloud-drive/internal/grpc/iam/client"
)

func main() {

	err := config.InitializeEnv()
	if err != nil {
		log.Fatal("env err")
	}

	iamClient := iamclient.NewIamServiceClient()

	ctx := context.Background()

	req, err := iamClient.GenerateNewApiKey(ctx, &iamv1.GenerateNewApiKeyRequest{
		UserId: "Rlk1p-iRXM4xC2jompmFp0rslbK9YMqpDjxk5ypqpHI=",
		ProjectId: "a0698c21-cc0a-4185-8b74-9bfa572b3a07",
		KeyName: "testkey",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(req.GetApiKey())
	fmt.Println(req.GetApiSecret())

	req2, err := iamClient.ValidateApiKeyPermission(ctx, &iamv1.ValidateApiKeyPermissionRequest{
		Permission: iamv1.ValidateApiKeyPermissionRequest_PERMISSION_DELETE,
		ApiKey: req.GetApiKey(),
		ApiSecret: req.GetApiSecret(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req2.Authorized)	
}