package projectrepository

import (
	"context"
	"database/sql"
	"errors"
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

	err = r.AddProjectUserRole(ctx, userId, pId, config.ADMIN_ROLE)

	return
}

func (r *ProjectRepository) GetAllProjects(ctx context.Context, userId uuid.UUID) ([]*projectmodels.ProjectModel, error) {

	filter := bson.M{
		"creator_id": userId,
	}

	cursor, err := r.mongodb.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*projectmodels.ProjectModel

	for cursor.Next(ctx) {
		var project projectmodels.ProjectModel

		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}

		projects = append(projects, &project)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) GetOneProjectByUserId(ctx context.Context, userId uuid.UUID) (*projectmodels.ProjectModel, error) {

	filter := bson.M{
		"creator_id": userId,
	}

	var project projectmodels.ProjectModel

	err := r.mongodb.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // or return mongo.ErrNoDocuments if you prefer
		}
		return nil, err
	}

	return &project, nil
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

func (r *ProjectRepository) IncrementProjectAssetsCount(
	ctx context.Context,
	projectId uuid.UUID,
	fileSize uint64,
) error {

	_, err := r.mongodb.UpdateOne(
		ctx,
		bson.M{
			"_id": projectId,
		},
		bson.M{
			"$inc": bson.M{
				"assets_amount": 1,
				"storage_bytes": int64(fileSize),
			},
		},
	)

	return err
}

func NewProjectRepository(mongoClient *mongo.Client, mysqlClient *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		mongodb: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
		sqldb: mysqlClient,
	}
}