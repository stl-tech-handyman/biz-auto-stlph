package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bizops360/go-api/internal/infra/stripe"
)

func TestStripeHandler_HandleDepositCalculate(t *testing.T) {
	// Set up test API key
	os.Setenv("SERVICE_API_KEY", "test-api-key")
	defer os.Unsetenv("SERVICE_API_KEY")

	handler := NewStripeHandler(stripe.NewStripePayments())

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "no params",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["ok"] != true {
					t.Errorf("expected ok=true, got %v", resp["ok"])
				}
			},
		},
		{
			name:           "with estimate",
			queryParams:    "?estimate=1000",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["ok"] != true {
					t.Errorf("expected ok=true, got %v", resp["ok"])
				}
				if resp["deposit"] == nil {
					t.Error("expected deposit in response")
				}
			},
		},
		{
			name:           "with deposit",
			queryParams:    "?deposit=500",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["ok"] != true {
					t.Errorf("expected ok=true, got %v", resp["ok"])
				}
				if resp["deposit"] == nil {
					t.Error("expected deposit in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/stripe/deposit/calculate"+tt.queryParams, nil)
			req.Header.Set("X-Api-Key", "test-api-key")
			w := httptest.NewRecorder()

			handler.HandleDepositCalculate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestStripeHandler_HandleDeposit(t *testing.T) {
	os.Setenv("SERVICE_API_KEY", "test-api-key")
	defer os.Unsetenv("SERVICE_API_KEY")

	handler := NewStripeHandler(stripe.NewStripePayments())

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request with estimate",
			body: map[string]interface{}{
				"email":          "test@example.com",
				"name":           "Test User",
				"estimatedTotal": 1000.0,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid request with deposit",
			body: map[string]interface{}{
				"email":        "test@example.com",
				"name":         "Test User",
				"depositValue": 500.0,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/stripe/deposit", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Api-Key", "test-api-key")
			w := httptest.NewRecorder()

			handler.HandleDeposit(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp["ok"] != true {
				t.Errorf("expected ok=true, got %v", resp["ok"])
			}
		})
	}
}

