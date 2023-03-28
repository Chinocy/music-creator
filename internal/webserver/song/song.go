package song

import (
	"fmt"
	"music-creator/internal/services/openai"
	"music-creator/internal/util"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

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

type Body struct {
	BPM      int    `json:"bpm" validate:"required,min=20,max=200"`
	Emotion  string `json:"emotion" validate:"required,is-emotion"`
	Language string `json:"language" validate:"required,is-language"`
	Genre    string `json:"genre" validate:"required,is-genre"`
	Subject  string `json:"subject" validate:"required"`
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

	message := compact(fmt.Sprintf(`
		Hello! Create a complete song of which tempo is around %d BPM.
		It must be about %s. It musical genre is %s. 
		It must transmit %s, finally this must be in %s. 
		Please add a chord progression with the lyrics to it as well. 
		Don't forget to add the title.
	`, body.BPM, body.Subject, body.Genre, body.Emotion, body.Language))

	fmt.Println(message)

	content, err := ctrl.openAIService.CreateSong(c, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func (ctrl *SongCtrl) GetChoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"emotions":  util.Emotions,
		"genres":    util.Genres,
		"languages": util.Languages})
}

// GetFieldNameFromJSONTag custom function tag for fix get field name from json tag
func GetFieldNameFromJSONTag(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	return name
}

var spaceEater = regexp.MustCompile(`\s+`)

func compact(q string) string {
	q = strings.TrimSpace(q)
	return spaceEater.ReplaceAllString(q, " ")
}
