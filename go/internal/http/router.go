package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bizops360/go-api/internal/app"
	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/http/handlers"
	"github.com/bizops360/go-api/internal/http/middleware"
	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/infra/stripe"
)


// Router sets up HTTP routes
type Router struct {
	formEventsHandler     *handlers.FormEventsHandler
	triggersHandler       *handlers.TriggersHandler
	stripeHandler         *handlers.StripeHandler
	stripeWebhookHandler  *handlers.StripeWebhookHandler
	estimateHandler       *handlers.EstimateHandler
	emailHandler          *handlers.EmailHandler
	calendarHandler       *handlers.CalendarHandler
	businessLeadHandler   *handlers.BusinessLeadHandler
	zapierHandler         *handlers.ZapierHandler
	healthHandler         *handlers.HealthHandler
	rootHandler           *handlers.RootHandler
	logger                *slog.Logger
	environment           string
}

// NewRouter creates a new router
func NewRouter(
	formEventsService *app.FormEventsService,
	triggersService *app.TriggersService,
	businessLoader *config.BusinessLoader,
	logger *slog.Logger,
	environment string,
) *Router {
	paymentsProvider := stripe.NewStripePayments()
	emailClient := email.NewEmailServiceClient()
	gmailSender, _ := email.NewGmailSender()
	
	return &Router{
		formEventsHandler:    handlers.NewFormEventsHandler(formEventsService),
		triggersHandler:      handlers.NewTriggersHandler(triggersService),
		stripeHandler:        handlers.NewStripeHandler(paymentsProvider),
		stripeWebhookHandler: handlers.NewStripeWebhookHandler(paymentsProvider, emailClient, gmailSender, logger),
		estimateHandler:      handlers.NewEstimateHandler(paymentsProvider),
		emailHandler:         handlers.NewEmailHandler(logger),
		calendarHandler:      handlers.NewCalendarHandler(logger),
		businessLeadHandler:  handlers.NewBusinessLeadHandler(businessLoader, logger),
		zapierHandler:        handlers.NewZapierHandler(logger),
		healthHandler:        handlers.NewHealthHandler(),
		rootHandler:          handlers.NewRootHandler(environment),
		logger:               logger,
		environment:          environment,
	}
}

// Handler returns the HTTP handler with all middleware applied
func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	// #region agent log
	if logFile, err := os.OpenFile(handlers.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "A", "location": "router.go:65", "message": "Router.Handler called - registering routes", "data": map[string]interface{}{"timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion

	// Swagger UI - MUST be before root handler to avoid conflicts
	openAPIPath := handlers.GetOpenAPIPath()
	
	// #region agent log
	if logFile, err := os.OpenFile(handlers.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "D", "location": "router.go:70", "message": "GetOpenAPIPath result", "data": map[string]interface{}{"openAPIPath": openAPIPath, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	swaggerHandler := handlers.NewSwaggerHandler(openAPIPath)
	mux.HandleFunc("/swagger", swaggerHandler.HandleSwaggerUI)
	mux.HandleFunc("/swagger-ui", swaggerHandler.HandleSwaggerUI)
	mux.HandleFunc("/swagger/", swaggerHandler.HandleSwaggerUI)
	mux.HandleFunc("/api/openapi.json", swaggerHandler.HandleOpenAPISpec)
	mux.HandleFunc("/api/openapi.yaml", swaggerHandler.HandleOpenAPISpec)
	
	// #region agent log
	if logFile, err := os.OpenFile(handlers.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "A", "location": "router.go:78", "message": "Swagger routes registered", "data": map[string]interface{}{"routes": []string{"/swagger", "/swagger-ui", "/swagger/", "/api/openapi.json", "/api/openapi.yaml"}, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	// Serve static Swagger UI file as fallback
	mux.HandleFunc("/swagger-simple", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./swagger-ui-simple.html")
	})
	
	// Serve direct Swagger HTML file
	mux.HandleFunc("/swagger.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./swagger.html")
	})
	
	// Serve test pages
	mux.HandleFunc("/test-final-invoice.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./test-final-invoice.html")
	})
	
	// Root endpoint - MUST be last to catch all other routes
	mux.HandleFunc("/", r.rootHandler.HandleRoot)

	// API v1 routes (new pipeline-based) - no auth required
	mux.Handle("/v1/form-events", r.formEventsHandler)
	mux.Handle("/v1/triggers", r.triggersHandler)

	// Calendar endpoint - no auth required
	mux.HandleFunc("/api/calendar/create", r.calendarHandler.HandleCreate)

	// Business-specific lead processing - no auth required
	mux.HandleFunc("/api/business/{businessId}/process-lead", r.businessLeadHandler.HandleProcessLead)

	// Zapier endpoint (legacy, matching Apps Script flow) - no auth required
	mux.HandleFunc("/api/zapier/process-lead", r.zapierHandler.HandleProcessLead)

	// Legacy API routes (matching JS API) - require API key
	mux.Handle("/api/stripe/deposit", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.stripeHandler.HandleDeposit)))
	mux.Handle("/api/stripe/deposit/calculate", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.stripeHandler.HandleDepositCalculate)))
	mux.Handle("/api/stripe/deposit/with-email", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.stripeHandler.HandleDepositWithEmail)))
	mux.Handle("/api/stripe/deposit/amount", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.stripeHandler.HandleGetDepositAmount)))
	mux.Handle("/api/stripe/final-invoice", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.stripeHandler.HandleFinalInvoice)))
	mux.Handle("/api/stripe/test", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.stripeHandler.HandleTest)))
	mux.Handle("/api/estimate", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.estimateHandler.HandleCalculate)))
	mux.Handle("/api/estimate/special-dates", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.estimateHandler.HandleSpecialDates)))
	mux.Handle("/api/email/test", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.emailHandler.HandleTest)))
	mux.Handle("/api/email/booking-deposit", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.emailHandler.HandleBookingDeposit)))
	mux.Handle("/api/email/final-invoice", middleware.APIKeyMiddleware(r.logger, http.HandlerFunc(r.emailHandler.HandleFinalInvoice)))
	
	// Stripe webhook - no auth required (Stripe signs the request)
	mux.HandleFunc("/api/stripe/webhook", r.stripeWebhookHandler.HandleWebhook)
	
	// Health endpoints - no auth required
	mux.HandleFunc("/api/health", r.healthHandler.HandleHealth)
	mux.HandleFunc("/api/health/ready", r.healthHandler.HandleReady)
	mux.HandleFunc("/api/health/live", r.healthHandler.HandleLive)

	// Wrap mux to log which handler is being called
	wrappedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// #region agent log
		if logFile, err := os.OpenFile(handlers.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "router.go:167", "message": "ServeMux routing request", "data": map[string]interface{}{"path": r.URL.Path, "method": r.Method, "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		mux.ServeHTTP(w, r)
	})
	
	// Apply middleware in order (outermost to innermost)
	handler := middleware.RequestIDMiddleware(wrappedMux)
	handler = middleware.SecurityHeadersMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.MaxRequestSizeMiddleware(middleware.DefaultMaxRequestSize)(handler)
	handler = middleware.RecoveryMiddleware(r.logger, handler)
	
	// Rate limiting (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute, r.logger)
	handler = middleware.RateLimitMiddleware(rateLimiter)(handler)
	
	// Logging middleware (should be last to capture all request details)
	handler = middleware.LoggingMiddleware(r.logger, handler)
	
	return handler
}
