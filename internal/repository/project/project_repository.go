package projectrepository

import (
	"context"
	"database/sql"
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

	return
}


func (r *ProjectRepository) AddProjectAdmin() {

}
func (r *ProjectRepository) IsProjectAdmin() {

}
func (r *ProjectRepository) DeleteProjectAdmin() {

}

func NewProjectRepository(mongoClient *mongo.Client, mysqlClient *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		mongodb: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
		sqldb: mysqlClient,
	}
}