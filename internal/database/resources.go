package database

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	apikeysv1 "github.com/lucas-woo/cloud-drive/api/apikeys/v1"
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"

	apikeysclient "github.com/lucas-woo/cloud-drive/internal/grpc/apikeys/client"
	authclient "github.com/lucas-woo/cloud-drive/internal/grpc/auth/client"
	iamclient "github.com/lucas-woo/cloud-drive/internal/grpc/iam/client"
	mediaclient "github.com/lucas-woo/cloud-drive/internal/grpc/media/client"
	apikeysrepository "github.com/lucas-woo/cloud-drive/internal/repository/apikeys"
	"github.com/lucas-woo/cloud-drive/internal/repository/auth"
	projectrepository "github.com/lucas-woo/cloud-drive/internal/repository/project"
	redisrepo "github.com/lucas-woo/cloud-drive/internal/repository/redis"
	s3repository "github.com/lucas-woo/cloud-drive/internal/repository/s3"
)

type AuthResources struct {
	RedisRepo *redisrepo.RedisRepository
	AuthRepo *authrepo.AuthRepo
}

type MediaResources struct {
	ProjectRepository *projectrepository.ProjectRepository
	S3Repository *s3repository.S3Repository
	AuthClient authv1.AuthServiceClient
}

type IamResources struct {
	ProjectRepository *projectrepository.ProjectRepository
	ApiKeysRepository *apikeysrepository.ApiKeysRepository
	AuthClient authv1.AuthServiceClient	
}

type GatewayResources struct {
	RedisRepo *redisrepo.RedisRepository
	AuthClient authv1.AuthServiceClient	
	IamClient iamv1.IAMServiceClient
	MediaClient mediav1.MediaServiceClient
	ApiKeysClient apikeysv1.ApiKeysServiceClient
}

func NewAuthResources() *AuthResources {

	mongoClient, err := ConnectMongo()
	if err != nil {
		log.Fatal(err.Error())
	}
	authRepo := authrepo.NewAuthRepo(mongoClient)

	redisClient, err := ConnectRedis()
	if err != nil {
		log.Fatal(err.Error())
	}
	redisRepo := redisrepo.NewRedisRepository(redisClient)

	return &AuthResources{
		AuthRepo: authRepo,
		RedisRepo: redisRepo,
	}
}

func NewMediaResources() *MediaResources {


	accessKey := os.Getenv("AWS_S3_ACCESS_KEY")
	secretKey := os.Getenv("AWS_S3_SECRET_KEY")
	region := os.Getenv("AWS_S3_REGION")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	staticProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(staticProvider),
	)
	
	if err != nil {
		log.Fatal("error with aws credentials")
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	uploader := transfermanager.New(s3Client)

	s3repo := s3repository.NewS3Repository(s3Client, presignClient, bucketName, uploader)

	mongoClient, err := ConnectMongo()
	if err != nil {
		log.Fatal(err.Error())
	}

	mysqlClient := ConnectMySql()

	projectRepo := projectrepository.NewProjectRepository(mongoClient, mysqlClient)
	
	authClient := authclient.NewAuthServiceClient()

	return &MediaResources{
		ProjectRepository: projectRepo,
		AuthClient: authClient,
		S3Repository: s3repo,
	}
}

func NewIamResources() *IamResources {

	mongoClient, err := ConnectMongo()
	if err != nil {
		log.Fatal(err.Error())
	}

	mysqlClient := ConnectMySql()

	projectRepo := projectrepository.NewProjectRepository(mongoClient, mysqlClient)
	apikeysRepo := apikeysrepository.NewApiKeysRepository(mysqlClient)

	authClient := authclient.NewAuthServiceClient()
	

	return &IamResources{
		ProjectRepository: projectRepo,
		ApiKeysRepository: apikeysRepo,
		AuthClient: authClient,

	}	
}


func NewGatewayResources() *GatewayResources {

	redisClient, err := ConnectRedis()
	if err != nil {
		log.Fatal(err.Error())
	}

	redisRepo := redisrepo.NewRedisRepository(redisClient)

	authClient := authclient.NewAuthServiceClient()

	iamClient := iamclient.NewIamServiceClient()

	mediaClient := mediaclient.NewMediaServiceClient()

	apikeysClient := apikeysclient.NewApiKeysClient()

	return &GatewayResources{
		ApiKeysClient: apikeysClient,
		IamClient: iamClient,
		AuthClient: authClient,
		RedisRepo: redisRepo,
		MediaClient: mediaClient,
	}
}