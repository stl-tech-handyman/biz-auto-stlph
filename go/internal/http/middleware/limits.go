package middleware

import (
	"net/http"
)

const (
	// DefaultMaxRequestSize is 10MB (matching JS API)
	DefaultMaxRequestSize = 10 * 1024 * 1024
)

// MaxRequestSizeMiddleware limits the size of request bodies
func MaxRequestSizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				http.Error(w, "Request entity too large", http.StatusRequestEntityTooLarge)
				return
			}

			// Wrap body reader to enforce limit
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			next.ServeHTTP(w, r)
		})
	}
}


