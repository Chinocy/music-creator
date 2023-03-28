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
	request := util.CompactString(fmt.Sprintf(`
	Hello! Create a complete song of which tempo is around %d BPM.
	It must be about %s. It musical genre is %s. 
	It must transmit %s and it must be in %s.
	Each line of the lyrics must be accompanied by its chord progression which must come before it excluding non-chord text.
	Each chords must be separated with a dash.
	Include title.
	Exclude non-song text.
	Every paragraph must be separated by 2 break lines and its lines must be separated by 1 break line.
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

		splitted := strings.Split(lines[0], ":")

		// is Title paragraph
		if strings.Contains(splitted[0], "Title") {
			if len(splitted) > 0 {
				title = strings.Replace(splitted[1], "\"", "", -1)
			} else {
				title = "Without title"
			}
			continue
		}

		pHeader := strings.Replace(splitted[0], "(", "", -1)
		pHeader = strings.Replace(pHeader, ")", "", -1)

		newP := Paragraph{
			Title:   pHeader,
			Phrases: []Phrase{},
		}

		lastChord := []string{}
		chordsAreFromHeader := false
		if len(splitted) > 1 && isChord(splitted[1]) {
			lastChord = getChords(splitted[1])
			chordsAreFromHeader = true
		}

		for i := 1; i < len(lines); i++ {
			if isChord(lines[i]) {
				if !chordsAreFromHeader && len(lastChord) > 0 {
					newP.Phrases = append(newP.Phrases, Phrase{
						Text:   "",
						Chords: lastChord,
					})
				}
				lastChord = getChords(lines[i])
			} else {
				newP.Phrases = append(newP.Phrases, Phrase{
					Text:   lines[i],
					Chords: lastChord,
				})
				lastChord = []string{}
			}
		}
		if !chordsAreFromHeader && len(lastChord) > 0 {
			newP.Phrases = append(newP.Phrases, Phrase{
				Text:   "",
				Chords: lastChord,
			})
		}
		ps = append(ps, newP)
	}
	return
}

func isChord(line string) bool {
	size := len(strings.Split(line, "-"))
	if size > 1 {
		return true
	}
	splitted := strings.Split(line, " ")
	for _, word := range splitted {
		if len(word) > 3 {
			return false
		}
	}
	return true
}

func getChords(line string) (chords []string) {
	chordText := strings.Replace(line, "(", " ", -1)
	chordText = strings.Replace(chordText, ")", " ", -1)
	chordText = strings.Replace(chordText, "|", " ", -1)
	chords = strings.Fields(strings.Replace(chordText, "-", " ", -1))
	return
}
