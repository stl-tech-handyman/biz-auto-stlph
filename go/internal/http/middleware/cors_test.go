package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORSMiddleware(nextHandler)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   func(*testing.T, http.Header)
	}{
		{
			name:           "regular request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, h http.Header) {
				if h.Get("Access-Control-Allow-Origin") != "*" {
					t.Error("expected Access-Control-Allow-Origin header")
				}
				if h.Get("Access-Control-Allow-Methods") == "" {
					t.Error("expected Access-Control-Allow-Methods header")
				}
			},
		},
		{
			name:           "OPTIONS preflight",
			method:         "OPTIONS",
			expectedStatus: http.StatusNoContent,
			checkHeaders: func(t *testing.T, h http.Header) {
				if h.Get("Access-Control-Allow-Origin") != "*" {
					t.Error("expected Access-Control-Allow-Origin header")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkHeaders != nil {
				tt.checkHeaders(t, w.Header())
			}
		})
	}
}


