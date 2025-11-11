package middleware

import (
	"context"
	"net/http"
	"taskforge/internal/logger"
	"time"
)

type requestIDKey struct{}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := generateRequestID()

			wrappedWriter := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
			r = r.WithContext(ctx)
			w.Header().Set("X-Request-ID", requestID)

			log.Debug("Request started",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", requestID)

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)

			log.Info("HTTP request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrappedWriter.statusCode,
				"duration", duration.String(),
				"request_id", requestID,
				"user_agent", r.UserAgent(),
				"ip", r.RemoteAddr)
		})
	}
}
