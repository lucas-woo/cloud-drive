package middlewares

import redisrepo "github.com/lucas-woo/cloud-drive/internal/repository/redis"


type Middlewares struct {
	redisRepository *redisrepo.RedisRepository
}

func NewMiddlewares(redisRepo *redisrepo.RedisRepository) *Middlewares {
	return &Middlewares{
		redisRepository: redisRepo,
	}
}