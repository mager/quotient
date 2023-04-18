package main

import (
	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/mager/quotient/config"
	"github.com/mager/quotient/database"
	"github.com/mager/quotient/handler"
	"github.com/mager/quotient/logger"
	openaiClient "github.com/mager/quotient/openai"
	"github.com/mager/quotient/router"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			config.Options,
			database.Options,
			logger.Options,
			router.Options,

			// APIs
			openaiClient.Options,
		),
		fx.Invoke(
			Register,
		),
	).Run()
}

func Register(db *firestore.Client, log *zap.SugaredLogger, router *mux.Router, openai *openai.Client) {
	p := handler.Handler{
		Database: db,
		Log:      log,
		Router:   router,

		OpenAI: openai,
	}

	handler.New(p)
}
