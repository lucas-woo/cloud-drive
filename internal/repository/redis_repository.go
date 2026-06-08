package repository

import (
	"context"
	"time"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/utils"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

const ( 
	yesRememberMe = time.Second * 60 * 60 * 24 * 7
	noRememberMe = time.Second * 60 * 60 * 24
)	

// returns generated sessionId and error
func (r *RedisRepository) SetUserSession(ctx context.Context, userId string, rememberMe bool) (string, error) {
	sessionId, err := utils.GenerateSessionId()
	if err != nil {
		return "", err
	}
	if rememberMe {
		r.client.Set(ctx, config.SessionPrefix + sessionId, userId, yesRememberMe)
		} else {
		r.client.Set(ctx, config.SessionPrefix + sessionId, userId, noRememberMe)
	}
	return sessionId, nil
}
// func (r *RedisRepository) GetUserIdBySession(ctx context.Context, sessionId string)

func (r *RedisRepository) RemoveUserSession(ctx context.Context, sessionId string) (bool, error) {
	count, err := r.client.Del(ctx, config.SessionPrefix + sessionId).Result()	
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}