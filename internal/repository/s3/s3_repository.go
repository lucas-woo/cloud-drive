package s3repository

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lucas-woo/cloud-drive/internal/config"
)

type S3Repository struct {
	s3Client *s3.Client
	presignClient *s3.PresignClient
	bucketName string
	uploader *transfermanager.Client
}


func (r *S3Repository) GetPreSignedUploadUrl(ctx context.Context, objectId string) (string, error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(objectId), 
	}

	req, err := r.presignClient.PresignPutObject(ctx, params, func(opts *s3.PresignOptions) {
		opts.Expires = config.PreSignedUrlTime
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (r *S3Repository) UploadStreamImage(ctx context.Context, reader io.ReadCloser, objectId, contentType string) error {	
	input := &transfermanager.UploadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(objectId),
		Body: reader,
		ContentType: aws.String(contentType),
	}

	_, err := r.uploader.UploadObject(ctx, input)

	if err != nil {
		return err
	}

	return nil	
}

func (r *S3Repository) UploadFileStream(ctx context.Context, reader io.Reader, objectId, contentType string) error {
	
	input := &transfermanager.UploadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(objectId),
		Body: reader,
		ContentType: aws.String(contentType),
	}

	_, err := r.uploader.UploadObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}


func (r *S3Repository) DeleteObject(ctx context.Context, objectName string) error {
	_, err := r.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(objectName),
	})
	if err != nil {
		return err
	}

	return nil
}

func NewS3Repository(	s3Client *s3.Client, presignClient *s3.PresignClient, bucketName string, uploader *transfermanager.Client) *S3Repository {
	return &S3Repository{
		uploader: uploader,
		s3Client: s3Client,
		presignClient: presignClient,
		bucketName: bucketName,
	}
}