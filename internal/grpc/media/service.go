package media

import (
	"context"

	"github.com/google/uuid"
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/database"
	projectmodels "github.com/lucas-woo/cloud-drive/internal/models/project"
)

type Service struct {
	mediaResources *database.MediaResources
}

func (s *Service) CreateNewProject(ctx context.Context, createNewProjectRequest *projectmodels.CreateNewProjectRequest) (projectName string, projectId string, err error) {

	res, err := s.mediaResources.AuthClient.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: createNewProjectRequest.SessionId,
	})

	if err != nil {
		return
	}

	userId, err := uuid.Parse(res.GetUserId())

	if err != nil {
		return
	}


	return
}

func NewMediaService(mediaResources *database.MediaResources) *Service {
	return &Service{
		mediaResources: mediaResources,
	}
}