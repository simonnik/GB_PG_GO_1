package config

import (
	"github.com/kelseyhightower/envconfig"
)

type DB struct {
	Host     string `envconfig:"DB_HOST"`
	Port     string `envconfig:"DB_PORT"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Name     string `envconfig:"DB_NAME"`
}

// Config contains application settings
type Config struct {
	DB   DB
	Port string `envconfig:"APP_PORT"`
}

// BuildConfig creates a configuration structure using environment variables
func BuildConfig() (*Config, error) {
	var c Config

	err := envconfig.Process("", &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
