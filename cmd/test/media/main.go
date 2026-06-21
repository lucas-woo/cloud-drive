package main

import (
	"context"
	"fmt"
	"log"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	authclient "github.com/lucas-woo/cloud-drive/internal/grpc/auth/client"
	mediaclient "github.com/lucas-woo/cloud-drive/internal/grpc/media/client"
)

func main() {

	err := config.InitializeEnv()
	if err != nil {
		log.Fatal("env err")
	}

	mediaClient := mediaclient.NewMediaServiceClient()
	authClient := authclient.NewAuthServiceClient()

	ctx := context.Background()

	authResponse, err := authClient.SignUpUser(ctx, &authv1.SignUpUserRequest{
		Username: "lucas",
		Email: "test3@gmail.com",
		Password: "1234",
		RememberMe: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(authResponse.GetSessionId())

	res, err := mediaClient.CreateNewProject(ctx, &mediav1.CreateNewProjectRequest{
		SessionId: authResponse.GetSessionId(),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.ProjectId, res.ProjectName)
}