package s3repository

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	S3Client *s3.Client
	PresignClient *s3.PresignClient
	BucketName string
}


func (r *S3Repository) GetSignedUploadUrl() {

}

func NewS3Repository(	s3Client *s3.Client, presignClient *s3.PresignClient, bucketName string) *S3Repository {
	return &S3Repository{
		S3Client: s3Client,
		PresignClient: presignClient,
		BucketName: bucketName,
	}
}