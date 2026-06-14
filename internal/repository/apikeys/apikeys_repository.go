package apikeysrepository

import (
	"context"
	"database/sql"
)

type ApiKeysRepository struct {
	mysql *sql.DB
}



func (r *ApiKeysRepository) GenerateApiKeyAndSecret(ctx context.Context) (error) {


}

func NewApiKeysRepository(mySqlClient *sql.DB) *ApiKeysRepository {
	return &ApiKeysRepository{
		mysql: mySqlClient,
	}
}