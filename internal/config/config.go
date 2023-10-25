// Package config with environment variables
package config

import (
	env "github.com/caarlos0/env/v8"
)

// Config is a struct with environment variables
type Config struct {
	ConnectionString string `env:"APOD_CONNECTION_STRING"`
	Port             int    `env:"APOD_PORT"`
	NasaApiKey       string `env:"NASA_API_KEY"`
	NasaApiUrl       string `env:"NASA_API_URL"`
}

// New returns parsed object of config
func New() (*Config, error) {
	cfg := new(Config)
	err := env.Parse(cfg)
	return cfg, err
}
