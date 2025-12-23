package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/bizops360/go-api/internal/util"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := util.GetRequestID(r.Context())
				stack := debug.Stack()

				logger.Error("panic_recovered",
					"requestId", requestID,
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
					"stack", string(stack),
				)

				util.WriteError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}


