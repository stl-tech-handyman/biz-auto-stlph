package middleware

import (
	"net/http"
)

// SecurityHeadersMiddleware adds security headers (helmet equivalent)
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Remove server header (set by Go by default)
		w.Header().Del("Server")

		next.ServeHTTP(w, r)
	})
}


