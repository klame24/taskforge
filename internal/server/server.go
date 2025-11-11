package server

import (
	"net/http"
	"taskforge/internal/config"
	"taskforge/internal/logger"
)

type Server struct {
	Config *config.Config
	Logger *logger.Logger
}

func New(config *config.Config) *Server {
	log := logger.New(config.LogLevel)

	return &Server{
		Config: config,
		Logger: log,
	}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Health check called",
			"method", r.Method,
			"path", r.URL.Path,
			"user_agent", r.UserAgent())
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Ready check called",
			"method", r.Method,
			"path", r.URL.Path)
		w.Write([]byte("ready"))
	})
	return mux
}
