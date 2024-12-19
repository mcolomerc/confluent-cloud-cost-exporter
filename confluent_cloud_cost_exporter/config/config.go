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
	Cron        Cron `yaml:"cron,omitempty" omitempty:"true"`
	Web         Web  `yaml:"web,omitempty" omitempty:"true"`
}

type Cron struct {
	Expression string `yaml:"expression,omitempty" env:"CRON_EXPRESSION"`
	Target     `yaml:"target,omitempty"`
}

type Target struct {
	KafkaConfig `yaml:"kafka,omitempty"`
}

type Web struct {
	Cache      `yaml:"cache,omitempty"`
	PromConfig `yaml:"prometheus,omitempty"`
}

type Cache struct {
	Expiration time.Duration `yaml:"expiration,omitempty" env:"CACHE_EXPIRATION" env-default:"60m"`
}

// Create a enum for the different types of exporters
type Exporter string

const (
	PROMETHEUS Exporter = "prometheus"
	JSON       Exporter = "json"
	KAFKA      Exporter = "kafka"
)

// App -.
type Endpoints struct {
	CostsUrl string `yaml:"costs" env:"CONFLUENT_COSTS_URL" env-default:"https://api.confluent.cloud/billing/v1/costs"`
}

// Credentials -.
type Credentials struct {
	Key    string `yaml:"key" env:"CONFLUENT_CLOUD_API_KEY"`
	Secret string `yaml:"secret" env:"CONFLUENT_CLOUD_API_SECRET"`
}

// Kafka Config -.
type KafkaConfig struct {
	Bootstrap      string `yaml:"bootstrap,omitempty" env:"BOOTSTRAP"`
	Credentials    `yaml:"credentials,omitempty"`
	SchemaRegistry `yaml:"schemaRegistry,omitempty"`
	Topic          string `yaml:"topic,omitempty" env:"TOPIC"`
}

// Kafka Credentials -.
type KafkaCredentials struct {
	Key    string `yaml:"key" env:"KAFKA_API_KEY"`
	Secret string `yaml:"secret" env:"KAFKA_API_SECRET"`
}

// Schema Registry -.
type SchemaRegistry struct {
	Endpoint       string `yaml:"endpoint,omitempty" env:"SCHEMA_REGISTRY_URL"`
	SR_Credentials `yaml:"credentials,omitempty"`
}

// Credentials -.
type SR_Credentials struct {
	Key    string `yaml:"key" env:"SCHEMA_REGISTRY_API_KEY"`
	Secret string `yaml:"secret" env:"SCHEMA_REGISTRY_API_SECRET"`
}

// NewConfig returns app config.
func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	fmt.Println("Using Confluent Cloud Key: " + cfg.Credentials.Key)

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cfg, nil
}
