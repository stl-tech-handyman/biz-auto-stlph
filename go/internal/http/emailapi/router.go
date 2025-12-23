package emailapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bizops360/go-api/internal/http/emailapi/handlers"
	"github.com/bizops360/go-api/internal/http/middleware"
	"github.com/bizops360/go-api/internal/infra/email"
)

// Router sets up HTTP routes for email API
type Router struct {
	emailHandler *handlers.EmailHandler
	logger       *slog.Logger
}

// NewEmailAPIRouter creates a new email API router
func NewEmailAPIRouter(gmailSender *email.GmailSender, logger *slog.Logger) *Router {
	return &Router{
		emailHandler: handlers.NewEmailHandler(gmailSender, logger),
		logger:       logger,
	}
}

// Handler returns the HTTP handler with all middleware applied
func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	// Email endpoints
	mux.HandleFunc("/api/email/send", r.emailHandler.HandleSend)
	mux.HandleFunc("/api/email/draft", r.emailHandler.HandleDraft)
	mux.HandleFunc("/api/email/test", r.emailHandler.HandleTest) // Test endpoint for compatibility
	mux.HandleFunc("/api/health", r.emailHandler.HandleHealth)
	mux.HandleFunc("/", r.emailHandler.HandleRoot)

	// Apply middleware in order (outermost to innermost)
	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.SecurityHeadersMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.MaxRequestSizeMiddleware(middleware.DefaultMaxRequestSize)(handler)
	handler = middleware.RecoveryMiddleware(r.logger, handler)
	
	// Rate limiting (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute, r.logger)
	handler = middleware.RateLimitMiddleware(rateLimiter)(handler)
	
	// API key authentication (required for all endpoints)
	handler = middleware.APIKeyMiddleware(r.logger, handler)
	
	// Logging middleware (should be last to capture all request details)
	handler = middleware.LoggingMiddleware(r.logger, handler)

	return handler
}

