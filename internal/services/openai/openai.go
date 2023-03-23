package openai

import (
	"context"

	goopenai "github.com/sashabaranov/go-openai"
)

type Config struct {
	Token string `json:"token"`
}

type OpenAIService struct {
	token  string
	client *goopenai.Client
}

func NewOrderAIService(cfg Config) *OpenAIService {
	return &OpenAIService{
		token:  cfg.Token,
		client: goopenai.NewClient(cfg.Token),
	}
}

func (o *OpenAIService) CreateSong(ctx context.Context, request string) (content string, err error) {
	resp, err := o.client.CreateChatCompletion(
		ctx,
		goopenai.ChatCompletionRequest{
			Model: goopenai.GPT3Dot5Turbo,
			Messages: []goopenai.ChatCompletionMessage{
				{
					Role:    goopenai.ChatMessageRoleUser,
					Content: request,
				},
			},
		},
	)

	if err != nil {
		return
	}
	content = resp.Choices[0].Message.Content
	return
}
