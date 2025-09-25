package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host        string `envconfig:"HOST" default:"localhost"`
	Port        string `envconfig:"PORT" default:"8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgres://postgres:31313758@localhost:5432/taskdb?sslmode=disable"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
	JwtSecret   string `envconfig:"JWT_SECRET" default:"secret"`
	JwtTtlMin   int    `envconfig:"JWT_TTL_MINUTES" default:"60"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
