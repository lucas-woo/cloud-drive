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
	ProjectFoldersTable = "project_folders"
	ApiKeysTable = "api_keys"
	ApiKeyPermissionsTable = "api_key_permissions"
	ProjectObjectsTable = "project_objects"

	ADMIN_ROLE = "admin"

	PreSignedUrlTime time.Duration = time.Minute * 3

	SessionPrefix string = "session:"

	//assets
	AmountImagesToFetch int = 40

	//cookies
	CookieSession string = "session_id"
	CookieSessionMaxAge int
	CookieSessionPath string
	CookieSessionDomain string
	CookieSessionSecure bool
	CookieSessionHttpOnly bool 	

	//gin set/get
	GinUserId string = "user_id"


	//folder prefix
	FolderPrefix string = "/"

	S3PublicPrefix string = "public"
	S3PrivatePrefix string = "private"
	S3TransformPrefix string = "transform"
	S3ConfirmPrefix = "confirm"
)


func InitializeEnv() error {
	return godotenv.LoadEnv()
}