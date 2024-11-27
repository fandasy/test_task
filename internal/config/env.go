package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string `env:"ADDR"`
	Slog        string `env:"SLOG"`
	DBHost      string `env:"DB_HOST"`
	DBPort      int    `env:"DB_PORT"`
	DBName      string `env:"DB_NAME"`
	DBUser      string `env:"DB_USER"`
	DBPassword  string `env:"DB_PASSWORD"`
	YourAPIHost string `env:"YOUR_API_HOST"`
}

func LoadEnvConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
