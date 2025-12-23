package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bizops360/go-api/internal/app"
	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/infra/db"
	logger "github.com/bizops360/go-api/internal/infra/log"
)

func TestRouter_AllEndpoints(t *testing.T) {
	// Set up test environment
	os.Setenv("SERVICE_API_KEY", "test-api-key")
	os.Setenv("STRIPE_SECRET_KEY_PROD", "sk_test_fake_key_for_testing")
	defer func() {
		os.Unsetenv("SERVICE_API_KEY")
		os.Unsetenv("STRIPE_SECRET_KEY_PROD")
	}()

	// Initialize minimal dependencies
	cfg := config.LoadConfig()
	log := logger.NewLogger("debug")
	businessLoader := config.NewBusinessLoader(cfg)
	jobsRepo := db.NewMemoryJobsRepo()

	actions := map[string]domain.Action{
		"normalize_input": &app.NormalizeInputAction{},
	}
	pipelineRunner := domain.NewPipelineRunner(actions)

	formEventsService := app.NewFormEventsService(businessLoader, pipelineRunner, jobsRepo)
	triggersService := app.NewTriggersService(businessLoader, pipelineRunner, jobsRepo)

	router := NewRouter(formEventsService, triggersService, log, "dev")
	handler := router.Handler()

	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		body           string
		expectedStatus int
		skipAuth       bool
	}{
		// Root endpoint
		{
			name:           "GET /",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			skipAuth:       true,
		},
		// Health endpoints
		{
			name:           "GET /api/health",
			method:         "GET",
			path:           "/api/health",
			expectedStatus: http.StatusOK,
			skipAuth:       true,
		},
		{
			name:           "GET /api/health/ready",
			method:         "GET",
			path:           "/api/health/ready",
			expectedStatus: http.StatusOK,
			skipAuth:       true,
		},
		{
			name:           "GET /api/health/live",
			method:         "GET",
			path:           "/api/health/live",
			expectedStatus: http.StatusOK,
			skipAuth:       true,
		},
		// Stripe endpoints (require auth)
		{
			name:           "GET /api/stripe/deposit/calculate",
			method:         "GET",
			path:           "/api/stripe/deposit/calculate?estimate=1000",
			headers:        map[string]string{"X-Api-Key": "test-api-key"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/stripe/deposit",
			method:         "POST",
			path:           "/api/stripe/deposit",
			headers:        map[string]string{"X-Api-Key": "test-api-key", "Content-Type": "application/json"},
			body:           `{"email":"test@example.com","name":"Test","estimatedTotal":1000}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/stripe/test",
			method:         "POST",
			path:           "/api/stripe/test",
			headers:        map[string]string{"X-Api-Key": "test-api-key", "Content-Type": "application/json"},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest, // Will fail because Stripe API key is fake
		},
		// Estimate endpoints (require auth)
		{
			name:           "POST /api/estimate",
			method:         "POST",
			path:           "/api/estimate",
			headers:        map[string]string{"X-Api-Key": "test-api-key", "Content-Type": "application/json"},
			body:           `{"eventDate":"2025-06-15","durationHours":4,"numHelpers":2}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/estimate/special-dates",
			method:         "GET",
			path:           "/api/estimate/special-dates?years=2",
			headers:        map[string]string{"X-Api-Key": "test-api-key"},
			expectedStatus: http.StatusOK,
		},
		// Email endpoints (require auth)
		{
			name:           "POST /api/email/test",
			method:         "POST",
			path:           "/api/email/test",
			headers:        map[string]string{"X-Api-Key": "test-api-key", "Content-Type": "application/json"},
			body:           `{"to":"test@example.com","subject":"Test","html":"<p>Test</p>"}`,
			expectedStatus: http.StatusInternalServerError, // Will fail because email service is not running
		},
		{
			name:           "POST /api/email/booking-deposit",
			method:         "POST",
			path:           "/api/email/booking-deposit",
			headers:        map[string]string{"X-Api-Key": "test-api-key", "Content-Type": "application/json"},
			body:           `{"name":"Test User","email":"test@example.com"}`,
			expectedStatus: http.StatusInternalServerError, // Will fail because email service is not running
		},
		// Auth tests
		{
			name:           "POST /api/estimate without API key",
			method:         "POST",
			path:           "/api/estimate",
			headers:        map[string]string{"Content-Type": "application/json"},
			body:           `{"eventDate":"2025-06-15","durationHours":4,"numHelpers":2}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "POST /api/estimate with wrong API key",
			method:         "POST",
			path:           "/api/estimate",
			headers:        map[string]string{"X-Api-Key": "wrong-key", "Content-Type": "application/json"},
			body:           `{"eventDate":"2025-06-15","durationHours":4,"numHelpers":2}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader *strings.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			} else {
				bodyReader = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

