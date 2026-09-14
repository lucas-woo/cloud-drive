package services

import (
	"context"
	"io"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
	"github.com/lucas-woo/cloud-drive/internal/utils"
)

type MediaApiService struct {
	mediaClient mediav1.MediaServiceClient
	iamClient iamv1.IAMServiceClient
	mapper utils.RestMapper
}

func (s *MediaApiService) UploadFile(
	ctx context.Context,
	projectId string,
	objectName string,
	originalFileName string,
	folderId string,
	isActive bool,
	contentType string,
	file io.Reader,
) (*api.UploadObjectApiResponse, error) {

	format := s.mapper.FindFormat(originalFileName)

	stream, err := s.mediaClient.UploadFileApi(ctx)
	if err != nil {
		return nil, err
	}

	err = stream.Send(&mediav1.UploadFileApiRequest{
		Payload: &mediav1.UploadFileApiRequest_UploadInfo{
			UploadInfo: &mediav1.FileUploadInfo{
				ProjectId:   projectId,
				ObjectName:  objectName,
				FolderId:    folderId,
				ContentType: contentType,
				Format:      format,
				IsActive:    isActive,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 64*1024)

	for {
		n, err := file.Read(buf)

		if n > 0 {
			if err := stream.Send(&mediav1.UploadFileApiRequest{
				Payload: &mediav1.UploadFileApiRequest_FileChunk{
					FileChunk: buf[:n],
				},
			}); err != nil {
				return nil, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return &api.UploadObjectApiResponse{
		ObjectId: res.GetObjectId(),
	}, nil
}

func (s *MediaApiService) UploadImage(
	ctx context.Context,
	projectId string,
	objectName string,
	originalFileName string,
	folderId string,
	isActive bool,
	contentType string,
	transformations *api.ImageTransformations,
	image io.Reader,
) (*api.UploadImageApiResponse, error) {

	format := s.mapper.FindFormat(originalFileName)

	stream, err := s.mediaClient.UploadImageApi(ctx)
	if err != nil {
		return nil, err
	}

	protoTransformations := s.mapper.ToProtoTransformations(transformations)

	err = stream.Send(&mediav1.UploadImageApiRequest{
		Payload: &mediav1.UploadImageApiRequest_UploadInfo{
			UploadInfo: &mediav1.ImageUploadInfo{
				ProjectId:       projectId,
				ObjectName:      objectName,
				FolderId:        folderId,
				ContentType:     contentType,
				Format:          format,
				IsActive:        isActive,
				Transformations: protoTransformations,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 64*1024)

	for {
		n, err := image.Read(buf)

		if n > 0 {
			if err := stream.Send(&mediav1.UploadImageApiRequest{
				Payload: &mediav1.UploadImageApiRequest_ImageChunk{
					ImageChunk: buf[:n],
				},
			}); err != nil {
				return nil, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return &api.UploadImageApiResponse{
		ObjectId: res.GetObjectId(),
	}, nil
}

func (s *MediaApiService) ValidateApiKey(ctx context.Context, apiKey, apiSecret string, permission iamv1.ValidateApiKeyPermissionRequest_Permission) (bool, error) {
	res, err := s.iamClient.ValidateApiKeyPermission(ctx, &iamv1.ValidateApiKeyPermissionRequest{
		Permission: permission,
		ApiKey: apiKey,
		ApiSecret: apiSecret,
	})
	return res.GetAuthorized(), err
}

func (s *MediaApiService) GetProjectId(ctx context.Context, apiKey string) (string, error) {
	res, err := s.iamClient.GetProjectId(ctx, &iamv1.GetProjectIdRequest{
		ApiKey: apiKey,
	})
	if err != nil {
		return "", err
	}
	return res.GetProjectId(), err
}

func (s *MediaApiService) GetAllFolders(ctx context.Context, projectId string) {
	
}

func NewMediaApiService(mediaClient mediav1.MediaServiceClient, iamClient iamv1.IAMServiceClient, mapper utils.RestMapper ) *MediaApiService {
	return &MediaApiService{
		mediaClient: mediaClient,
		iamClient: iamClient,
		mapper: mapper,
	}
}