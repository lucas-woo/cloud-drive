package apikeysrepository

import (
	"context"
	"database/sql"

	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type ApiKeysRepository struct {
	mysql *sql.DB
}

func (r *ApiKeysRepository) GenerateApiKeyAndSecret(ctx context.Context, req *dto.CreateNewProjectRequest) (error) {
	
	return nil
}

func NewApiKeysRepository(mySqlClient *sql.DB) *ApiKeysRepository {
	return &ApiKeysRepository{
		mysql: mySqlClient,
	}
}