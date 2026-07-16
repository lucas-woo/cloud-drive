package projectrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
)

func (r *ProjectRepository) CreateNewObject(ctx context.Context, projectId uuid.UUID, objectId uuid.UUID, folder string) error {
	timeNow := time.Now().UTC()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			project_id,
			object_id,
			folder,
			modified_at
		) VALUES (?, ?, ?, ?)
	`, config.ProjectObjectsTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		projectId[:],
		objectId[:],
		folder,
		timeNow,
	)

	return err
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