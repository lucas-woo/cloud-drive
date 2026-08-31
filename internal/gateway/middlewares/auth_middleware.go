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
		id, err := m.RedisRepository.FindUserId(c.Request.Context(), sessionID);
		if err != nil || len(id) == 0 {
			c.Next()
			return 			
		}
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (m *AuthMiddleware) IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie(config.CookieSession)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userId, err := m.RedisRepository.FindUserId(c.Request.Context(), sessionId)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(config.GinUserId, userId)
		
		c.Next()
	}
}

func (m *AuthMiddleware)EventBridgeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func NewAuthMiddleware(redisRepo *redisrepo.RedisRepository) *AuthMiddleware {
	return &AuthMiddleware{
		RedisRepository: redisRepo,
	}
}