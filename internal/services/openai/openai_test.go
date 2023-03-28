package openai

import (
	"context"
	"music-creator/internal/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSong(t *testing.T) {
	test.SetupTest()
	// mock response
	openAIService := NewOrderAIService(
		Config{
			Token: os.Getenv("OPENAI_TOKEN"),
		},
	)
	_, err := openAIService.CreateSong(context.Background(), 60, "Relax And Chill", "joy", "spanish", "reggae")
	assert.Nil(t, err)
	// assert.Equal(t, true, len(song) > 1)
}
