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
	openAIService := NewOrderAIService(
		Config{
			Token: os.Getenv("OPENAI_TOKEN"),
		},
	)
	song, err := openAIService.CreateSong(context.Background(), 60, "Relax And Chill", "joy", "spanish", "reggae")
	assert.Nil(t, err)
	assert.Equal(t, true, len(song.Paragraphs) > 1)
	assert.Equal(t, true, song.Title != "")
}
