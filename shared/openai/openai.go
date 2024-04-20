package openai

import (
	"context"
	"encoding/json"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAi struct {
	client *openai.Client
}

type Story struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Completion struct {
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	MainPoints    []string `json:"main_points"`
	ActionItems   []string `json:"action_items"`
	FollowUp      []string `json:"follow_up"`
	Arguments     []string `json:"arguments"`
	RelatedTopics []string `json:"related_topics"`
	Sentiment     string   `json:"sentiment"`
}

func NewOpenAiClient() (*OpenAi, error) {
	client := openai.NewClient(os.Getenv("OPEN_AI_KEY"))
	return &OpenAi{
		client: client,
	}, nil
}

func (ai *OpenAi) CreateTranscription(ctx context.Context, filePath string) (string, error) {
	request := openai.AudioRequest{
		Model:       openai.Whisper1,
		FilePath:    filePath,
		Temperature: 0.2,
	}
	response, err := ai.client.CreateTranscription(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Text, nil
}

func (ai *OpenAi) CreateCompletion(ctx context.Context, transcript string) (string, error) {
	defaultPrompt := `Please analyze the provided transcript and generate a JSON file with the following structure:

	"title": Assign a descriptive title based on the transcript content.
	"summary": Summarize the transcript in a concise paragraph.
	"main_points": List up to 10 main points extracted from the transcript. Limit each point to 100 words.
	"action_items": Identify up to 5 actionable steps suggested by the transcript. Limit each action to 100 words.
	"follow_up": Formulate up to 5 follow-up questions that could deepen understanding or inquiry based on the transcript. Limit each question to 100 words.
	"arguments": Outline up to 5 potential arguments against the points made in the transcript. Limit each argument to 100 words.
	"related_topics": List up to 5 topics related to the transcript content. Limit each topic to 100 words.
	"sentiment": Perform a sentiment analysis and summarize the overall tone of the transcript.

	Ensure proper JSON formatting: Do not include trailing commas in any arrays. Please see the transcript below:`

	defaultPrompt += transcript

	request := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser, Content: defaultPrompt,
			},
			{
				Role: openai.ChatMessageRoleSystem, Content: `You are an assistant that only speaks JSON. Do not write normal text.

				Example formatting: {
						"title": "Notion Buttons",
						"summary": "A collection of buttons for Notion",
						"action_items": [
						  "item 1",
						  "item 2",
						  "item 3"
						],
						"follow_up": [
						  "item 1",
						  "item 2",
						  "item 3"
						],
						"arguments": [
						  "item 1",
						  "item 2",
						  "item 3"
						],
						"related_topics": [
						  "item 1",
						  "item 2",
						  "item 3"
						]
				"sentiment": "positive"
					  }
				  `,
			},
		},
		MaxTokens:   2000,
		Temperature: 0.2,
	}
	response, err := ai.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

func (ai *OpenAi) AudioToCompletion(ctx context.Context, f string) (Completion, error) {
	completion := Completion{}

	text, err := ai.CreateTranscription(ctx, f)
	if err != nil {
		return completion, err
	}
	completionText, err := ai.CreateCompletion(ctx, text)
	if err != nil {
		return completion, err
	}
	err = json.Unmarshal([]byte(completionText), &completion)
	if err != nil {
		return completion, err
	}
	return completion, nil
}
