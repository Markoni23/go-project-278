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

	return Config{
		Env: env,
		Server: ServerConfig{
			BasePath:    os.Getenv("BASE_PATH"),
			Port:        os.Getenv("APP_PORT"),
			SentryDSN:   os.Getenv("SENTRY_DSN"),
			FrontendUrl: os.Getenv("FRONTEND_URL"),
		},
		Database: DBConfig{
			DatabaseUrl: os.Getenv("DATABASE_URL"),
		},
	}
}
