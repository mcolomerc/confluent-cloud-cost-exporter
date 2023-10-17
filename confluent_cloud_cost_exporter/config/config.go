package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Endpoints   `yaml:"endpoints"`
	Credentials `yaml:"credentials"`
	Cache       `yaml:"cache"`
	PromConfig  `yaml:"prometheus"`
}

type Cache struct {
	Expiration time.Duration `yaml:"expiration" env:"CACHE_EXPIRATION" env-default:"60"`
}

// Create a enum for the different types of exporters
type Format string

const (
	PROMETHEUS Format = "prometheus"
	JSON       Format = "json"
)

// App -.
type Endpoints struct {
	CostsUrl string `env-required:"true" yaml:"costs" env:"CONFLUENT_COSTS_URL"`
}

// Credentials -.
type Credentials struct {
	Key    string `env-required:"true" yaml:"key" env:"CONFLUENT_CLOUD_API_KEY"`
	Secret string `env-required:"true" yaml:"secret" env:"CONFLUENT_CLOUD_API_SECRET"`
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cfg, nil
}
