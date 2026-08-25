package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"

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
	mongoClient *mongo.Client
	redisClient *redis.Client
	RedisRepo *redisrepo.RedisRepository
	AuthRepo *authrepo.AuthRepo
}

type MediaResources struct {
	mongoClient *mongo.Client
	mysqlClient *sql.DB
	ProjectRepository *projectrepository.ProjectRepository
	S3Repository *s3repository.S3Repository
}

type IamResources struct {
	mongoClient *mongo.Client
	mysqlClient *sql.DB
	ProjectRepository *projectrepository.ProjectRepository
	ApiKeysRepository *apikeysrepository.ApiKeysRepository
}

type GatewayResources struct {
	redisClient *redis.Client
	RedisRepo *redisrepo.RedisRepository
	AuthClient authv1.AuthServiceClient	
	IamClient iamv1.IAMServiceClient
	MediaClient mediav1.MediaServiceClient
	authConn *grpc.ClientConn
	iamConn *grpc.ClientConn
	mediaConn *grpc.ClientConn
}

func NewAuthResources(ctx context.Context) *AuthResources {

	mongoClient, err := ConnectMongo(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	authRepo := authrepo.NewAuthRepo(mongoClient)

	redisClient, err := ConnectRedis(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	redisRepo := redisrepo.NewRedisRepository(redisClient)

	return &AuthResources{
		redisClient: redisClient,
		mongoClient: mongoClient,
		AuthRepo: authRepo,
		RedisRepo: redisRepo,
	}
}

func (r *AuthResources) Close() {
	if r.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting Mongo client: %v", err)
		} else {
			log.Println("MongoDB client disconnected successfully.")
		}
	}	

	if r.redisClient != nil {
		if err := r.redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		} else {
			log.Println("Redis client disconnected successfully.")
		}
	}	
}


func NewMediaResources(ctx context.Context) *MediaResources {


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

	mongoClient, err := ConnectMongo(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	mysqlClient := ConnectMySql()

	projectRepo := projectrepository.NewProjectRepository(mongoClient, mysqlClient)
	

	return &MediaResources{
		mysqlClient: mysqlClient,
		mongoClient: mongoClient,
		ProjectRepository: projectRepo,
		S3Repository: s3repo,
	}
}

func (r *MediaResources) Close() {

	if r.mysqlClient != nil {
		r.mysqlClient.Close()
	}

	if r.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting Mongo client: %v", err)
		} else {
			log.Println("MongoDB client disconnected successfully.")
		}
	}	
}

func NewIamResources(ctx context.Context) *IamResources {

	mongoClient, err := ConnectMongo(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	mysqlClient := ConnectMySql()

	projectRepo := projectrepository.NewProjectRepository(mongoClient, mysqlClient)
	apikeysRepo := apikeysrepository.NewApiKeysRepository(mysqlClient)

	

	return &IamResources{
		ProjectRepository: projectRepo,
		ApiKeysRepository: apikeysRepo,
		mysqlClient: mysqlClient,
	}	
}

func (r *IamResources) Close() {

	if r.mysqlClient != nil {
		r.mysqlClient.Close()
	}	

	if r.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting Mongo client: %v", err)
		} else {
			log.Println("MongoDB client disconnected successfully.")
		}
	}	
}

func NewGatewayResources(ctx context.Context) *GatewayResources {

	redisClient, err := ConnectRedis(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	redisRepo := redisrepo.NewRedisRepository(redisClient)

	authClient, authConn := authclient.NewAuthServiceClient()

	iamClient, iamConn := iamclient.NewIamServiceClient()

	mediaClient, mediaConn := mediaclient.NewMediaServiceClient()


	return &GatewayResources{
		IamClient: iamClient,
		AuthClient: authClient,
		RedisRepo: redisRepo,
		MediaClient: mediaClient,
		authConn: authConn,
		iamConn: iamConn,
		mediaConn: mediaConn,
	}
}

func (r *GatewayResources) Close() {
	if r.redisClient != nil {
		if err := r.redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		} else {
			log.Println("Redis client disconnected successfully.")
		}
	}

	var errs error

	if err := r.authConn.Close(); err != nil {
		errs = errors.Join(errs, err)
	}

	if err := r.iamConn.Close(); err != nil {
		errs = errors.Join(errs, err)
	}
	if err := r.mediaConn.Close(); err != nil {
		errs = errors.Join(errs, err)
	}
	
	if errs != nil {
		log.Printf("error closing grpc clients:\n %v", errs)
	}
}