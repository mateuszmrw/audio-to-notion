package s3

import (
	"os"
)

type S3Config struct {
	Endpoint         string
	Region           string
	S3ForcePathStyle bool
}

func LoadConfig() *S3Config {
	return &S3Config{
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: os.Getenv("AWS_ENDPOINT"),
	}
}
