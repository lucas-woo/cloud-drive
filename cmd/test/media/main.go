package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	authclient "github.com/lucas-woo/cloud-drive/internal/grpc/auth/client"
	mediaclient "github.com/lucas-woo/cloud-drive/internal/grpc/media/client"
)

func main() {

	_ = config.InitializeEnv()


	mediaClient, conn1 := mediaclient.NewMediaServiceClient()
	authClient, conn2 := authclient.NewAuthServiceClient()

	defer func(){
		conn1.Close()
		conn2.Close()
	}()

	ctx := context.Background()

	authResponse, err := authClient.SignUpUser(ctx, &authv1.SignUpUserRequest{
		Username: "lucass",
		Email: "8@gmail.com",
		Password: "1234",
		RememberMe: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	userIdReq, err := authClient.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: authResponse.GetSessionId(),
	})
	if err != nil {
		log.Fatal(err)
	}

	proj, err := mediaClient.CreateNewProject(ctx, &mediav1.CreateNewProjectRequest{
		UserId: userIdReq.GetUserId(),
	})
	if err != nil {
		log.Fatal(err)
	}
	

	stream, err := mediaClient.UploadImageApi(ctx)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("image.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = stream.Send(&mediav1.UploadImageApiRequest{
		Payload: &mediav1.UploadImageApiRequest_UploadInfo{
			UploadInfo: &mediav1.ImageUploadInfo{
				ProjectId: proj.GetProjectId(),
				ObjectName: "t",
				FolderId: "",
				ContentType: "image/jpeg",
				Transformations: &mediav1.ImageTransformations{
					Scale: &mediav1.Scale{
						Width: 10,
						Height: 10,
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 64*1024)

	for {
		n, err := file.Read(buffer)

		if n > 0 {
			err = stream.Send(&mediav1.UploadImageApiRequest{
				Payload: &mediav1.UploadImageApiRequest_ImageChunk{
					ImageChunk: buffer[:n],
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			fmt.Println("finished")
			break;
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("here 1")
	err = stream.CloseSend()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("here 2")
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}