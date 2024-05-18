package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BaseURL        string        `envconfig:"HF_BASE_URL" default:"https://www.hellofresh.de/gw/api"`
	AccessToken    string        `envconfig:"HF_ACCESS_TOKEN"`
	DefaultTimeout time.Duration `default:"10s"`
}

func New(filenames ...string) (Config, error) {
	conf := Config{}
	_ = godotenv.Load(filenames...) //nolint:errcheck // ignore error

	err := envconfig.Process("", &conf)
	if err != nil {
		return conf, fmt.Errorf("process envconfig :%w", err)
	}

	return conf, nil
}
