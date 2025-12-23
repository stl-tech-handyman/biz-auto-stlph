package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// RequestIDKey is the context key for request ID
type RequestIDKey struct{}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey{}, requestID)
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), randomString(9))
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

