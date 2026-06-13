package projectrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	projectmodels "github.com/lucas-woo/cloud-drive/internal/models/project"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProjectRepository struct {
	mongodb *mongo.Collection
	sqldb *sql.DB
}

func (r *ProjectRepository) CreateNewProject(ctx context.Context, userId uuid.UUID, projectRequest *projectmodels.CreateNewProjectRequest) (projectName string, projectId string, err error) {

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

	err = r.AddProjectAdmin(userId, pId)

	return
}


func (r *ProjectRepository) AddProjectAdmin(userId uuid.UUID, projectId uuid.UUID) error {
    query := fmt.Sprintf(`
        INSERT INTO %s (user_id, project_id)
        VALUES (?, ?)
    `, config.ProjectAdminTable)

    _, err := r.sqldb.Exec(query, userId[:], projectId[:])
    return err
}

func (r *ProjectRepository) IsProjectAdmin(userId uuid.UUID, projectId uuid.UUID) (bool, error) {

	query := fmt.Sprintf(`
			SELECT EXISTS(
					SELECT 1
					FROM %s
					WHERE user_id = ?
						AND project_id = ?
			)
	`, config.ProjectAdminTable)

	var exists bool

	err := r.sqldb.QueryRow(
			query,
			userId[:],
			projectId[:],
	).Scan(&exists)

	return exists, err	
}

func (r *ProjectRepository) DeleteProjectAdmin(userId uuid.UUID, projectId uuid.UUID) (error) {
	query := fmt.Sprintf(`
			DELETE FROM %s
			WHERE user_id = ?
				AND project_id = ?
	`, config.ProjectAdminTable)

	_, err := r.sqldb.Exec(query, userId[:], projectId[:])
	return err
}

func NewProjectRepository(mongoClient *mongo.Client, mysqlClient *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		mongodb: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
		sqldb: mysqlClient,
	}
}