// Package config with environment variables
package config

import (
	env "github.com/caarlos0/env/v8"
)

// Config is a struct with environment variables
type Config struct {
	ConnectionString string `env:"APOD_CONNECTION_STRING"`
	Port             string `env:"APOD_PORT"`
	NasaAPIKey       string `env:"NASA_API_KEY"`
	NasaAPIURL       string `env:"NASA_API_URL"`
	ApodDB           string `env:"APOD_DB"`
	ApodUser         string `env:"APOD_USER"`
	ApodPassword     string `env:"APOD_PASSWORD"`
	ApodDBPort       string `env:"APOD_DB_PORT"`
}

// New returns parsed object of config
func New() (*Config, error) {
	cfg := new(Config)
	err := env.Parse(cfg)
	return cfg, err
}
