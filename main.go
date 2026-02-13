package main

import (
	"context"
	"fmt"
	"log"
	"markoni23/url-shortener/internal/app"
	"markoni23/url-shortener/internal/config"
	"markoni23/url-shortener/internal/db"

	"github.com/getsentry/sentry-go"
)

func main() {
	cfg := config.LoadEnv()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Server.SentryDSN,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	ctx := context.Background()
	db, err := db.InitDB(cfg.Database.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err := app.Run(ctx, cfg, db); err != nil {
		log.Fatal(err)
	}
}
