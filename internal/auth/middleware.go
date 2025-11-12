package auth

import (
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const UserContextKey ContextKey = "user"

// AuthMiddleware возвращает middleware для проверки JWT
func (manager *JWTManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": true, "message": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		// Формат: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": true, "message": "Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := manager.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, `{"error": true, "message": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Добавляем claims в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext получает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}
