package main

import (
	"log"

	"markoni23/url-shortener/internal/app"
	"markoni23/url-shortener/internal/config"
	"markoni23/url-shortener/internal/db"

	"github.com/getsentry/sentry-go"
)

func main() {
	cfg := config.LoadEnv()

	if cfg.Server.SentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: cfg.Server.SentryDSN,
		}); err != nil {
			log.Printf("Sentry initialization failed: %v\n", err)
		}
	} else {
		log.Println("No SentryDSN")
	}

	database, err := db.InitDB(cfg.Database.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err := app.Run(cfg, database); err != nil {
		log.Fatal(err)
	}
}
