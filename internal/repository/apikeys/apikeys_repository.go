package apikeysrepository

import (
	"context"
	"database/sql"
	"errors"
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

	apiSecret, dbHash, err := utils.GenerateApiSecret()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			project_id,
			name,
			api_key,
			api_secret_hash,
			is_active,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`, config.ApiKeysTable)

	_, err = r.mysql.ExecContext(
		ctx,
		query,
		projectId[:],
		req.KeyName,
		apiKey[:],
		dbHash,
		true,
		createdAt,
	)

	if err != nil {
		return nil, err
	}

	return &dto.GenerateNewApiKeyResponse{
		CreatedAt: createdAt,
		ApiKey: apiKey,
		ApiSecret: apiSecret,
	}, nil
}


func (r *ApiKeysRepository) AddAPIKeyPermission(ctx context.Context, apiKey uuid.UUID, permission string) error {

	query := fmt.Sprintf(`
		INSERT IGNORE INTO %s (
			api_key,
			permission
		) VALUES (?, ?)
	`, config.ApiKeyPermissionsTable)

	_, err := r.mysql.ExecContext(ctx, query, apiKey[:],permission,)

	return err
}

func (r *ApiKeysRepository) ValidateApiKeyPermission(ctx context.Context,req *dto.ValidateApiKeyPermissionRequest) (bool, error) {

	apiKey, err := uuid.Parse(req.ApiKey)
	if err != nil {
		return false, err
	}
	query := fmt.Sprintf(`
		SELECT 
			api_secret_hash,
			is_active
		FROM %s
		WHERE api_key = ?
	`, config.ApiKeysTable)

	var (
		apiSecretHash []byte
		isActive bool
	)

	err = r.mysql.QueryRowContext(ctx, query, apiKey[:]).Scan(&apiSecretHash,&isActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if !isActive {
		return false, nil
	}

	if !utils.CompareSecret(req.ApiSecret, apiSecretHash) {
		return false, nil
	}

	permissionQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s
			WHERE api_key = ?
			  AND permission = ?
		)
	`, config.ApiKeyPermissionsTable)

	var hasPermission bool

	err = r.mysql.QueryRowContext(
		ctx,
		permissionQuery,
		req.ApiKey[:],
		req.PermissionRequest,
	).Scan(&hasPermission)

	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

func (r *ApiKeysRepository) GetAllApiKeys(ctx context.Context, projectId uuid.UUID) ([]*dto.ApiKey, error) {
	
}

func NewApiKeysRepository(mySqlClient *sql.DB) *ApiKeysRepository {
	return &ApiKeysRepository{
		mysql: mySqlClient,
	}
}