package dto

type SignUpUserRequest struct {
	Username string
	Password string
	Email string
	RememberMe bool
}

type LoginUserRequest struct {
	Username string
	Password string
	Email string
	RememberMe bool	
}