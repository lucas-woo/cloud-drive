package database

import (
	"log"

	"github.com/lucas-woo/cloud-drive/internal/repository"
)

type AuthResources struct {
	RedisRepo *repository.RedisRepository
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
	redisRepo := repository.NewRedisRepository(redisClient)

	return &AuthResources{
		AuthRepo: authRepo,
		RedisRepo: redisRepo,
	}
}