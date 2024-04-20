package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SessionConfig struct {
	Endpoint string
	Region   string
}

type AwsSession struct {
	Session *session.Session
}

func NewSession(sessionConfig *SessionConfig) (*AwsSession, error) {
	session, err := session.NewSession(&aws.Config{
		Region:   aws.String(sessionConfig.Region),
		Endpoint: aws.String(sessionConfig.Endpoint),
	})
	if err != nil {
		return nil, err
	}
	return &AwsSession{
		Session: session,
	}, nil
}
