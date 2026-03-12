package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type Config struct {
	Env      string
	Server   ServerConfig
	Database DBConfig
}

func (c *Config) IsDevelopmentEnv() bool {
	return c.Env == envDev
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

	env, exists := os.LookupEnv("ENV")
	if !exists {
		env = envDev
	}

	port, exists := os.LookupEnv("APP_PORT")
	if !exists {
		port = "8080"
	}

	basePath, exists := os.LookupEnv("BASE_PATH")
	if !exists {
		basePath = "http://localhost:8080"
	}

	frontendUrl, exists := os.LookupEnv("FRONTEND_URL")
	if !exists {
		frontendUrl = "http://localhost:5123"
	}

	databaseURL, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		databaseURL = "postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable"
	}

	return Config{
		Env: env,
		Server: ServerConfig{
			BasePath:    basePath,
			Port:        port,
			SentryDSN:   os.Getenv("SENTRY_DSN"),
			FrontendUrl: frontendUrl,
		},
		Database: DBConfig{
			DatabaseUrl: databaseURL,
		},
	}
}
