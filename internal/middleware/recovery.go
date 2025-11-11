package middleware

import (
	"net/http"
	"taskforge/internal/logger"
)

func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := GetRequestID(r.Context())

					log.Error(
						"Recovered from panic",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
						"request_id", requestID,
						"ip", r.RemoteAddr,
					)

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
