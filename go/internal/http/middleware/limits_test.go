package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMaxRequestSizeMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test with 1KB limit
	handler := MaxRequestSizeMiddleware(1024)(nextHandler)

	tests := []struct {
		name           string
		contentLength  int64
		body           string
		expectedStatus int
	}{
		{
			name:           "small request",
			contentLength:  100,
			body:           strings.Repeat("a", 100),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "request at limit",
			contentLength:  1024,
			body:           strings.Repeat("a", 1024),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "request exceeds limit",
			contentLength:  2048,
			body:           strings.Repeat("a", 2048),
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			req.ContentLength = tt.contentLength
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}


