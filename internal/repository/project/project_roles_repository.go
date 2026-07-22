package projectrepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
)


func (r *ProjectRepository) AddProjectUserRole(ctx context.Context,userId uuid.UUID, projectId uuid.UUID, role string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, project_id, role)
		VALUES (?, ?, ?)
	`, config.ProjectUserRolesTable)

	_, err := r.sqldb.ExecContext(ctx ,query, userId[:], projectId[:], role)
	return err
}

func (r *ProjectRepository) CheckProjectUserRole(ctx context.Context, userId uuid.UUID, projectId uuid.UUID, role string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1
			FROM %s
			WHERE user_id = ?
			AND project_id = ?
			AND role = ?
		)
	`, config.ProjectUserRolesTable)

	var exists bool

	err := r.sqldb.QueryRowContext(ctx, query,userId[:],projectId[:],role,).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ProjectRepository) DeleteProjectUserRole(ctx context.Context, userId uuid.UUID, projectId uuid.UUID, role string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id = ?
		  AND project_id = ?
		  AND role = ?
	`, config.ProjectUserRolesTable)

	_, err := r.sqldb.ExecContext(ctx, query, userId[:], projectId[:], role)
	return err
}
