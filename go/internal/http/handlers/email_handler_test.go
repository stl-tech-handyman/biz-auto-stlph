package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEmailHandler_HandleTest(t *testing.T) {
	// Set environment variable to skip actual email service call
	os.Setenv("EMAIL_SERVICE_URL", "http://localhost:9999")
	defer os.Unsetenv("EMAIL_SERVICE_URL")

	handler := NewEmailHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing to field",
			body: map[string]interface{}{
				"subject": "Test",
				"html":    "<p>Test</p>",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing subject and html",
			body: map[string]interface{}{
				"to": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "valid request with html",
			body: map[string]interface{}{
				"to":    "test@example.com",
				"html":  "<p>Test</p>",
				"subject": "Test Subject",
			},
			expectedStatus: http.StatusInternalServerError, // Will fail because email service is not running
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/email/test", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleTest(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestEmailHandler_HandleBookingDeposit(t *testing.T) {
	os.Setenv("EMAIL_SERVICE_URL", "http://localhost:9999")
	defer os.Unsetenv("EMAIL_SERVICE_URL")

	handler := NewEmailHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing name field",
			body: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "valid request",
			body: map[string]interface{}{
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedStatus: http.StatusInternalServerError, // Will fail because email service is not running
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/email/booking-deposit", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleBookingDeposit(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}


