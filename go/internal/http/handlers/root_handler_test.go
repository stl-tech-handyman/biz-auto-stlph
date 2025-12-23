package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler_HandleRoot(t *testing.T) {
	handler := NewRootHandler("dev")

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["service"] != "bizops360-api-go" {
					t.Errorf("expected service=bizops360-api-go, got %v", resp["service"])
				}
				if resp["status"] != "running" {
					t.Errorf("expected status=running, got %v", resp["status"])
				}
				if endpoints, ok := resp["endpoints"].(map[string]interface{}); ok {
					if _, ok := endpoints["health"]; !ok {
						t.Error("expected endpoints.health")
					}
					if _, ok := endpoints["stripe"]; !ok {
						t.Error("expected endpoints.stripe")
					}
				} else {
					t.Error("expected endpoints in response")
				}
			},
		},
		{
			name:           "POST request should fail",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			handler.HandleRoot(w, req)

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

