package song

import (
	"music-creator/internal/services/openai"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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
	}
	content, err := ctrl.openAIService.CreateSong(c, params.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}
