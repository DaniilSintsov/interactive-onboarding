package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
)

type (
	HTTPConfig struct {
		Host              string        `env:"HTTP_HOST" envDefault:"localhost"`
		Port              string        `env:"HTTP_PORT" envDefault:"8080"`
		ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
		ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"15s"`
		WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		HTTPShutdownTime  time.Duration `env:"HTTP_SHUTDOWN_TIME" envDefault:"30s"`
	}
	Config struct {
		PG struct {
			Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
			Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
			DB       string `env:"POSTGRES_DB" envDefault:"onboarding"`
			User     string `env:"POSTGRES_USER" envDefault:"onboarding_user"`
			Password string `env:"POSTGRES_PASSWORD" envDefault:"12345"`
		}
	}
)

func (c HTTPConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *Config) ConstructPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.PG.User,
		c.PG.Password,
		net.JoinHostPort(c.PG.Host, c.PG.Port),
		c.PG.DB,
	)
}

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
