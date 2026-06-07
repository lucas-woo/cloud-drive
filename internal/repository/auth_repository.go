package repository

import (
	"context"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuthRepo struct {
	db *mongo.Collection
	
}

func (r *AuthRepo) EmailTaken(ctx context.Context, email string) (bool, error) {
	filter := bson.D{
		bson.E{Key: "email", Value: email},
	}
	opts := options.Count().SetLimit(1)
	count, err := r.db.CountDocuments(ctx, filter, opts)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func NewAuthRepo(mongoClient *mongo.Client) *AuthRepo {
	return &AuthRepo{
		db: mongoClient.Database(config.UserDatabaseName).Collection(config.UserCollectionName),
	}
}