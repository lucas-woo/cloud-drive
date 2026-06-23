package mediagrpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	authv1 "github.com/lucas-woo/cloud-drive/api/auth/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type Service struct {
	mediaResources *database.MediaResources
}

func (s *Service) CreateNewProject(ctx context.Context, createNewProjectRequest *dto.CreateNewProjectRequest) (projectName string, projectId string, err error) {

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

	projectName, projectId, err = s.mediaResources.ProjectRepository.CreateNewProject(ctx, userId, createNewProjectRequest)

	return
}

func (s *Service) GetUploadObjectSignedUrl(ctx context.Context, req *dto.UploadObjectRequest) (url string, err error) {
	res, err := s.mediaResources.AuthClient.ValidateUserSession(ctx, &authv1.ValidateUserSessionRequest{
		SessionId: req.SessionId,
	})
	if err != nil{
		return
	}
	userId, err := uuid.Parse(res.GetUserId())
	if err != nil{
		return 
	}
	projectId, err := uuid.Parse(res.GetUserId())
	if err != nil{
		return 
	}
	ok, err := s.mediaResources.ProjectRepository.CheckProjectUserRole(ctx, userId, projectId, config.ADMIN_ROLE)
	if err != nil{
		return 
	}
	if !ok {
		return "", errors.New("doesn't have role")
	}
	objectId := uuid.New()
	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, req.Folder)
	if err != nil{
		return 
	}

	


	return "", nil
}


func NewMediaService(mediaResources *database.MediaResources) *Service {
	return &Service{
		mediaResources: mediaResources,
	}
}