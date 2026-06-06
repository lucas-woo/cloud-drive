package repository

import (
	"github.com/lucas-woo/cloud-drive/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepo struct {
	db *mongo.Collection
	
}

func NewAuthRepo(mongoClient *mongo.Client) *AuthRepo {
	return &AuthRepo{
		db: mongoClient.Database(config.CredentialDatabaseName).Collection(config.CredentialCollection),
	}
}