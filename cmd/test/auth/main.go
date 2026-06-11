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
		Password: "hello",
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
		log.Fatal("err 2")
	}
	fmt.Println("res2:", res2)	

	res3, err := client.LogoutUser(ctx, &authv1.LogoutUserRequest{
		SessionId: res.SessionId,
	})
	if err != nil {
		log.Fatal("err 3")
	}
	fmt.Println("res3:", res3.LoggedOut)	

	res4, err := client.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: res.SessionId,
	})
	if err != nil {
		log.Fatal("err 4")
	}
	fmt.Println("res4:", res4)		
	

	res5, err := client.LoginUser(ctx, &authv1.LoginUserRequest{
		Email: "lucas@gmail.com",
		Password: "hello",
		RememberMe: false,
	})
	if err != nil {
		log.Fatal("err 5")
	}
	fmt.Println("res5:", res5.SessionId)			


	res6, err := client.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: res5.SessionId,
	})
	if err != nil {
		log.Fatal("err 6")
	}
	fmt.Println("res6:", res6)	

	res7, err := client.LogoutUser(ctx, &authv1.LogoutUserRequest{
		SessionId: res5.SessionId,
	})
	if err != nil {
		log.Fatal("err 7")
	}
	fmt.Println("res7:", res7.LoggedOut)	

	res8, err := client.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: res5.SessionId,
	})
	if err != nil {
		log.Fatal("err 8")
	}
	fmt.Println("res8:", res8)			
}