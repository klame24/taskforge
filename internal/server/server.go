package server

import (
	"net/http"
	"taskforge/internal/config"
)

type Server struct{
    Config *config.Config
}

func New(config *config.Config) *Server {
    return &Server{
        Config: config,
    }
}

func (s *Server) Router() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })
    mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ready"))
    })
    return mux
}
