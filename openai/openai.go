package openai

import (
	"github.com/mager/quotient/config"
	openai "github.com/sashabaranov/go-openai"
)

func ProvideOpenAI(cfg config.Config) *openai.Client {
	client := openai.NewClient(cfg.OpenAIKey)
	return client
}

var Options = ProvideOpenAI
