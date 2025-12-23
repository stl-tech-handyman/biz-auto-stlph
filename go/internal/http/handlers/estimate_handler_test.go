package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bizops360/go-api/internal/infra/stripe"
)

func TestEstimateHandler_HandleCalculate(t *testing.T) {
	handler := NewEstimateHandler(stripe.NewStripePayments())

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "valid request",
			body: map[string]interface{}{
				"eventDate":     "2025-06-15",
				"durationHours": 4.0,
				"numHelpers":    2,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["ok"] != true {
					t.Error("expected ok=true")
				}
				if data, ok := resp["data"].(map[string]interface{}); ok {
					if data["totalCost"] == nil {
						t.Error("expected totalCost in response")
					}
					if deposit, ok := data["deposit"].(map[string]interface{}); ok {
						if deposit["recommended"] == nil {
							t.Error("expected deposit.recommended")
						}
						if deposit["range"] == nil {
							t.Error("expected deposit.range")
						}
						if deposit["calculation"] == nil {
							t.Error("expected deposit.calculation")
						}
					} else {
						t.Error("expected deposit in response")
					}
				} else {
					t.Error("expected data in response")
				}
			},
		},
		{
			name: "missing eventDate",
			body: map[string]interface{}{
				"durationHours": 4.0,
				"numHelpers":    2,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid durationHours",
			body: map[string]interface{}{
				"eventDate":     "2025-06-15",
				"durationHours": -1.0,
				"numHelpers":    2,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid numHelpers",
			body: map[string]interface{}{
				"eventDate":     "2025-06-15",
				"durationHours": 4.0,
				"numHelpers":    0,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/estimate", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleCalculate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil && w.Code == http.StatusOK {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestEstimateHandler_HandleSpecialDates(t *testing.T) {
	handler := NewEstimateHandler(stripe.NewStripePayments())

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "default request",
			query:          "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["ok"] != true {
					t.Error("expected ok=true")
				}
				if resp["data"] == nil {
					t.Error("expected data in response")
				}
			},
		},
		{
			name:           "with years parameter",
			query:          "?years=2",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["yearsAhead"] != float64(2) {
					t.Errorf("expected yearsAhead=2, got %v", resp["yearsAhead"])
				}
			},
		},
		{
			name:           "invalid years parameter",
			query:          "?years=50",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				// Should default to 5 years
				if resp["yearsAhead"] != float64(5) {
					t.Errorf("expected yearsAhead=5 (default), got %v", resp["yearsAhead"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/estimate/special-dates"+tt.query, nil)
			w := httptest.NewRecorder()

			handler.HandleSpecialDates(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				tt.checkResponse(t, resp)
			}
		})
	}
}


