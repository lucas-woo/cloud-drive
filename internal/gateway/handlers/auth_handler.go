package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var signUpReq dto.GatewaySignUpRequest
	if err := c.ShouldBindBodyWithJSON(&signUpReq); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sessionId, err := h.authService.SignUp(c.Request.Context(), &signUpReq)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	
	c.SetCookie(config.CookieSession, sessionId, config.CookieSessionMaxAge, config.CookieSessionPath, config.CookieSessionDomain, config.CookieSessionSecure, config.CookieSessionHttpOnly)

	c.JSON(http.StatusCreated, "created")

}

func (h *AuthHandler) Login(c *gin.Context) {
}


func (h *AuthHandler) Logout(c *gin.Context) {
	
}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}