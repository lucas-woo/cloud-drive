package apikeysrepository

import (
	"github.com/lucas-woo/cloud-drive/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ApiKeysRepository struct {
	db *mongo.Collection
}

func NewApiKeysRepository(mongoClient *mongo.Client) *ApiKeysRepository {
	return &ApiKeysRepository{
		db: mongoClient.Database(config.IamDatabaseName).Collection(config.ApiKeysCollectionName),
	}
}