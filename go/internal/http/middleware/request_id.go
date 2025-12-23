package middleware

import (
	"net/http"

	"github.com/bizops360/go-api/internal/util"
)

// RequestIDMiddleware adds a request ID to the context and response headers
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = util.GenerateRequestID()
		}

		ctx := util.WithRequestID(r.Context(), requestID)
		w.Header().Set("X-Request-Id", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

