package main

import (
	"github.com/gorilla/mux"
	"github.com/mager/quotient/handler"
	"github.com/mager/quotient/logger"
	"github.com/mager/quotient/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			logger.Options,
			router.Options,
		),
		fx.Invoke(
			handler.New,
		),
	).Run()
}

func Register(router *mux.Router, log *zap.SugaredLogger) {
	p := handler.Handler{
		Router: router,
		Log:    log,
	}

	handler.New(p)
}
