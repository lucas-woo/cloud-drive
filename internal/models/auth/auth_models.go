package authmodels

import (
	"github.com/google/uuid"
)


type UserModel struct {
	UserID uuid.UUID `bson:"_id"`
	Email  string    `bson:"email,omitempty"`
	Hash   []byte    `bson:"hash,omitempty"`
}