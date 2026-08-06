package handlers

import (
	"log"
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

	sessionId, userId, err := h.authService.SignUp(c.Request.Context(), &signUpReq)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}

	projectId, err := h.authService.CreateNewProject(c.Request.Context(), userId)
	if err != nil {
		//should either retry or delete the user
		log.Printf("error creating new user project: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}	
	
	err = h.authService.AddAdminRole(c.Request.Context(), userId, projectId)
	if err != nil {
		//should either retry or delete the user
		log.Printf("error adding admin role to project creator: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}		

	c.SetCookie(config.CookieSession, sessionId, config.CookieSessionMaxAge, config.CookieSessionPath, config.CookieSessionDomain, config.CookieSessionSecure, config.CookieSessionHttpOnly)

	c.JSON(http.StatusCreated, "created")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq dto.GatewayLoginRequest

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
	sessCookie, _ := c.Get(config.CookieSession)
	sessionId, ok := sessCookie.(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return					
	}
	ok, err := h.authService.Logout(c.Request.Context(), sessionId)

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