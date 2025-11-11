package main

import (
	"log"
	"net/http"
	"time"

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

	// Простая версия с таймаутами
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
