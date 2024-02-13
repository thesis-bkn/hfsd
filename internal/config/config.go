package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port int    `env:"PORT"         envDefault:"2002"`
	URI  string `env:"DATABASE_URL"`

	Authenticate struct {
		JwtSecret string `env:"JWT_SECRET"`
	} `envPrefix:"AUTHENTICATE_"`
}

func LoadConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
