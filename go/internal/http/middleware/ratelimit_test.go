package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	logger := slog.Default()
	limiter := NewRateLimiter(3, 1*time.Minute, logger)

	ip := "127.0.0.1"

	// Should allow first 3 requests
	for i := 0; i < 3; i++ {
		if !limiter.Allow(ip) {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// 4th request should be blocked
	if limiter.Allow(ip) {
		t.Error("4th request should be blocked")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	logger := slog.Default()
	limiter := NewRateLimiter(2, 1*time.Minute, logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RateLimitMiddleware(limiter)(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("request %d should succeed, got status %d", i+1, w.Code)
		}
	}

	// 3rd request should be rate limited
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		remoteAddr     string
		expectedIP     string
	}{
		{
			name:       "X-Forwarded-For takes precedence",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.1"},
			remoteAddr: "127.0.0.1:12345",
			expectedIP: "192.168.1.1",
		},
		{
			name:       "X-Real-IP as fallback",
			headers:    map[string]string{"X-Real-IP": "10.0.0.1"},
			remoteAddr: "127.0.0.1:12345",
			expectedIP: "10.0.0.1",
		},
		{
			name:       "RemoteAddr as last resort",
			headers:    map[string]string{},
			remoteAddr: "127.0.0.1:12345",
			expectedIP: "127.0.0.1:12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			ip := getClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}


