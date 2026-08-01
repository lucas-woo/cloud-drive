package projectrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	projectmodels "github.com/lucas-woo/cloud-drive/internal/models/project"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	pId, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}
	projectId = pId.String()

	newProject := &projectmodels.ProjectModel{
		CreatorId: userId,
		CreatedAt: time.Now(),
		ProjectId: pId,
		ProjectName: projectName,
		Description: projectRequest.Description,
		IsActive: true,
		Transformations: 0,
		StorageBytes: 0,
		AssetsAmount: 0,
	}

	_, err = r.mongodb.InsertOne(ctx, newProject)
	if err != nil {
		return
	}

	_, err = r.CreateRootFolder(ctx, pId)
	
	if err != nil {
		return 
	}

	err = r.AddProjectUserRole(ctx, userId, pId, config.ADMIN_ROLE)

	return
}

func (r *ProjectRepository) IncrementTransformationCount(ctx context.Context, projectId uuid.UUID) error {
	filter := bson.M{"_id": projectId}
	update := bson.M{
		"$inc": bson.M{
			"transformations": 1,
		},
	}
	_, err := r.mongodb.UpdateOne(ctx, filter, update)
	return err
}

func NewProjectRepository(mongoClient *mongo.Client, mysqlClient *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		mongodb: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
		sqldb: mysqlClient,
	}
}