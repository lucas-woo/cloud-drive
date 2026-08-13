package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/gateway/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var signUpReq api.SignUpRequest
	if err := c.ShouldBindBodyWithJSON(&signUpReq); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sessionId, userId, err := h.authService.SignUp(c.Request.Context(), &signUpReq)
	if err != nil {
		log.Printf("error signing up: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}

	_, err = h.authService.CreateNewProject(c.Request.Context(), userId)
	if err != nil {
		//should either retry or delete the user
		log.Printf("error creating new user project: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	
	

	c.SetCookie(config.CookieSession, sessionId, config.CookieSessionMaxAge, config.CookieSessionPath, config.CookieSessionDomain, config.CookieSessionSecure, config.CookieSessionHttpOnly)

	c.JSON(http.StatusCreated, "created")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq api.LoginRequest

	if err := c.ShouldBindBodyWithJSON(&loginReq); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	
	sessionId, err := h.authService.Login(c.Request.Context(), &loginReq)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}
	
	c.SetCookie(config.CookieSession, sessionId, config.CookieSessionMaxAge, config.CookieSessionPath, config.CookieSessionDomain, config.CookieSessionSecure, config.CookieSessionHttpOnly)

	c.JSON(http.StatusCreated, "ok")
}


func (h *AuthHandler) Logout(c *gin.Context) {
	sessionId, err := c.Cookie(config.CookieSession)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return					
	}
	
	ok, err := h.authService.Logout(c.Request.Context(), sessionId)
	fmt.Println(ok, err)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return				
	}
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return					
	}
	c.JSON(http.StatusOK, "ok")
}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}