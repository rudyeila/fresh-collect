package hellofresh

import "time"

type Config struct {
	BaseURL        string        `envconfig:"HF_BASE_URL" default:"https://www.hellofresh.de/gw/api"`
	AccessToken    string        `envconfig:"HF_ACCESS_TOKEN"`
	DefaultTimeout time.Duration `default:"10s"`
}
