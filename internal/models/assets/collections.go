package assetsmodels

import (
	"time"

	"github.com/google/uuid"
)

type CollectionModel struct {
	CollectionId uuid.UUID `bson:"_id"`
	ProjectId uuid.UUID `bson:"project_id"`
	CreatedBy uuid.UUID `bson:"creator_id"`
	Name        string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
	CreatedAt    time.Time `bson:"created_at"`
	LastModified time.Time `bson:"last_modified"`
	IsPublic bool `bson:"is_public"`
}