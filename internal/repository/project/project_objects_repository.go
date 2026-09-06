package projectrepository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
)

func (r *ProjectRepository) CreateNewObject(
    ctx context.Context,
    projectId uuid.UUID,
    objectId uuid.UUID,
    folderId uuid.UUID,
    isActive bool,
    format string,
) error {
    timeNow := time.Now().UTC()

    query := fmt.Sprintf(`
        INSERT INTO %s (
            project_id,
            folder_id,
            object_id,
            is_active,
            format,
            modified_at
        ) VALUES (?, ?, ?, ?, ?, ?)
    `, config.ProjectObjectsTable)

    _, err := r.sqldb.ExecContext(
        ctx,
        query,
        projectId[:],
        folderId[:],
        objectId[:],
        isActive,
        format,
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
			folder_name,
			folder_size,
			asset_count,
			last_upload,
			created_at,
			modified_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, config.ProjectFoldersTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		folderId[:],
		projectId[:],
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

func (r *ProjectRepository) ConfirmObjectInfo(
	ctx context.Context,
	objectId uuid.UUID,
	fileSize uint64,
) (updated bool, err error) {
	timeNow := time.Now().UTC()

	query := fmt.Sprintf(`
		UPDATE %s
		SET
			file_size = ?,
			is_pending = FALSE,
			created_at = ?,
			modified_at = ?
		WHERE object_id = ?
		  AND is_pending = TRUE
	`, config.ProjectObjectsTable)

	result, err := r.sqldb.ExecContext(
		ctx,
		query,
		fileSize,
		timeNow,
		timeNow,
		objectId[:],
	)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
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

func (r *ProjectRepository) DeleteObjectWithOwnContext(objectId uuid.UUID) {
	cleanupCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := r.DeleteObject(cleanupCtx, objectId); err != nil {
		log.Printf("failed to cleanup object %s: %v", objectId, err)
	}	
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

func (r *ProjectRepository) GetProjectIdFromObjectId(
	ctx context.Context,
	objectId uuid.UUID,
) (uuid.UUID, error) {

	query := fmt.Sprintf(`
		SELECT project_id
		FROM %s
		WHERE object_id = ?
	`, config.ProjectObjectsTable)

	var projectId uuid.UUID

	err := r.sqldb.QueryRowContext(
		ctx,
		query,
		objectId[:],
	).Scan(&projectId)
	if err != nil {
		return uuid.Nil, err
	}

	return projectId, nil
}