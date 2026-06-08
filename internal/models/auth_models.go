package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type UserModel struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	Email string `bson:"email,omitempty"`
	Hash []byte `bson:"hash,omitempty"`
	UserID uuid.UUID `bson:"user_id,omitempty"` 
}

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