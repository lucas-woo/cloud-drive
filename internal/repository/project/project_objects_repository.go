package projectrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
)

func (r *ProjectRepository) CreateNewObject(ctx context.Context, projectId uuid.UUID, objectId uuid.UUID, folder string) error {
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
		time.Now().UTC(),
	)

	return err
}