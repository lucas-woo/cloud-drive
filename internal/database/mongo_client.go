package database

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)


func ConnectMongo(ctx context.Context) (client *mongo.Client, err error) {

	uri, found := os.LookupEnv("MONGO_URI");
	if !found {
		err = errors.New("no mongo uri")
		return 
	}

	options := options.Client().ApplyURI(uri);
	client, err = mongo.Connect(options);

	if err != nil {
		return
	}

	err = client.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		return nil, err
	}

	return client, nil
}