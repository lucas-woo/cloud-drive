package projectmodels

import (
	"time"

	"github.com/google/uuid"
)


type ProjectModel struct {
	ProjectId uuid.UUID `bson:"_id"`
	CreatorId uuid.UUID `bson:"creator_id,omitempty"`
	ProjectName string `bson:"project_name,omitempty"`
	Description string `bson:"description,omitempty"`
	IsActive bool `bson:"active,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
}
