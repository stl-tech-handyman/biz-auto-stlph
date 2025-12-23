package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

// LoggingMiddleware logs HTTP requests with structured logging
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := util.GetRequestID(r.Context())

		// Get client IP (check headers for proxy/load balancer)
		clientIP := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			clientIP = xff
		} else if xri := r.Header.Get("X-Real-IP"); xri != "" {
			clientIP = xri
		}

		// Log request start with full context
		logger.Info("request_start",
			"requestId", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"userAgent", r.UserAgent(),
			"ip", clientIP,
			"contentLength", r.ContentLength,
			"contentType", r.Header.Get("Content-Type"),
		)

		// Wrap response writer to capture status code and response size
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK, bytesWritten: 0}

		next.ServeHTTP(wrapped, r)

		// Log request end with performance metrics
		duration := time.Since(start)
		logger.Info("request_end",
			"requestId", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", duration.String(),
			"durationMs", duration.Milliseconds(),
			"bytesWritten", wrapped.bytesWritten,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

