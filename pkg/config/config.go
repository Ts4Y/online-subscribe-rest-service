package config

import (
	"errors"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP     HTTP
	Postgres Postgres
	Logger   Logger
}

type HTTP struct {
	Port         int           `env:"HTTP_PORT"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT"`
}

type Postgres struct {
	DSN string `env:"POSTGRES_DSN"`
}

type Logger struct {
	Mode string `env:"LOGGER_MODE"`
}

func New(envPath string) (Config, error) {
	var c Config

	err := godotenv.Load(envPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	err = env.Parse(&c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
