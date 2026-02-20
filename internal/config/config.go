package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	Server   ServerConfig
	Database DBConfig
}

type ServerConfig struct {
	BasePath    string
	Port        string
	SentryDSN   string
	FrontendUrl string
}

type DBConfig struct {
	DatabaseUrl string
}

func LoadEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}

	return Config{
		Env: os.Getenv("ENV"),
		Server: ServerConfig{
			BasePath:    os.Getenv("BASE_PATH"),
			Port:        os.Getenv("PORT"),
			SentryDSN:   os.Getenv("SENTRY_DSN"),
			FrontendUrl: os.Getenv("FRONTEND_URL"),
		},
		Database: DBConfig{
			DatabaseUrl: os.Getenv("DATABASE_URL"),
		},
	}
}
