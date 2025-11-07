package main

import (
    "log"
    "net/http"

    "taskforge/internal/server"
)

func main() {
    srv := server.New()
    log.Println("TaskForge API starting on :8080")
    http.ListenAndServe(":8080", srv.Router())
}
