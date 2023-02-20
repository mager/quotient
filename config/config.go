package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	OpenAIKey string
}

func ProvideConfig(log *zap.SugaredLogger) Config {
	var cfg Config
	err := envconfig.Process("quotient", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	return cfg
}

var Options = ProvideConfig
