package main

import (
	"log"

	"github.com/getsentry/sentry-go"
	"markoni23/url-shortener/internal/app"
	"markoni23/url-shortener/internal/config"
)

func main() {
	cfg := config.LoadEnv()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Server.SentryDSN,
	}); err != nil {
		log.Printf("Sentry initialization failed: %v\n", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
