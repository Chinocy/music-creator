package song

import (
	"music-creator/internal/services/openai"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var emotions = []string{
	"amusement", "joy", "eroticism", "euphoric", "relaxation", "sadness",
	"dreaminess", "triumph", "anxiety", "scariness", "annoyance", "defiance",
}

var genres = []string{
	"rock", "pop", "hip hop", "reggae", "salsa", "merengue", "bachata", "funk", "soul", "reggaeton",
}

var languages = []string{
	"english", "spanish",
}

// bpm 20-200

type SongCtrl struct {
	openAIService *openai.OpenAIService
}

func NewSongCtrl() SongCtrl {
	openAIService := openai.NewOrderAIService(
		openai.Config{
			Token: os.Getenv("OPENAI_TOKEN"),
		},
	)
	ctrl := SongCtrl{
		openAIService: openAIService,
	}
	return ctrl
}

type Params struct {
	Message string
}

func (ctrl *SongCtrl) CreateSong(c *gin.Context) {
	params := new(Params)
	err := c.Bind(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content, err := ctrl.openAIService.CreateSong(c, params.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func (ctrl *SongCtrl) GetChoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"emotions":  emotions,
		"genres":    genres,
		"languages": languages})
}
