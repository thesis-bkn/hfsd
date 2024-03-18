package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port        int    `env:"PORT"         envDefault:"2002"`
	DatabaseURL string `env:"DATABASE_URL"`
}

func LoadConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
