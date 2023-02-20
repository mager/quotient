package openai

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/mager/quotient/config"
)

func ProvideOpenAI(cfg config.Config) gpt3.Client {
	opts := []gpt3.ClientOption{
		gpt3.WithDefaultEngine("text-davinci-003"),
	}
	return gpt3.NewClient(cfg.OpenAIKey, opts...)
}

var Options = ProvideOpenAI
