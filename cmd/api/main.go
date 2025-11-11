package main

import (
	"log"
	"net/http"

	"taskforge/internal/config"
	"taskforge/internal/db"
	"taskforge/internal/server"
)

func main() {
	cfg := config.Load()

    db, err := db.ConnectDB(cfg.DB_DSN)
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
    defer db.Close()

	srv := server.New(cfg, db)
	addr := ":" + cfg.HTTPPort
	srv.Logger.Info(
		"Starting TaskForge API",
		"port", cfg.HTTPPort,
		"log_level", cfg.LogLevel,
		"environment", "development")

	if err := http.ListenAndServe(addr, srv.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
