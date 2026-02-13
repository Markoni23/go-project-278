package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	envDevelop = "development"
	envProd    = "production"
)

type Config struct {
	Env      string
	Server   ServerConfig
	Database DBConfig
}

type ServerConfig struct {
	BasePath  string
	Port      string
	SentryDSN string
}

type DBConfig struct {
	DatabaseUrl string
}

func LoadEnv() Config {
	godotenv.Load()
	return Config{
		Env: os.Getenv("ENV"),
		Server: ServerConfig{
			BasePath:  os.Getenv("BASE_PATH"),
			Port:      os.Getenv("PORT"),
			SentryDSN: os.Getenv("SENTRY_DSN"),
		},
		Database: DBConfig{
			DatabaseUrl: os.Getenv("DATABASE_URL"),
		},
	}
}
