package main

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/gorilla/mux"
	"github.com/mager/quotient/config"
	"github.com/mager/quotient/handler"
	"github.com/mager/quotient/logger"
	"github.com/mager/quotient/openai"
	"github.com/mager/quotient/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			config.Options,
			logger.Options,
			router.Options,

			// APIs
			openai.Options,
		),
		fx.Invoke(
			Register,
		),
	).Run()
}

func Register(log *zap.SugaredLogger, router *mux.Router, openai gpt3.Client) {
	p := handler.Handler{
		Log:    log,
		Router: router,

		OpenAI: openai,
	}

	handler.New(p)
}
