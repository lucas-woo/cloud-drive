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
		fmt.Println(err)
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

	err = s.mediaResources.ProjectRepository.CreateNewObject(ctx, projectId, oId, folderId, req.IsActive, req.Format)
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

func (s *Service) UploadImageApiService(stream mediav1.MediaService_UploadImageApiServer) (objectID string, err error) {

	ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		return "", fmt.Errorf("failed to receive upload info: %w", err)
	}

	imageInfo := req.GetUploadInfo()
	if imageInfo == nil {
		return "", errors.New("missing upload info")
	}

	projectID, err := uuid.Parse(imageInfo.GetProjectId())
	if err != nil {
		return "", fmt.Errorf("invalid project id: %w", err)
	}

	folderID, err := uuid.Parse(imageInfo.GetFolderId())
	if err != nil {
		return "", fmt.Errorf("invalid folder id: %w", err)
	}

	objectUUID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate object id: %w", err)
	}

	objectID = objectUUID.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(
		ctx,
		projectID,
		objectUUID,
		folderID,
		imageInfo.GetIsActive(),
		imageInfo.GetFormat(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create new object: %w", err)
	}

	uploadSuccessful := false
	defer func() {
		if !uploadSuccessful {
			s.mediaResources.ProjectRepository.DeleteObjectWithOwnContext(objectUUID)
		}
	}()

	pipeReader, pipeWriter := io.Pipe()

	uploadErrChan := make(chan error, 1)

	go func() {
		uploadErrChan <- s.mediaResources.S3Repository.UploadStreamImage(
			ctx,
			pipeReader,
			objectID,
			imageInfo.GetProjectId(),
			imageInfo.GetContentType(),
			imageInfo.GetIsActive(),
			imageInfo.GetTransformations(),
		)
	}()

	for {
		req, err = stream.Recv()

		if err == io.EOF {
			if closeErr := pipeWriter.Close(); closeErr != nil {
				uploadErr := <-uploadErrChan

				if uploadErr != nil {
					return "", fmt.Errorf(
						"failed to close upload pipe: %w; s3 upload also failed: %v",
						closeErr,
						uploadErr,
					)
				}

				return "", fmt.Errorf("failed to close upload pipe: %w", closeErr)
			}

			break
		}

		if err != nil {
			pipeErr := pipeWriter.CloseWithError(err)

			uploadErr := <-uploadErrChan

			if pipeErr != nil {
				return "", fmt.Errorf(
					"stream receive error: %w; failed to close pipe: %v",
					err,
					pipeErr,
				)
			}

			if uploadErr != nil {
				return "", fmt.Errorf(
					"stream receive error: %w; s3 upload also failed: %v",
					err,
					uploadErr,
				)
			}

			return "", fmt.Errorf("stream receive error: %w", err)
		}

		chunk := req.GetImageChunk()
		if chunk == nil {
			chunkErr := errors.New("empty image chunk received")

			pipeErr := pipeWriter.CloseWithError(chunkErr)
			uploadErr := <-uploadErrChan

			if pipeErr != nil {
				return "", fmt.Errorf(
					"%w; failed to close pipe: %v",
					chunkErr,
					pipeErr,
				)
			}

			if uploadErr != nil {
				return "", fmt.Errorf(
					"%w; s3 upload also failed: %v",
					chunkErr,
					uploadErr,
				)
			}

			return "", chunkErr
		}

		if _, err = pipeWriter.Write(chunk); err != nil {
			pipeErr := pipeWriter.CloseWithError(err)
			uploadErr := <-uploadErrChan

			if pipeErr != nil {
				return "", fmt.Errorf(
					"failed to write image chunk: %w; failed to close pipe: %v",
					err,
					pipeErr,
				)
			}

			if uploadErr != nil {
				return "", fmt.Errorf(
					"failed to write image chunk: %w; s3 upload also failed: %v",
					err,
					uploadErr,
				)
			}

			return "", fmt.Errorf("failed to write image chunk: %w", err)
		}
	}

	if err := <-uploadErrChan; err != nil {
		return "", fmt.Errorf("s3 upload failed: %w", err)
	}

	uploadSuccessful = true

	return objectID, nil
}


