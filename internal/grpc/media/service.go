package mediagrpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type Service struct {
	mediaResources *database.MediaResources
}

func (s *Service) CreateNewProject(ctx context.Context, createNewProjectRequest *dto.CreateNewProjectRequest) (projectName string, projectId string, err error) {

	userId, err := uuid.Parse(createNewProjectRequest.UserId)

	if err != nil {
		return
	}

	projectName, projectId, err = s.mediaResources.ProjectRepository.CreateNewProject(ctx, userId, createNewProjectRequest)

	return
}

func (s *Service) GetUploadObjectSignedUrl(ctx context.Context, req *dto.UploadObjectRequest) (url string, objectId string, err error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil{
		return 
	}
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil{
		return 
	}
	ok, err := s.mediaResources.ProjectRepository.CheckProjectUserRole(ctx, userId, projectId, config.ADMIN_ROLE)
	if err != nil{
		return 
	}
	if !ok {
		return "", "", errors.New("doesn't have role")
	}
	oId := uuid.New()
	objectId = oId.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, oId, req.Folder)
	if err != nil{
		return 
	}

	url, err = s.mediaResources.S3Repository.GetPreSignedUploadUrl(ctx, objectId)

	return 
}

func (s *Service) ConfirmObjectUpload(ctx context.Context, req *dto.LambdaS3UploadConfirmationRequest) error {
	objectId, err := uuid.Parse(req.ObjectId)
	if err != nil {
		return err
	}
	err = s.mediaResources.ProjectRepository.ConfirmObjectInfo(ctx, objectId, req.FileSize, req.Format)
	return err 
}



func (s *Service) UploadImageApiService(stream mediav1.MediaService_UploadImageApiServer) (string, error) {

	req, err  := stream.Recv()

	if err != nil {
		return "", err
	}

	imageInfo, ok := req.GetPayload().(*mediav1.UploadImageApiRequest_UploadInfo)
	
	if !ok || imageInfo == nil {
		return "", errors.New("invalid")
	}

	oId := uuid.New()
	objectId := oId.String()	

	return objectId, nil
}


func NewMediaService(mediaResources *database.MediaResources) *Service {
	return &Service{
		mediaResources: mediaResources,
	}
}