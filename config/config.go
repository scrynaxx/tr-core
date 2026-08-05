package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Rabbit   RabbitConfig
	Redis    RedisConfig
	SMTP     SMTPConfig
	JWT      JWTConfig
	Tracing  TracingConfig
}

type AppConfig struct {
	Environment Environment `env:"APP_ENVIRONMENT"`
}

type PostgresConfig struct {
	User     string `env:"PG_USER,required"`
	Password string `env:"PG_PASSWORD,required"`
	Database string `env:"PG_DATABASE,required"`
	Host     string `env:"PG_HOST,required"`
	Port     uint16 `env:"PG_PORT,required"`
	SSLMode  string `env:"PG_SSL_MODE" envDefault:"require"`
}

type RabbitConfig struct {
	User     string `env:"RABBIT_USER,required"`
	Password string `env:"RABBIT_PASSWORD,required"`
	Address  string `env:"RABBIT_ADDRESS,required"`
	Vhost    string `env:"RABBIT_VHOST,required"`
}

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD,required"`
	Address  string `env:"REDIS_ADDRESS,required"`
}

type SMTPConfig struct {
	Host     string `env:"SMTP_HOST,required"`
	Port     uint16 `env:"SMTP_PORT,required"`
	Username string `env:"SMTP_USERNAME,required"`
	Password string `env:"SMTP_PASSWORD,required"`
	Sender   string `env:"SMTP_SENDER,required"`
	Name     string `env:"SMTP_NAME,required"`
}

type JWTConfig struct {
	Issuer string `env:"JWT_ISSUER,required"`
	Secret string `env:"JWT_SECRET,required"`
}

type TracingConfig struct {
	Endpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,required"`
	Protocol string `env:"OTEL_EXPORTER_OTLP_PROTOCOL,required"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return &cfg, nil
}
