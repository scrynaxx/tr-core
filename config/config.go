package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
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

type HTTPConfig struct {
	Address           string        `env:"HTTP_ADDRESS" envDefault:":8080"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"1m"`
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
	Issuer          string        `env:"JWT_ISSUER,required"`
	Secret          string        `env:"JWT_SECRET,required"`
	AccessLifetime  time.Duration `env:"JWT_ACCESS_LIFETIME" envDefault:"15m"`
	RefreshLifetime time.Duration `env:"JWT_REFRESH_LIFETIME" envDefault:"360h"`
}

type TracingConfig struct {
	Endpoint string `env:"OTEL_ENDPOINT,required"`
	Protocol string `env:"OTEL_PROTOCOL,required"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}
