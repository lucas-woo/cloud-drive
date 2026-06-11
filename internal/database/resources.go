package database

import (
	"log"

	"github.com/lucas-woo/cloud-drive/internal/repository/auth"
	redisrepo "github.com/lucas-woo/cloud-drive/internal/repository/redis"
)

type AuthResources struct {
	RedisRepo *redisrepo.RedisRepository
	AuthRepo *authrepo.AuthRepo
}

type MediaResources struct {
	
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
	return &MediaResources{}
}