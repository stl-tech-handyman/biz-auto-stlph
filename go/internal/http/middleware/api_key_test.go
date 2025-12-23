package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAPIKeyMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		headerValue    string
		expectedStatus int
		setEnv         bool
	}{
		{
			name:           "valid API key",
			apiKey:         "test-key-123",
			headerValue:    "test-key-123",
			expectedStatus: http.StatusOK,
			setEnv:         true,
		},
		{
			name:           "invalid API key",
			apiKey:         "test-key-123",
			headerValue:    "wrong-key",
			expectedStatus: http.StatusUnauthorized,
			setEnv:         true,
		},
		{
			name:           "missing API key header",
			apiKey:         "test-key-123",
			headerValue:    "",
			expectedStatus: http.StatusUnauthorized,
			setEnv:         true,
		},
		{
			name:           "API key not configured",
			apiKey:         "",
			headerValue:    "test-key-123",
			expectedStatus: http.StatusInternalServerError,
			setEnv:         false,
		},
		{
			name:           "API key with whitespace",
			apiKey:         "test-key-123",
			headerValue:    "  test-key-123  ",
			expectedStatus: http.StatusOK,
			setEnv:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.setEnv {
				os.Setenv("SERVICE_API_KEY", tt.apiKey)
			} else {
				os.Unsetenv("SERVICE_API_KEY")
			}
			defer os.Unsetenv("SERVICE_API_KEY")

			// Create a simple handler that returns 200
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Create logger
			logger := slog.Default()

			// Create middleware
			handler := APIKeyMiddleware(logger, nextHandler)

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-Api-Key", tt.headerValue)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(w, req)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

