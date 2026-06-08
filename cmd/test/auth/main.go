package main

import (
	"context"
	"fmt"
	"log"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	authgrpc "github.com/lucas-woo/cloud-drive/internal/grpc/auth"
)

func main() {

	err := config.InitializeEnv()
	if err != nil {
		log.Fatal("env err")
	}

	client := authgrpc.NewAuthServiceClient()
	

	ctx := context.Background()

	res, err := client.SignUpUser(ctx, &authv1.SignUpUserRequest{
		Email: "lucas@gmail.com",
		Password: "1234",
		RememberMe: false,
	})

	if err != nil {
		log.Fatal("err signup")
	}
	fmt.Println("res1", res.SessionId)

	res2, err := client.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: res.SessionId,
	})
	if err != nil {
		log.Fatal("err signup")
	}
	fmt.Println("res2:", res2.LoggedIn)	

	res3, err := client.LogoutUser(ctx, &authv1.LogoutUserRequest{
		SessionId: res.SessionId,
	})
	if err != nil {
		log.Fatal("err signup")
	}
	fmt.Println("res3:", res3.LoggedOut)	
	

	// client.LoginUser(ctx, &authv1.LoginUserRequest{
	// 	Email: "lucas@gmail.com",
	// 	Password: "1234",
	// 	RememberMe: false,
	// })

}