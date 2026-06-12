package projectrepository

import (
	"context"

	"github.com/lucas-woo/cloud-drive/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProjectRepository struct {
	db *mongo.Collection
}

func (r *ProjectRepository) CreateNewProject(ctx context.Context) () {

}


func NewProjectRepository(mongoClient *mongo.Client) *ProjectRepository {
	return &ProjectRepository{
		db: mongoClient.Database(config.MediaDatabaseName).Collection(config.ProjectCollectionName),
	}
}