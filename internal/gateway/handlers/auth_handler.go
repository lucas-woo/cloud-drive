package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func (h *AuthHandler) Login(ctx *gin.Context) {

}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}