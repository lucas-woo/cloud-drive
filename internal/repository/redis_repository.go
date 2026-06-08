package repository

import (
	"context"
	"time"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

const ( 
	yesRememberMe = time.Second * 60 * 60 * 24 * 7
	noRememberMe = time.Second * 60 * 60 * 24
)	

func (r *RedisRepository)SetUserSession (ctx context.Context, userId string, rememberMe bool) {
	if rememberMe {
		r.client.Set(ctx, config.SessionPrefix + userId, userId, yesRememberMe)
		} else {
		r.client.Set(ctx, config.SessionPrefix + userId, userId, noRememberMe)
	}

}


func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}