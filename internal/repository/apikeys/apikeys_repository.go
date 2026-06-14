package apikeysrepository

import (
	"database/sql"
)

type ApiKeysRepository struct {
	mysql *sql.DB
}


func NewApiKeysRepository(mySqlClient *sql.DB) *ApiKeysRepository {
	return &ApiKeysRepository{
		mysql: mySqlClient,
	}
}