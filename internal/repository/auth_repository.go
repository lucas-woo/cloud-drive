package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

func (r *AuthRepo) CreateNewUser(ctx context.Context, user models.SignUpUserRequest) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUserId := uuid.New()

	newUser := &models.UserModel{
		Hash: string(hash),
		Email: user.Email,
		UserID: newUserId,
	}

	_, err = r.db.InsertOne(ctx, newUser)

	if err != nil {
		return "", err
	}

	return newUserId.String(), nil
}

func (r *AuthRepo) LoginUser(ctx context.Context, user models.LoginUserRequest) (string, error) {
	
	return "", nil
}

func NewAuthRepo(mongoClient *mongo.Client) *AuthRepo {
	return &AuthRepo{
		db: mongoClient.Database(config.UserDatabaseName).Collection(config.UserCollectionName),
	}
}