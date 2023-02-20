package openai

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/mager/quotient/config"
)

func ProvideOpenAI(cfg config.Config) gpt3.Client {
	return gpt3.NewClient(cfg.OpenAIKey)
}

var Options = ProvideOpenAI
