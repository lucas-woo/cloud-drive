package projectrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

func (r *ProjectRepository) CreateNewObject(
	ctx context.Context,
	projectId uuid.UUID,
	objectId uuid.UUID,
	folderId uuid.UUID,
	isActive bool,
) error {
	timeNow := time.Now().UTC()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			project_id,
			folder_id,
			object_id,
			is_active,
			modified_at
		) VALUES (?, ?, ?, ?, ?)
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		projectId[:],
		folderId[:],
		objectId[:],
		isActive,
		timeNow,
	)

	return err
}

func (r *ProjectRepository) CreateRootFolder(ctx context.Context, projectId uuid.UUID) (uuid.UUID, error) {
	folderId := uuid.New()
	now := time.Now().UTC()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			folder_id,
			project_id,
			parent_folder_id,
			folder_name,
			folder_size,
			asset_count,
			last_upload,
			created_at,
			modified_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, config.ProjectFoldersTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		folderId[:],
		projectId[:],
		nil,     
		"Home",
		0,       
		0,       
		nil,     
		now,
		now,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return folderId, nil
}

func (r *ProjectRepository) ConfirmObjectInfo(ctx context.Context, objectId uuid.UUID, fileSize uint64, format string) error {
	timeNow := time.Now().UTC()
	query := fmt.Sprintf(`
		UPDATE %s
		SET
			file_size = ?,
			format = ?,
			created_at = ?,
			modified_at = ?
		WHERE object_id = ?
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		fileSize,
		format,
		timeNow,
		timeNow,
		objectId[:],
	)

	return err
}

func (r *ProjectRepository) UpdateObjectInfo(ctx context.Context, objectId uuid.UUID, fileSize uint64, format string) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET
			file_size = ?,
			format = ?,
			modified_at = ?
		WHERE object_id = ?
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		fileSize,
		format,
		time.Now().UTC(),
		objectId[:],
	)
	return err
}

func (r *ProjectRepository) DeleteObject(ctx context.Context, objectId uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE object_id = ?
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		objectId[:],
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) ActivateObject(ctx context.Context, objectId uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET
			is_active = TRUE,
			modified_at = ?
		WHERE object_id = ?
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		time.Now().UTC(),
		objectId[:],
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) DisableObject(ctx context.Context, objectId uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET
			is_active = FALSE,
			modified_at = ?
		WHERE object_id = ?
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		time.Now().UTC(),
		objectId[:],
	)
	if err != nil {
		return err
	}

	return nil
}


func (r *ProjectRepository) GetAssets(
	ctx context.Context,
	projectId uuid.UUID,
	cursor *dto.AssetCursor,
	limit int,
) ([]*dto.ProjectObject, error) {

	if limit <= 0 {
		limit = 40
	}

	query := `
		SELECT
			project_id,
			collection_id,
			folder_id,
			object_id,
			file_size,
			format,
			is_active,
			created_at,
			modified_at
		FROM project_objects
		WHERE project_id = ?
	`

	args := []any{projectId[:]}

	if cursor != nil {
		query += `
			AND (
				created_at < ?
				OR (created_at = ? AND object_id < ?)
			)
		`
		args = append(args,
			cursor.CreatedAt,
			cursor.CreatedAt,
			cursor.ObjectId[:],
		)
	}

	query += `
		ORDER BY created_at DESC, object_id DESC
		LIMIT ?
	`

	args = append(args, limit)

	rows, err := r.sqldb.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*dto.ProjectObject

	for rows.Next() {
		var asset dto.ProjectObject

		err := rows.Scan(
			&asset.ProjectId,
			&asset.CollectionId,
			&asset.FolderId,
			&asset.ObjectId,
			&asset.FileSize,
			&asset.Format,
			&asset.IsActive,
			&asset.CreatedAt,
			&asset.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *ProjectRepository) GetAllFolders(ctx context.Context, projectId uuid.UUID) {

}

func (r *ProjectRepository) GetAllCollections() {

}