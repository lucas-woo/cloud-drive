package s3repository

import (
	"context"
	"fmt"
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

type UnseekableReader struct {
	R io.ReadCloser
}

func (u UnseekableReader) Read(p []byte) (int, error) {
	return u.R.Read(p)
}

func (u UnseekableReader) Close() error {
	return u.R.Close()
}

func (r *S3Repository) GetPreSignedUploadUrl(ctx context.Context, objectId, projectId string, isActive bool) (string, error) {
	prefix := config.S3PrivatePrefix
	if isActive {
		prefix = config.S3PublicPrefix
	}

	key := fmt.Sprintf("%s/%s/%s", prefix, projectId, objectId)

	params := &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	}

	req, err := r.presignClient.PresignPutObject(ctx, params, func(opts *s3.PresignOptions) {
		opts.Expires = config.PreSignedUrlTime
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (r *S3Repository) UploadStreamImage(ctx context.Context, reader io.ReadCloser, objectId, projectId, contentType string, isActive bool) error {	

	prefix := config.S3PrivatePrefix
	if isActive {
		prefix = config.S3PublicPrefix
	}

	key := fmt.Sprintf("%s/%s/%s", prefix, projectId, objectId)

	input := &transfermanager.UploadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(key),
		Body: reader,
		ContentType: aws.String(contentType),
	}

	_, err := r.uploader.UploadObject(ctx, input)

	if err != nil {
		return err
	}

	return nil	
}

func (r *S3Repository) UploadFileStream(ctx context.Context, reader io.Reader, objectId, projectId, contentType string, isActive bool) error {
	
	prefix := config.S3PrivatePrefix
	if isActive {
		prefix = config.S3PublicPrefix
	}


	key := fmt.Sprintf("%s/%s/%s", prefix, projectId, objectId)
	
	input := &transfermanager.UploadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(key),
		Body: reader,
		ContentType: aws.String(contentType),
	}

	_, err := r.uploader.UploadObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (r *S3Repository) ActivateObject(ctx context.Context, objectId string) (error) {
	_, err := r.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(r.bucketName),
			CopySource: aws.String(r.bucketName + "/private/" + objectId),
			Key:        aws.String("public/" + objectId),
	})
	if err != nil {
			return err
	}

	_, err = r.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(r.bucketName),
			Key:    aws.String("public/" + objectId),
	})
	if err != nil {
			return err
	}
	return nil
}	


func (r *S3Repository) DisactivateObject(ctx context.Context, objectId string) (error) {
	_, err := r.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(r.bucketName),
			CopySource: aws.String(r.bucketName + "/public/" + objectId),
			Key:        aws.String("private/" + objectId),
	})
	if err != nil {
			return err
	}

	_, err = r.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(r.bucketName),
			Key:    aws.String("public/" + objectId),
	})
	if err != nil {
			return err
	}
	return nil
} 



func (r *S3Repository) DeleteObject(ctx context.Context, objectPath string) error {
	_, err := r.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(objectPath),
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