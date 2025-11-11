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
    log.Printf("TaskForge API starting on %s", addr)
    if err := http.ListenAndServe(addr, srv.Router()); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
