package s3

import (
	"audio-to-notion/shared/session"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	*s3.S3
}

func NewS3Client() (*S3Client, error) {
	s3Config := LoadConfig()
	session, err := session.NewSession(&session.SessionConfig{Endpoint: s3Config.Endpoint, Region: s3Config.Region})
	if err != nil {
		return nil, err
	}
	s3 := s3.New(session.Session, &aws.Config{S3ForcePathStyle: aws.Bool(s3Config.S3ForcePathStyle)})
	return &S3Client{s3}, nil
}

func (client *S3Client) GetFile(bucket string, key string) (io.ReadCloser, error) {
	output, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}
