package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port               int    `env:"PORT"                  envDefault:"2002"`
	DatabaseURL        string `env:"DATABASE_URL"`
	Bucket             string `env:"BUCKET_NAME"`
	AwsAccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	EndpointUrl        string `env:"S3_ENDPOINT_URL"`
	ImagePath          string `env:"IMAGE_PATH"`
	MaskPath           string `env:"MASK_PATH"`
}

func LoadConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	cfg.ImagePath = "images"
	cfg.MaskPath = "masks"

	return &cfg
}
