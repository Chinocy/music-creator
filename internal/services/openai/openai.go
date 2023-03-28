package openai

import (
	"context"
	"fmt"
	"music-creator/internal/util"
	"strings"

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

type Song struct {
	Title      string      `json:"title"`
	BPM        int         `json:"bpm"`
	Emotion    string      `json:"emotion"`
	Language   string      `json:"language"`
	Genre      string      `json:"genre"`
	Subject    string      `json:"subject"`
	Paragraphs []Paragraph `json:"paragraphs"`
	Response   string      `json:"response"`
	Request    string      `json:"request"`
}

type Paragraph struct {
	Title   string   `json:"title"`
	Phrases []Phrase `json:"phrases"`
}

type Phrase struct {
	Text   string   `json:"text"`
	Chords []string `json:"chords"`
}

func (o *OpenAIService) CreateSong(ctx context.Context,
	bpm int, subject, emotion, language, genre string,
) (song Song, err error) {
	// Please add a chord progression with the lyrics in every line to it as well.
	// All chord progressions are without parenthesis and must be separated with an white space.

	request := util.CompactString(fmt.Sprintf(`
	Hello! Create a complete song of which tempo is around %d BPM.
	It must be about %s. It musical genre is %s. 
	It must transmit %s and it must be in %s.
	Include title. Exclude intro. 
	Exclude non-song text.
	Every paragraph must be separated by 2 break lines and its lines must be separated by 1 break line.
	Each line of the lyric must be accompanied by its chord progression.
	All chords must be separated with a dash.
	Repeat paragraphs or stanzas or verses or phrases as much as needed, don't use multiplication.
	`, bpm, subject, genre, emotion, language))

	fmt.Println(request)

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

	response := resp.Choices[0].Message.Content
	title, paragraphs := GetTitleAndParagraphs(response)
	song = Song{
		BPM:        bpm,
		Subject:    subject,
		Genre:      genre,
		Emotion:    emotion,
		Language:   language,
		Paragraphs: paragraphs,
		Title:      title,
		Response:   response,
		Request:    request,
	}
	return
}

func GetTitleAndParagraphs(content string) (title string, ps []Paragraph) {
	paragraphs := strings.Split(content, "\n\n")
	for _, p := range paragraphs {
		lines := strings.Split(p, "\n")
		if strings.Contains(lines[0], "Title:") {
			title = strings.Split(lines[0], ": ")[1]
			continue
		}

		newP := Paragraph{
			Title:   lines[0],
			Phrases: []Phrase{},
		}

		for i := 1; i+1 < len(lines); i = i + 2 {
			chords := strings.Replace(lines[i], "-", " ", -1)
			phrase := Phrase{
				Text:   lines[i+1],
				Chords: strings.Fields(chords),
			}
			newP.Phrases = append(newP.Phrases, phrase)
		}

		ps = append(ps, newP)
	}
	return
}
