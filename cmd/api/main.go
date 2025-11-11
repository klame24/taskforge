package main

import (
	"log"
	"net/http"

	"taskforge/internal/config"
	"taskforge/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.New(cfg)
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
