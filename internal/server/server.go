package server

import (
	"database/sql"
	"net/http"
	"taskforge/internal/config"
	"taskforge/internal/logger"
	"taskforge/internal/middleware"
)

type Server struct {
	Config *config.Config
	Logger *logger.Logger
	DB     *sql.DB
}

func New(config *config.Config, db *sql.DB) *Server {
	log := logger.New(config.LogLevel)

	return &Server{
		Config: config,
		Logger: log,
		DB: db,
	}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Health check called",
			"method", r.Method,
			"path", r.URL.Path,
			"user_agent", r.UserAgent())
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Ready check called",
			"method", r.Method,
			"path", r.URL.Path)

		err := s.DB.Ping()
		if err != nil {
			s.Logger.Error("Database is inavailable", "error", err)
			http.Error(w, "Database is unavailable", http.StatusServiceUnavailable)
			return 
		}

		w.Write([]byte("ready\n"))
	})

	handler := middleware.RecoveryMiddleware(s.Logger)(mux)
	handler = middleware.LoggerMiddleware(s.Logger)(handler)

	return handler
}
