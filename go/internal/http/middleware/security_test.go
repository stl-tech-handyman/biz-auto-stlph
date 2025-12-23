package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeadersMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	headers := w.Header()

	// Check security headers
	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":        "1; mode=block",
		"Referrer-Policy":         "strict-origin-when-cross-origin",
		"Content-Security-Policy": "default-src 'self'",
		"Permissions-Policy":     "geolocation=(), microphone=(), camera=()",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := headers.Get(header)
		if actualValue != expectedValue {
			t.Errorf("expected %s=%s, got %s", header, expectedValue, actualValue)
		}
	}

	// Check that Server header is removed
	if headers.Get("Server") != "" {
		t.Error("expected Server header to be removed")
	}
}


