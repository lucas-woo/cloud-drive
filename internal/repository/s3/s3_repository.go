package s3repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lucas-woo/cloud-drive/internal/config"
)

type S3Repository struct {
	s3Client *s3.Client
	presignClient *s3.PresignClient
	bucketName string
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

func NewS3Repository(	s3Client *s3.Client, presignClient *s3.PresignClient, bucketName string) *S3Repository {
	return &S3Repository{
		s3Client: s3Client,
		presignClient: presignClient,
		bucketName: bucketName,
	}
}