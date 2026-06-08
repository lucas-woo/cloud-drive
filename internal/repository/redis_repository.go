package repository

import "github.com/redis/go-redis/v9"

type RedisRepository struct {
	client *redis.Client
}

func (r *RedisRepository)SetUserSession (userId string) {
	
}


func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}