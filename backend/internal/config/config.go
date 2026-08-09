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
		AllowedOrigins    []string      `env:"HTTP_ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:5173"`
		ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
		ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"15s"`
		WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		HTTPShutdownTime  time.Duration `env:"HTTP_SHUTDOWN_TIME" envDefault:"30s"`
		AdminToken        string        `env:"ADMIN_TOKEN" envDefault:"admin_token"`
	}

	PGConfig struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
		DB       string `env:"POSTGRES_DB" envDefault:"onboarding"`
		User     string `env:"POSTGRES_USER" envDefault:"onboarding_user"`
		Password string `env:"POSTGRES_PASSWORD" envDefault:"12345"`
	}

	Config struct {
		HTTPConfig HTTPConfig
		PGConfig   PGConfig
	}
)

func (c *HTTPConfig) HTTPAddress() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *Config) ConstructPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.PGConfig.User,
		c.PGConfig.Password,
		net.JoinHostPort(c.PGConfig.Host, c.PGConfig.Port),
		c.PGConfig.DB,
	)
}

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
