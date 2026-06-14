package projectrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	projectmodels "github.com/lucas-woo/cloud-drive/internal/models/project"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProjectRepository struct {
	mongodb *mongo.Collection
	sqldb *sql.DB
}

func (r *ProjectRepository) CreateNewProject(ctx context.Context, userId uuid.UUID, projectRequest *dto.CreateNewProjectRequest) (projectName string, projectId string, err error) {

	if projectRequest.ProjectName == "" {
		projectName = "New_Project"
	} else {
		projectName = projectRequest.ProjectName
	}

	pId := uuid.New()
	projectId = pId.String()

	newProject := &projectmodels.ProjectModel{
		CreatorId: userId,
		CreatedAt: time.Now(),
		ProjectId: pId,
		ProjectName: projectName,
		Description: projectRequest.Description,
		IsActive: true,
	}

	_, err = r.mongodb.InsertOne(ctx, newProject)
	if err != nil {
		return
	}

	err = r.AddProjectUserRole(ctx, userId, pId, config.ADMIN_ROLE)

	return
}


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

func NewProjectRepository(mongoClient *mongo.Client, mysqlClient *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		mongodb: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
		sqldb: mysqlClient,
	}
}