package main

import (
	"audio-to-notion/shared/notion"
	openai "audio-to-notion/shared/openai"
	"audio-to-notion/shared/s3"
	"context"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const AUDIO_FILE_PATH = "/tmp/audio.mp3"

func HandleRequest(ctx context.Context, event events.S3Event) (string, error) {
	record := event.Records[0]

	ai, err := openai.NewOpenAiClient()
	if err != nil {
		return "Ai Config Error", err
	}
	s3Client, err := s3.NewS3Client()
	if err != nil {
		return "S3 Client Error", err
	}
	body, err := s3Client.GetFile(record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		return "", err
	}
	outFile, err := os.Create(AUDIO_FILE_PATH)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, body)
	if err != nil {
		return "", err
	}
	completion, err := ai.AudioToCompletion(ctx, AUDIO_FILE_PATH)
	if err != nil {
		return "", err
	}

	notionClient := notion.NewNotionClient()

	page, err := notionClient.CreatePage(ctx, os.Getenv("NOTION_DATABASE_ID"), completion)
	if err != nil {
		return "", err
	}

	return page.ID, nil
}

func main() {
	lambda.Start(HandleRequest)
}
