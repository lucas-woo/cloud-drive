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
	assetsmodels "github.com/lucas-woo/cloud-drive/internal/models/assets"
	projectmodels "github.com/lucas-woo/cloud-drive/internal/models/project"
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

	pId, err := uuid.Parse(projectId)

	if err != nil {
		return
	}
	_, err = s.mediaResources.ProjectRepository.CreateRootFolder(ctx, pId)

	return
}

func (s *Service) GetUploadObjectSignedUrl(ctx context.Context, req *dto.UploadObjectRequest) (url string, objectId string, err error) {
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil{
		return 
	}
	oId, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}

	objectId = oId.String()

	folderId, err := uuid.Parse(req.FolderId)
	if err != nil {
		return "", "", err
	}

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, oId, folderId, req.IsActive)
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

	updated, err := s.mediaResources.ProjectRepository.ConfirmObjectInfo(ctx, objectId, req.FileSize)

	if err != nil {
		return err
	}

	if !updated {
		return nil
	}

	projectId, err := s.mediaResources.ProjectRepository.GetProjectIdFromObjectId(ctx, objectId)

	if err != nil {
		return err
	}

	err = s.mediaResources.ProjectRepository.IncrementProjectAssetsCount(ctx, projectId, req.FileSize)

	return err 
}

func (s *Service) GetDashboard(ctx context.Context, req *dto.GetDashboardRequest) (*projectmodels.ProjectModel, error) {
	userId, err := uuid.Parse(req.UserId)

	if err != nil {
		return nil, err
	}

	projectInfo, err := s.mediaResources.ProjectRepository.GetOneProjectByUserId(ctx, userId)

	if err != nil {
		return nil, err
	}		

	return projectInfo, nil
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

	folderId, err := uuid.Parse(imageInfo.FolderId)

	if err != nil {
		return "", err
	}

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, folderId, imageInfo.GetIsActive())

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



	folderId, err := uuid.Parse(fileInfo.GetFolderId())

	if err != nil {
		return "", err
	}

	contentType := fileInfo.GetContentType()

	objectId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	objectIdString := objectId.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, objectId, folderId, fileInfo.GetIsActive())
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

//change config.AmountImagesToFetch to change number of assets returned
func (s *Service) GetAssets(ctx context.Context, req *dto.GetAssetsRequest) ([]*dto.ProjectObject, *dto.AssetCursor, error) {

	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, nil, err
	}
	var assets []*dto.ProjectObject
	if req.AssetCursor != nil {
		lastObjectId, err := uuid.Parse(req.AssetCursor.ObjectId)
		if err != nil {
			return nil,nil, err
		}
		assets, err = s.mediaResources.ProjectRepository.GetAssets(ctx, projectId, &dto.AssetCursor{
			ObjectId: lastObjectId,
			CreatedAt: req.AssetCursor.CreatedAt,
		}, config.AmountImagesToFetch)
	} else {
		assets, err = s.mediaResources.ProjectRepository.GetAssets(ctx, projectId, nil, config.AmountImagesToFetch)
	}
	nextCursor := &dto.AssetCursor{
			CreatedAt: assets[len(assets)-1].CreatedAt,
			ObjectId:  assets[len(assets)-1].ObjectId,
	}
	return assets, nextCursor, err
}

func (s *Service) GetAllFolders(ctx context.Context, pId string) ([]*dto.ProjectFolder, error) {

	projectId, err := uuid.Parse(pId)
	if err != nil {
		return nil, err
	}

	return s.mediaResources.ProjectRepository.GetAllFolders(ctx, projectId)
}

func (s *Service) GetAllCollections(ctx context.Context, pId string) ([]*assetsmodels.CollectionModel, error ) {
	projectId, err := uuid.Parse(pId)
	if err != nil {
		return nil, err
	}

	return s.mediaResources.ProjectRepository.GetAllCollections(ctx, projectId)
}


func (s *Service) GetAssetsInFolder(ctx context.Context, req *dto.GetAssetsInFolderRequest) ([]*dto.ProjectObject, *dto.AssetCursor, error) {

	projectId, err := uuid.Parse(req.ProjectId)

	if err != nil {
		return nil, nil, err
	}

	folderId, err := uuid.Parse(req.FolderId)
	if err != nil {
		return nil, nil, err
	}

	var assets []*dto.ProjectObject

	if req.AssetCursor != nil {
		lastObjectId, err := uuid.Parse(req.AssetCursor.ObjectId)
		if err != nil {
			return nil,nil, err
		}
		assets, err = s.mediaResources.ProjectRepository.GetAssetsInFolder(ctx, projectId, folderId, &dto.AssetCursor{
			ObjectId: lastObjectId,
			CreatedAt: req.AssetCursor.CreatedAt,
		}, config.AmountImagesToFetch)
	} else {
		assets, err = s.mediaResources.ProjectRepository.GetAssetsInFolder(ctx, projectId, folderId, nil, config.AmountImagesToFetch)
	}	

	nextCursor := &dto.AssetCursor{
			CreatedAt: assets[len(assets)-1].CreatedAt,
			ObjectId:  assets[len(assets)-1].ObjectId,
	}
	return assets, nextCursor, err	
}

func (s *Service) GetAssetsInCollection(ctx context.Context, req *dto.GetAssetsInCollectionRequest) ([]*dto.ProjectObject, *dto.AssetCursor, error) {

	projectId, err := uuid.Parse(req.ProjectId)

	if err != nil {
		return nil, nil, err
	}

	collectionId, err := uuid.Parse(req.CollectionId)
	if err != nil {
		return nil, nil, err
	}	

	var assets []*dto.ProjectObject

	if req.AssetCursor != nil {
		lastObjectId, err := uuid.Parse(req.AssetCursor.ObjectId)
		if err != nil {
			return nil,nil, err
		}
		assets, err = s.mediaResources.ProjectRepository.GetAssetsInCollection(ctx, projectId, collectionId, &dto.AssetCursor{
			ObjectId: lastObjectId,
			CreatedAt: req.AssetCursor.CreatedAt,
		}, config.AmountImagesToFetch)
	} else {
		assets, err = s.mediaResources.ProjectRepository.GetAssetsInCollection(ctx, projectId, collectionId, nil, config.AmountImagesToFetch)
	}	

	nextCursor := &dto.AssetCursor{
			CreatedAt: assets[len(assets)-1].CreatedAt,
			ObjectId:  assets[len(assets)-1].ObjectId,
	}
	return assets, nextCursor, err	
}

func(s *Service) CreateNewFolder(ctx context.Context, req *dto.CreateNewFolderRequest) (*dto.CreateNewFolderResponse, error) {
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	folderId, err := s.mediaResources.ProjectRepository.CreateNewFolder(ctx, projectId, req.FolderName)
	if err != nil {
		return nil, err
	}
	
	return &dto.CreateNewFolderResponse{FolderId: folderId}, nil
}

func(s *Service) CreateNewCollection(ctx context.Context, req *dto.CreateNewCollectionRequest) (*dto.CreateNewCollectionResponse, error) {
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	creatorId, err := uuid.Parse(req.CreatorId)
	if err != nil {
		return nil, err
	}

	collectionId, err := s.mediaResources.ProjectRepository.CreateNewCollection(ctx, projectId, creatorId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &dto.CreateNewCollectionResponse{CollectionId: collectionId}, err
}

func NewMediaService(mediaResources *database.MediaResources) *Service {
	return &Service{
		mediaResources: mediaResources,
	}
} 