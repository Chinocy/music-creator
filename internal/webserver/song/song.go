package song

import (
	"context"
	"music-creator/internal/services/openai"
	"music-creator/internal/util"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type OpenAIService interface {
	CreateSong(ctx context.Context, bpm int, subject, emotion, language, genre string) (song openai.Song, err error)
}

type SongCtrl struct {
	openAIService OpenAIService
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

type Body struct {
	BPM      int    `json:"bpm" validate:"required,min=60,max=180"`
	Emotion  string `json:"emotion" validate:"required,is-emotion"`
	Language string `json:"language" validate:"required,is-language"`
	Genre    string `json:"genre" validate:"required,is-genre"`
	Subject  string `json:"subject" validate:"require,is-subject"`
}

func (ctrl *SongCtrl) CreateSong(c *gin.Context) {
	body := new(Body)
	err := c.Bind(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	clientErrs := util.ValidateStruct(body)
	if len(clientErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": clientErrs})
		return
	}

	song, err := ctrl.openAIService.CreateSong(c,
		body.BPM,
		body.Subject,
		body.Emotion,
		body.Language,
		body.Genre,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"song": song})
}

func (ctrl *SongCtrl) GetChoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"emotions":  util.Emotions,
		"genres":    util.Genres,
		"languages": util.Languages,
		"bpm":       util.BPMChoices,
		"subjects":  util.Subjects,
	})
}

// GetFieldNameFromJSONTag custom function tag for fix get field name from json tag
func GetFieldNameFromJSONTag(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	return name
}
