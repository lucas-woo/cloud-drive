package projectrepository

import "go.mongodb.org/mongo-driver/v2/mongo"

type ProjectRepository struct {
	db *mongo.Collection
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		
	}
}