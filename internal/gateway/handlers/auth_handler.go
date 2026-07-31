package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	
}

func (h *AuthHandler) Login(ctx *gin.Context) {
}


func (h *AuthHandler) Logout(ctx *gin.Context) {
	
}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}