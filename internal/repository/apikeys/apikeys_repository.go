package apikeysrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"github.com/lucas-woo/cloud-drive/internal/utils"
)

type ApiKeysRepository struct {
	mysql *sql.DB
}


func (r *ApiKeysRepository) CreateAPIKey(ctx context.Context, req *dto.GenerateNewApiKeyRequest, projectId uuid.UUID) (*dto.GenerateNewApiKeyResponse, error) {

	apiKey, err := utils.GenerateApiKey()
	if err != nil {
		return nil, err
	}
	apiSecret, err := utils.GenerateApiSecret()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id,
			project_id,
			name,
			api_key,
			api_secret,
			is_active,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, config.ApiKeysTable)

	apiId := uuid.New()

	_, err = r.mysql.ExecContext(ctx, query,apiId[:],projectId[:],req.KeyName, apiKey, apiSecret,true,createdAt)

	if err != nil {
		return nil, err
	}

	return &dto.GenerateNewApiKeyResponse{
		CreatedAt: createdAt,
		ApiKey: apiKey,
		ApiSecret: apiSecret,
		ApiId: apiId,
	}, nil
}


func (r *ApiKeysRepository) AddAPIKeyPermission(ctx context.Context, apiKeyID uuid.UUID, permission string) error {

	query := fmt.Sprintf(`
		INSERT IGNORE INTO %s (
			api_key_id,
			permission
		) VALUES (?, ?)
	`, config.ApiKeyPermissionsTable)

	_, err := r.mysql.ExecContext(ctx, query, apiKeyID[:],permission,)

	return err
}

func (r *ApiKeysRepository) ValidateApiKeyPermission(ctx context.Context, req *dto.ValidateApiKeyPermissionRequest) (bool, error) {

	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s ak
			INNER JOIN %s ap
				ON ak.id = ap.api_key_id
			WHERE ak.api_key = ?
			  AND ak.api_secret = ?
			  AND ak.is_active = TRUE
			  AND ap.permission = ?
		)
	`, config.ApiKeysTable, config.ApiKeyPermissionsTable)

	var exists bool

	err := r.mysql.QueryRowContext(ctx, query, req.ApiKey, req.ApiSecret, req.PermissionRequest).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}


func NewApiKeysRepository(mySqlClient *sql.DB) *ApiKeysRepository {
	return &ApiKeysRepository{
		mysql: mySqlClient,
	}
}