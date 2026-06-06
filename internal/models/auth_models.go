package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type UserLogin struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	Email string `bson:"email,omitempty"`
	Hash string `bson:"hash,omitempty"`
	UserID uuid.UUID `bson:"user_id,omitempty"` 
}