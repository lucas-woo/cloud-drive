package config

import (
	"time"

	"github.com/lucas-woo/godotenv"
)

var (
	UserDatabaseName string = "user_database"
	IamDatabaseName string = "iam_database";
	MediaDatabaseName string = "media_database"

	UserCollectionName string = "user_credentials";
	ProjectCollectionName string = "projects";
	ApiKeysCollectionName string = "api_keys";


	ProjectUserRolesTable = "project_user_roles"
	ApiKeysTable = "api_keys"
	ApiKeyPermissionsTable = "api_key_permissions"
	ProjectObjectsTable = "project_objects"

	UploadPermission = "upload"
	DeletePermission = "delete"

	ADMIN_ROLE = "admin"

	PreSignedUrlTime time.Duration = time.Minute * 3

	SessionPrefix string = "session:"

	//cookies
	CookieSession string = "session_id"
	CookieSessionMaxAge int
	CookieSessionPath string
	CookieSessionDomain string
	CookieSessionSecure bool
	CookieSessionHttpOnly bool 	
)


func InitializeEnv() error {
	return godotenv.LoadEnv()
}