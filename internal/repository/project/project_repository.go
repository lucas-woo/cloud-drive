package projectrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	projectmodels "github.com/lucas-woo/cloud-drive/internal/models/project"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProjectRepository struct {
	db *mongo.Collection
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

	_, err = r.db.InsertOne(ctx, newProject)

	return
}


func NewProjectRepository(mongoClient *mongo.Client) *ProjectRepository {
	return &ProjectRepository{
		db: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
	}
}