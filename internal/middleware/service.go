package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

func generateRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "fallback-id"
	}

	return hex.EncodeToString(bytes)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}

	return ""
}
