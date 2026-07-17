package mediagrpc

import (
	"context"
	"errors"
	"io"

	"github.com/google/uuid"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"github.com/lucas-woo/cloud-drive/internal/utils"
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
	oId := uuid.Must(uuid.NewV7())
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

	if req.ErrorStatus != nil {
		s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
		return req.ErrorStatus
	}
	
	err = s.mediaResources.ProjectRepository.ConfirmObjectInfo(ctx, objectId, req.FileSize, req.Format)
	return err 
}



func (s *Service) UploadImageApiService(stream mediav1.MediaService_UploadImageApiServer) (string, error) {

	ctx := stream.Context()
	req, err  := stream.Recv()

	if err != nil {
		return "", err
	}

	imageInfo := req.GetUploadInfo()
	
	if imageInfo == nil {
		return "", errors.New("no image info")
	}

	metadata, err := utils.ExtractImageParams() //this needs to be completed later
	if err != nil {
		return "", errors.New("error with metadata")
	}	
	
	pId := imageInfo.GetProjectId()
	projectId, err := uuid.Parse(pId)
	if err != nil {
		return "", errors.New("invalid project id")
	}

	folder := imageInfo.GetFolder()

	contentType := imageInfo.GetContentType()

	objectId := uuid.Must(uuid.NewV7())
	objectIdString := objectId.String()
	
	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, folder)
	if err != nil{
		return "", errors.New("error creating new object")
	}

	pr, pw := io.Pipe()

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.mediaResources.S3Repository.UploadStreamImage(ctx, pr, objectIdString, contentType, metadata)
	}()	

	for {
		req, err = stream.Recv()

		if err == io.EOF {
			pw.Close()
			break;
		}

		if err != nil {
			s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
			pw.CloseWithError(err)
			return "", err
		}

		chunk := req.GetImageChunk()
		if chunk == nil {
			s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
			err = errors.New("no chunks")
			pw.CloseWithError(err)
			return "", err
		}

		_, err = pw.Write(chunk)
		if err != nil {
			s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
			return "", err
		}
	}

	if err = <-errChan; err != nil {
		s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
		return "", err
	}

	return objectIdString, nil
}


func (s *Service) UploadFileApiService(stream mediav1.MediaService_UploadFileApiServer) (string, error) {

	ctx := stream.Context()

	req, err  := stream.Recv()

	if err != nil {
		return "", err
	}

	fileInfo := req.GetUploadInfo()

	if fileInfo == nil {
		return "", errors.New("no image info")
	}

	pId := fileInfo.GetProjectId()
	projectId, err := uuid.Parse(pId)
	if err != nil {
		return "", errors.New("invalid project id")
	}

	folder := fileInfo.GetFolder()

	contentType := fileInfo.GetContentType()

	objectId := uuid.Must(uuid.NewV7())
	objectIdString := objectId.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, folder)
	if err != nil{
		return "", errors.New("error creating new object")
	}

	pr, pw := io.Pipe()

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.mediaResources.S3Repository.UploadFileStream(ctx, pr, objectIdString, contentType)
	}()	

	for {
		req, err = stream.Recv()

		if err == io.EOF {
			pw.Close()
			break;
		}

		if err != nil {
			s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
			pw.CloseWithError(err)
			return "", err
		}

		chunk := req.GetFileChunk()

		if chunk == nil {
			err = errors.New("no file chunk")
			s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
			pw.CloseWithError(err)
			return "", err
		}

		_, err = pw.Write(chunk)
		if err != nil {
			s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
			pw.CloseWithError(err)
			return "", err
		}
	}

	if err = <-errChan; err != nil {
		s.mediaResources.ProjectRepository.DeleteObject(ctx, objectId)
		return "", err
	}

	return objectIdString, nil	

}

func NewMediaService(mediaResources *database.MediaResources) *Service {
	return &Service{
		mediaResources: mediaResources,
	}
}