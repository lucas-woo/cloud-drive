package dto

type GatewaySignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email string `json:"email" binding:"required"`
	RememberMe bool `json:"rememberMe" binding:"required"`
}


type GatewayLoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RememberMe bool `json:"rememberMe" binding:"required"`
}