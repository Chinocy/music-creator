package openai

import (
	"context"
	"music-creator/internal/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloName(t *testing.T) {
	test.SetupTest()
	openAIService := NewOrderAIService(
		Config{
			Token: os.Getenv("OPENAI_TOKEN"),
		},
	)
	song, err := openAIService.CreateSong(context.Background(), "Create a song")
	assert.Nil(t, err)
	assert.Equal(t, true, len(song) > 1)

}
