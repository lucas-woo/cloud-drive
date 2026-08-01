package mediagrpc

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/database"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	s3repository "github.com/lucas-woo/cloud-drive/internal/repository/s3"
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
	oId, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}

	objectId = oId.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, oId, req.Folder, req.IsActive)
	if err != nil{
		return 
	}

	url, err = s.mediaResources.S3Repository.GetPreSignedUploadUrl(ctx, objectId, req.ProjectId, req.IsActive)

	return 
}

func (s *Service) ConfirmObjectUpload(ctx context.Context, req *dto.ObjectUploadConfirmationRequest) error {
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

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	req, err := stream.Recv()
	if err != nil {
		return "", err
	}

	imageInfo := req.GetUploadInfo()
	if imageInfo == nil {
		return "", errors.New("no image info")
	}

	projectId, err := uuid.Parse(imageInfo.GetProjectId())
	if err != nil {
		return "", errors.New("invalid project id")
	}

	objectId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	objectIdString := objectId.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, imageInfo.GetFolder(), imageInfo.GetIsActive())

	if err != nil {
		return "", errors.New("error creating new object")
	}

	
	var uploadSuccessful bool
	err = s.mediaResources.ProjectRepository.IncrementTransformationCount(ctx, projectId)

	if err != nil {
		return "", errors.New("error incrementing transformations")
	}
	
	defer func() {
		if !uploadSuccessful {

			s.mediaResources.ProjectRepository.DeleteObject(context.Background(), objectId)
		}
	}()


	cmd, err := utils.SelectImageProcessor(ctx, imageInfo.GetTransformations())
	if err != nil {
		return "", errors.New("error with transformations")
	}

	pipeReader, pipeWriter := io.Pipe()

	cmd.Stdin = pipeReader
	cppStdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	safeStdout := s3repository.UnseekableReader{R: cppStdout}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start processor: %w", err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.mediaResources.S3Repository.UploadStreamImage(ctx, safeStdout, objectIdString, imageInfo.GetProjectId(), imageInfo.GetContentType(), imageInfo.GetIsActive())
	}()

	for {
		req, err = stream.Recv()
		if err == io.EOF {
			pipeWriter.Close()
			break
		}
		if err != nil {
			pipeWriter.CloseWithError(err)
			return "", fmt.Errorf("stream receive error: %w", err)
		}

		chunk := req.GetImageChunk()
		if chunk == nil {
			err = errors.New("empty chunk received")
			pipeWriter.CloseWithError(err)
			return "", err
		}

		_, err = pipeWriter.Write(chunk)
		if err != nil {
			pipeWriter.CloseWithError(err)
			return "", fmt.Errorf("failed to write to processor: %w", err)
		}
	}

	if err := <-errChan; err != nil {
		return "", fmt.Errorf("s3 upload failed: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("image processor exited with error: %w", err)
	}

	uploadSuccessful = true
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

	objectId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	objectIdString := objectId.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, folder, fileInfo.GetIsActive())
	if err != nil{
		return "", errors.New("error creating new object")
	}

	pr, pw := io.Pipe()

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.mediaResources.S3Repository.UploadFileStream(ctx, pr, objectIdString, fileInfo.GetProjectId(), contentType, fileInfo.GetIsActive())
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