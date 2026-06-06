package database

import (
	"log"

	"github.com/lucas-woo/cloud-drive/internal/repository"
	"github.com/redis/go-redis/v9"
)

type AuthResources struct {
	RedisClient *redis.Client
	AuthRepo *repository.AuthRepo
}

func NewAuthResources() *AuthResources {

	mongoClient, err := ConnectMongo()
	if err != nil {
		log.Fatal(err.Error())
	}
	authRepo := repository.NewAuthRepo(mongoClient)

	redisClient, err := ConnectRedis()
	if err != nil {
		log.Fatal(err.Error())
	}

	return &AuthResources{
		AuthRepo: authRepo,
		RedisClient: redisClient,
	}
}