func (s *Service) UploadFileApiService(stream mediav1.MediaService_UploadFileApiServer) (objectID string, err error) {

	ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		return "", fmt.Errorf("failed to receive upload info: %w", err)
	}

	fileInfo := req.GetUploadInfo()
	if fileInfo == nil {
		return "", errors.New("missing upload info")
	}

	projectID, err := uuid.Parse(fileInfo.GetProjectId())
	if err != nil {
		return "", fmt.Errorf("invalid project id: %w", err)
	}

	folderID, err := uuid.Parse(fileInfo.GetFolderId())
	if err != nil {
		return "", fmt.Errorf("invalid folder id: %w", err)
	}

	objectUUID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate object id: %w", err)
	}

	objectID = objectUUID.String()

	err = s.mediaResources.ProjectRepository.CreateNewObject(
		ctx,
		projectID,
		objectUUID,
		folderID,
		fileInfo.GetIsActive(),
		fileInfo.GetFormat(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create object: %w", err)
	}

	uploadSuccessful := false
	defer func() {
		if !uploadSuccessful {
			s.mediaResources.ProjectRepository.DeleteObjectWithOwnContext(objectUUID)
		}
	}()

	pr, pw := io.Pipe()

	uploadErrChan := make(chan error, 1)

	go func() {
		uploadErrChan <- s.mediaResources.S3Repository.UploadFileStream(
			ctx,
			pr,
			objectID,
			fileInfo.GetProjectId(),
			fileInfo.GetContentType(),
			fileInfo.GetIsActive(),
		)
	}()

	for {
		req, err = stream.Recv()

		if err == io.EOF {
			if closeErr := pw.Close(); closeErr != nil {
				uploadErr := <-uploadErrChan

				if uploadErr != nil {
					return "", fmt.Errorf(
						"failed to close upload pipe: %w; upload also failed: %v",
						closeErr,
						uploadErr,
					)
				}

				return "", fmt.Errorf(
					"failed to close upload pipe: %w",
					closeErr,
				)
			}

			break
		}


		if err != nil {
			_ = pw.CloseWithError(err)

			uploadErr := <-uploadErrChan

			if uploadErr != nil {
				return "", fmt.Errorf(
					"stream receive error: %w; upload also failed: %v",
					err,
					uploadErr,
				)
			}

			return "", fmt.Errorf(
				"stream receive error: %w",
				err,
			)
		}

		chunk := req.GetFileChunk()
		if chunk == nil {
			chunkErr := errors.New("missing file chunk")

			_ = pw.CloseWithError(chunkErr)

			uploadErr := <-uploadErrChan

			if uploadErr != nil {
				return "", fmt.Errorf(
					"%w; upload also failed: %v",
					chunkErr,
					uploadErr,
				)
			}

			return "", chunkErr
		}

		if _, err = pw.Write(chunk); err != nil {
			_ = pw.CloseWithError(err)

			uploadErr := <-uploadErrChan

			if uploadErr != nil {
				return "", fmt.Errorf(
					"failed to write file chunk: %w; upload also failed: %v",
					err,
					uploadErr,
				)
			}

			return "", fmt.Errorf(
				"failed to write file chunk: %w",
				err,
			)
		}
	}

	if err = <-uploadErrChan; err != nil {
		return "", fmt.Errorf("S3 upload failed: %w", err)
	}

	uploadSuccessful = true

	return objectID, nil
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

	if err != nil {
		return nil,nil,err
	}

	if len(assets) == 0 {
		return assets, nil, nil
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
	if err != nil || len(assets) == 0{
		return nil,nil, err;
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

	if err != nil || len(assets) == 0 {
		return nil, nil, err
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