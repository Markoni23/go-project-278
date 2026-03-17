package main

import (
	"log"

	"markoni23/url-shortener/internal/app"
	"markoni23/url-shortener/internal/config"
	"markoni23/url-shortener/internal/db"
)

func main() {
	cfg := config.LoadEnv()

	database, err := db.InitDB(cfg.Database.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
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
