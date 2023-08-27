package translator

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type Translator struct {
	client *openai.Client
}

type Language string

func (l *Language) String() string {
	return string(*l)
}

func (l *Language) GetAntipode() Language {
	if *l == LanguageEnUs {
		return LanguageRuRu
	} else {
		return LanguageEnUs
	}
}

func (l *Language) Parse(s string) error {
	switch s {
	case "en-us":
		*l = LanguageEnUs
	case "ru-ru":
		*l = LanguageRuRu
	default:
		return fmt.Errorf("unsupported language %s", s)
	}

	return nil
}

const (
	LanguageEnUs Language = "en-us"
	LanguageRuRu Language = "ru-ru"
)

func NewTranslator(token string) *Translator {
	return &Translator{client: openai.NewClient(token)}
}

func (t *Translator) Translate(text string, from, to Language) ([]byte, error) {
	resp, err := t.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf("Переведи c %s на %s следующий текст", from, to),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return []byte(resp.Choices[0].Message.Content), nil
}
