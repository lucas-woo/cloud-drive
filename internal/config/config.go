package config

import "github.com/lucas-woo/godotenv"

var (
	UserDatabaseName string = "user_database"
	UserCollectionName string = "user_credentials";
)

func InitializeEnv() error {
	return godotenv.LoadEnv()
}