package authrepo

import "go.mongodb.org/mongo-driver/v2/mongo"

type AuthRepo struct {
	db *mongo.Collection
}