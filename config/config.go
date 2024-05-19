package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rudyeila/go-bring-api/bring"
	"github.com/rudyeila/hello-fresh-go-client/hellofresh"
)

type Config struct {
	HelloFresh hellofresh.Config
	Bring      bring.Config
}

func New(filenames ...string) (Config, error) {
	conf := Config{}
	_ = godotenv.Load(filenames...)

	err := envconfig.Process("", &conf)
	if err != nil {
		return conf, fmt.Errorf("process envconfig :%w", err)
	}

	return conf, nil
}
