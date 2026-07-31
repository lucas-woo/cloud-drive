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

func (m *AuthMiddleware) RedirectIfAuthenticated() gin.HandlerFunc {
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

func (m *AuthMiddleware) CheckIfSessionExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessCookie, exists := c.Get(config.CookieSession)

		if !exists {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}		
		sessionId, ok := sessCookie.(string)
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return					
		}		
		_, err := m.RedisRepository.FindUserId(c.Request.Context(), sessionId)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func NewAuthMiddleware(redisRepo *redisrepo.RedisRepository) *AuthMiddleware {
	return &AuthMiddleware{
		RedisRepository: redisRepo,
	}
}