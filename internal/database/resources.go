package database

import (
	"log"

	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"

	authclient "github.com/lucas-woo/cloud-drive/internal/grpc/auth/client"
	apikeysrepository "github.com/lucas-woo/cloud-drive/internal/repository/apikeys"
	"github.com/lucas-woo/cloud-drive/internal/repository/auth"
	projectrepository "github.com/lucas-woo/cloud-drive/internal/repository/project"
	redisrepo "github.com/lucas-woo/cloud-drive/internal/repository/redis"
)

type AuthResources struct {
	RedisRepo *redisrepo.RedisRepository
	AuthRepo *authrepo.AuthRepo
}

type MediaResources struct {
	ProjectRepository *projectrepository.ProjectRepository
	AuthClient authv1.AuthServiceClient
}

type IamResources struct {
	ProjectRepository *projectrepository.ProjectRepository
	ApiKeysRepository *apikeysrepository.ApiKeysRepository
	AuthClient authv1.AuthServiceClient	
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
	}
}

func NewIamResources() *IamResources {
	return &IamResources{}
}