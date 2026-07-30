package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/config"
	redisrepo "github.com/lucas-woo/cloud-drive/internal/repository/redis"
)

type AuthMiddleware struct {
	RedisRepository *redisrepo.RedisRepository
}

func (m *AuthMiddleware) IsAlreadyLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		sessionID, err := c.Cookie(config.CookieSession);
		
		if err != nil {
			c.Next()
			return 
		}

		_, err = m.RedisRepository.FindUserId(c.Request.Context(), sessionID);

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Next()
	}
}

func NewAuthMiddleware(redisRepo *redisrepo.RedisRepository) *AuthMiddleware {
	return &AuthMiddleware{
		RedisRepository: redisRepo,
	}
}