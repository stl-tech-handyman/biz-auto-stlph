package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/infra/firestore"
	"github.com/bizops360/go-api/internal/infra/geo"
	"github.com/bizops360/go-api/internal/infra/storage"
	"github.com/bizops360/go-api/internal/infra/stripe"
	"github.com/bizops360/go-api/internal/infra/weather"
	"github.com/bizops360/go-api/internal/ports"
	emailService "github.com/bizops360/go-api/internal/services/email"
	"github.com/bizops360/go-api/internal/services/pdf"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
)

// getStringFromMap safely gets a string value from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// EmailHandler handles email-related endpoints
type EmailHandler struct {
	emailClient          *email.EmailServiceClient
	gmailSender          *email.GmailSender
	geocodingService     *geo.GeocodingService
	distanceMatrixService *geo.DistanceMatrixService
	weatherService       *weather.WeatherService
	businessLoader       *config.BusinessLoader
	pdfService           *pdf.Service
	logger               *slog.Logger
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(logger *slog.Logger) *EmailHandler {
	return NewEmailHandlerWithBusinessLoader(logger, nil)
}

// NewEmailHandlerWithBusinessLoader creates a new email handler with business loader
func NewEmailHandlerWithBusinessLoader(logger *slog.Logger, businessLoader *config.BusinessLoader) *EmailHandler {
	handler := &EmailHandler{
		logger:         logger,
		businessLoader: businessLoader,
	}

	// Try to use email service client first (if EMAIL_SERVICE_URL is set)
	handler.emailClient = email.NewEmailServiceClient()
	if handler.emailClient != nil {
		logger.Info("Using email service API for email sending")
	}

	// Initialize geocoding service (optional - for weather)
	if geoService, err := geo.NewGeocodingService(); err == nil {
		handler.geocodingService = geoService
		logger.Info("Geocoding service initialized for weather")
	} else {
		logger.Warn("Geocoding service not available", "error", err)
	}

	// Initialize distance matrix service (optional - for driving distance)
	if distService, err := geo.NewDistanceMatrixService(); err == nil {
		handler.distanceMatrixService = distService
		logger.Info("Distance Matrix service initialized for driving distance")
	} else {
		logger.Warn("Distance Matrix service not available, will use Haversine formula as fallback", "error", err)
	}

	// Initialize weather service (optional)
	if weatherService, err := weather.NewWeatherService(); err == nil {
		handler.weatherService = weatherService
		logger.Info("Weather service initialized")
	} else {
		logger.Warn("Weather service not available", "error", err)
	}

	// Fall back to Gmail sender (if credentials are available)
	// #region agent log
	if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1,H2,H3,H4,H5",
			"location":     "email_handler.go:NewEmailHandler",
			"message":      "Attempting to create Gmail sender",
			"data":         map[string]interface{}{},
			"timestamp":    time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	if gmailSender, err := email.NewGmailSender(); err == nil {
		handler.gmailSender = gmailSender
		logger.Info("Using Gmail API for email sending")
		// #region agent log
		if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "H1,H2,H3,H4,H5",
				"location":     "email_handler.go:NewEmailHandler",
				"message":      "Gmail sender created successfully",
				"data":         map[string]interface{}{},
				"timestamp":    time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
	} else {
		logger.Warn("Gmail API not available", "error", err)
		logger.Warn("Email functionality requires EMAIL_SERVICE_URL or GMAIL_CREDENTIALS_JSON to be configured")
		// #region agent log
		if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "H1,H2,H3,H4,H5",
				"location":     "email_handler.go:NewEmailHandler",
				"message":      "Gmail sender creation failed",
				"data": map[string]interface{}{
					"error": err.Error(),
				},
				"timestamp": time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
	}

	// Initialize PDF service (optional - for async PDF generation and storage)
	ctx := context.Background()
	var firestoreClient *firestore.Client
	var storageClient *storage.Client

	if projectID := os.Getenv("GCP_PROJECT_ID"); projectID != "" {
		if client, err := firestore.NewClient(ctx, projectID); err == nil {
			firestoreClient = client
			logger.Info("Firestore client initialized for PDF service")
		} else {
			logger.Warn("Firestore client not available for PDF service", "error", err)
		}
	}

	if bucketName := os.Getenv("GCS_BUCKET_NAME"); bucketName != "" {
		if client, err := storage.NewClient(ctx, bucketName); err == nil {
			storageClient = client
			logger.Info("Cloud Storage client initialized for PDF service")
		} else {
			logger.Warn("Cloud Storage client not available for PDF service", "error", err)
		}
	}

	if firestoreClient != nil && storageClient != nil {
		handler.pdfService = pdf.NewService(firestoreClient, storageClient, logger)
		logger.Info("PDF service initialized")
	} else {
		logger.Warn("PDF service not available - GCP_PROJECT_ID and GCS_BUCKET_NAME must be set")
	}

	return handler
}

// IsEmailServiceAvailable checks if email service is configured and available
func (h *EmailHandler) IsEmailServiceAvailable() bool {
	return h.gmailSender != nil || h.emailClient != nil
}

// HandleTest handles POST /api/email/test
func (h *EmailHandler) HandleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		To       string `json:"to"`
		Subject  string `json:"subject"`
		HTML     string `json:"html"`
		Text     string `json:"text"`
		From     string `json:"from"`
		FromName string `json:"fromName"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.To == "" {
		util.WriteError(w, http.StatusBadRequest, "to (recipient email) is required")
		return
	}

	if body.Subject == "" && body.HTML == "" {
		util.WriteError(w, http.StatusBadRequest, "either subject+html or html is required")
		return
	}

	req := &ports.SendEmailRequest{
		To:       body.To,
		Subject:  body.Subject,
		HTMLBody: body.HTML,
		TextBody: body.Text,
		From:     body.From,
		FromName: body.FromName,
	}

	var result *ports.SendEmailResult
	var err error

	// #region agent log
	if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1,H2,H3,H4,H5",
			"location":     "email_handler.go:HandleTest",
			"message":      "Checking email service availability",
			"data": map[string]interface{}{
				"gmailSenderIsNil": h.gmailSender == nil,
				"emailClientIsNil": h.emailClient == nil,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Use Gmail sender if available, otherwise use HTTP client
	if h.gmailSender != nil {
		result, err = h.gmailSender.SendEmail(r.Context(), req)
	} else if h.emailClient != nil {
		result, err = h.emailClient.SendEmail(r.Context(), req)
	} else {
		util.WriteError(w, http.StatusServiceUnavailable, "email service is not configured. Please set GMAIL_CREDENTIALS_JSON or EMAIL_SERVICE_URL")
		return
	}

	if err != nil {
		h.logger.Error("failed to send email", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to send email: "+err.Error())
		return
	}

	if !result.Success {
		errorMsg := "unknown error"
		if result.Error != nil {
			errorMsg = *result.Error
		}
		h.logger.Error("email sending failed", "error", errorMsg)
		util.WriteError(w, http.StatusInternalServerError, "email sending failed: "+errorMsg)
		return
	}

	h.logger.Info("email sent successfully", "messageId", result.MessageID, "to", req.To)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Email sent successfully",
		"result": map[string]interface{}{
			"messageId": result.MessageID,
			"success":   result.Success,
		},
	})
}

// HandleBookingDeposit handles POST /api/email/booking-deposit
func (h *EmailHandler) HandleBookingDeposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body map[string]interface{}
	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	name, _ := body["name"].(string)
	if name == "" {
		util.WriteError(w, http.StatusBadRequest, "Missing required field: name")
		return
	}

	// Convert body to SendEmailRequest for booking deposit
	emailReq := &ports.SendEmailRequest{
		To:       getStringFromMap(body, "email"),
		Subject:  "Booking Deposit Confirmation",
		HTMLBody: fmt.Sprintf("<p>Hello %s,</p><p>Your booking deposit has been processed.</p>", getStringFromMap(body, "name")),
		FromName: "BizOps360",
	}

	var emailResult *ports.SendEmailResult
	var err error

	if h.gmailSender != nil {
		emailResult, err = h.gmailSender.SendEmail(r.Context(), emailReq)
	} else if h.emailClient != nil {
		emailResult, err = h.emailClient.SendEmail(r.Context(), emailReq)
	} else {
		util.WriteError(w, http.StatusServiceUnavailable, "email service is not configured")
		return
	}

	if err != nil {
		h.logger.Error("failed to send booking deposit email", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to send booking deposit email: "+err.Error())
		return
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("booking deposit email sending failed", "error", errorMsg)
		util.WriteError(w, http.StatusInternalServerError, "email sending failed: "+errorMsg)
		return
	}

	result := map[string]interface{}{
		"messageId": emailResult.MessageID,
		"success":   emailResult.Success,
	}

	h.logger.Info("booking deposit email sent successfully", "messageId", emailResult.MessageID, "to", emailReq.To)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Booking deposit email sent successfully",
		"result":  result,
	})
}

// HandleFinalInvoice handles POST /api/email/final-invoice
func (h *EmailHandler) HandleFinalInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Name             string  `json:"name"`
		Email            string  `json:"email"`
		TotalAmount      float64 `json:"totalAmount"`
		DepositPaid      float64 `json:"depositPaid"`
		RemainingBalance float64 `json:"remainingBalance"`
		InvoiceURL       string  `json:"invoiceUrl"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.Name == "" {
		util.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	if body.Email == "" {
		util.WriteError(w, http.StatusBadRequest, "email is required")
		return
	}

	if body.InvoiceURL == "" {
		util.WriteError(w, http.StatusBadRequest, "invoiceUrl is required")
		return
	}

	// Generate email HTML from template
	templateService := emailService.NewTemplateService()
	// For the standalone email endpoint, use defaults for missing fields
	htmlBody, textBody, err := templateService.GenerateFinalInvoiceEmail(
		body.Name, "Event", "", nil, body.TotalAmount, body.DepositPaid, body.RemainingBalance, body.InvoiceURL, true, "")
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to generate email template: "+err.Error())
		return
	}

	emailReq := &ports.SendEmailRequest{
		To:       body.Email,
		Subject:  "Final Invoice - STL Party Helpers",
		HTMLBody: htmlBody,
		TextBody: textBody,
		FromName: "STL Party Helpers",
	}

	var emailResult *ports.SendEmailResult
	err = nil

	if h.gmailSender != nil {
		emailResult, err = h.gmailSender.SendEmail(r.Context(), emailReq)
	} else if h.emailClient != nil {
		emailResult, err = h.emailClient.SendEmail(r.Context(), emailReq)
	} else {
		util.WriteError(w, http.StatusServiceUnavailable, "email service is not configured")
		return
	}

	if err != nil {
		h.logger.Error("failed to send final invoice email", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to send final invoice email: "+err.Error())
		return
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("final invoice email sending failed", "error", errorMsg)
		util.WriteError(w, http.StatusInternalServerError, "email sending failed: "+errorMsg)
		return
	}

	result := map[string]interface{}{
		"messageId": emailResult.MessageID,
		"success":   emailResult.Success,
	}

	h.logger.Info("final invoice email sent successfully", "messageId", emailResult.MessageID, "to", emailReq.To)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Final invoice email sent successfully",
		"result":  result,
	})
}

// SendFinalInvoiceEmail is a helper method that can be called from other handlers
// Returns (success bool, errorMessage string)
func (h *EmailHandler) SendFinalInvoiceEmail(ctx context.Context, name, email, eventType, eventDate string, helpersCount *int, originalQuote, depositPaid, remainingBalance float64, invoiceURL string, showGratuity bool, saveAsDraft bool, templateName string) (bool, string) {
	if name == "" || email == "" || invoiceURL == "" {
		return false, "name, email, and invoiceUrl are required"
	}

	templateService := emailService.NewTemplateService()
	htmlBody, textBody, err := templateService.GenerateFinalInvoiceEmail(name, eventType, eventDate, helpersCount, originalQuote, depositPaid, remainingBalance, invoiceURL, showGratuity, templateName)
	if err != nil {
		return false, fmt.Sprintf("failed to generate email template: %v", err)
	}

	emailReq := &ports.SendEmailRequest{
		To:       email,
		Subject:  "Final Invoice - STL Party Helpers",
		HTMLBody: htmlBody,
		TextBody: textBody,
		FromName: "STL Party Helpers",
	}

	var emailResult *ports.SendEmailResult

	// #region agent log
	if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1",
			"location":     "email_handler.go:SendFinalInvoiceEmail",
			"message":      "Before calling email service",
			"data": map[string]interface{}{
				"saveAsDraft":    saveAsDraft,
				"hasGmailSender": h.gmailSender != nil,
				"hasEmailClient": h.emailClient != nil,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion

	if saveAsDraft {
		// Save as draft in Gmail
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmailDraft(ctx, emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmailDraft(ctx, emailReq)
		} else {
			return false, "email service is not configured"
		}
	} else {
		// Send email immediately
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmail(ctx, emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmail(ctx, emailReq)
		} else {
			return false, "email service is not configured"
		}
	}

	// #region agent log
	if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1",
			"location":     "email_handler.go:SendFinalInvoiceEmail",
			"message":      "After calling email service",
			"data": map[string]interface{}{
				"hasError": err != nil,
				"error": func() string {
					if err != nil {
						return err.Error()
					} else {
						return ""
					}
				}(),
				"hasResult": emailResult != nil,
				"resultSuccess": func() bool {
					if emailResult != nil {
						return emailResult.Success
					} else {
						return false
					}
				}(),
				"resultError": func() string {
					if emailResult != nil && emailResult.Error != nil {
						return *emailResult.Error
					} else {
						return ""
					}
				}(),
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion

	if err != nil {
		h.logger.Error("failed to send final invoice email", "error", err)
		return false, err.Error()
	}

	if emailResult == nil {
		h.logger.Error("final invoice email sending failed: emailResult is nil")
		return false, "email service returned nil result"
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("final invoice email sending failed", "error", errorMsg)
		return false, errorMsg
	}

	h.logger.Info("final invoice email sent successfully", "messageId", emailResult.MessageID, "to", email)
	return true, ""
}

// SendDepositEmail sends a deposit invoice email
// Returns (success bool, errorMessage string)
func (h *EmailHandler) SendDepositEmail(ctx context.Context, name, email string, depositAmount float64, invoiceURL string, saveAsDraft bool) (bool, string) {
	if name == "" || email == "" || invoiceURL == "" {
		return false, "name, email, and invoiceUrl are required"
	}

	templateService := emailService.NewTemplateService()
	htmlBody, textBody, err := templateService.GenerateDepositEmail(name, depositAmount, invoiceURL)
	if err != nil {
		return false, fmt.Sprintf("failed to generate email template: %v", err)
	}

	emailReq := &ports.SendEmailRequest{
		To:       email,
		Subject:  "Action needed to secure your reservation - STL Party Helpers",
		HTMLBody: htmlBody,
		TextBody: textBody,
		FromName: "STL Party Helpers",
	}

	var emailResult *ports.SendEmailResult

	if saveAsDraft {
		// Save as draft in Gmail
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmailDraft(ctx, emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmailDraft(ctx, emailReq)
		} else {
			return false, "email service is not configured"
		}
	} else {
		// Send email immediately
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmail(ctx, emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmail(ctx, emailReq)
		} else {
			return false, "email service is not configured"
		}
	}

	if err != nil {
		h.logger.Error("failed to send deposit email", "error", err)
		return false, err.Error()
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("deposit email sending failed", "error", errorMsg)
		return false, errorMsg
	}

	h.logger.Info("deposit email sent successfully", "messageId", emailResult.MessageID, "to", email)
	return true, ""
}

// HandleQuoteEmail handles POST /api/email/quote
func (h *EmailHandler) HandleQuoteEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		To            string  `json:"to"`
		ClientName    string  `json:"clientName"`
		EventDate     string  `json:"eventDate"` // Formatted date like "January 2, 2025"
		EventTime     string  `json:"eventTime"` // Time like "4:00 PM"
		EventLocation string  `json:"eventLocation"`
		Occasion      string  `json:"occasion"`
		GuestCount    int     `json:"guestCount"`
		Helpers       int     `json:"helpers"`
		Hours         float64 `json:"hours"`
		BaseRate      float64 `json:"baseRate"`
		HourlyRate    float64 `json:"hourlyRate"`
		TotalCost     float64 `json:"totalCost"`
		RateLabel     string  `json:"rateLabel"`
		DryRun        bool    `json:"dryRun"`
		SaveAsDraft   bool    `json:"saveAsDraft"`
		PayWithCheck  bool    `json:"payWithCheck"` // If true, attach PDF quote
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.To == "" {
		util.WriteError(w, http.StatusBadRequest, "to (recipient email) is required")
		return
	}

	// Parse event date to calculate correct rates for the year
	eventDate, parseErr := parseEventDateFromFormatted(body.EventDate)
	if parseErr != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid eventDate format: %v. Expected format: 'January 2, 2025'", parseErr))
		return
	}

	// Calculate estimate to get correct rates for the year
	estimate, calcErr := pricing.CalculateEstimate(eventDate, body.Hours, body.Helpers)
	if calcErr != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("failed to calculate estimate: %v", calcErr))
		return
	}

	// Use totalCost from body if provided, otherwise use calculated estimate
	totalCost := body.TotalCost
	if totalCost == 0 {
		totalCost = estimate.TotalCost
	}

	// Calculate deposit from total cost
	estimateCents := util.DollarsToCents(totalCost)
	depositCalc := stripe.CalculateDepositFromEstimate(estimateCents)
	depositAmount := util.CentsToDollars(depositCalc.Value)

	// Determine rate label
	rateLabel := body.RateLabel
	if rateLabel == "" {
		rateLabel = "Base Rate"
		if estimate.SpecialLabel != nil {
			rateLabel = *estimate.SpecialLabel
		}
	}

	// Calculate days until event and urgency level
	// Use calendar days (normalize to midnight for accurate day count)
	now := time.Now()
	location := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	eventDay := time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, location)
	daysUntilEvent := int(eventDay.Sub(today).Hours() / 24)
	if daysUntilEvent < 0 {
		daysUntilEvent = 0
	}

	// Determine urgency level
	urgencyLevel := util.CalculateUrgencyLevel(daysUntilEvent)

	// Calculate expiration date dynamically based on days until event
	expirationDate, expirationFormatted := util.CalculateExpirationDate(daysUntilEvent)

	// TODO: Create deposit invoice via Stripe to get actual payment link
	// For now, use placeholder - in production, create invoice and use HostedInvoiceURL
	depositLink := "#" // Placeholder - should be replaced with actual Stripe invoice URL

	// Generate confirmation number
	confirmationNumber := util.GenerateConfirmationNumber(body.To, body.Occasion, eventDate)

	// Fetch weather forecast if event is < 10 days away
	var weatherForecast *util.WeatherForecastData
	if daysUntilEvent < 10 && h.weatherService != nil && h.geocodingService != nil && body.EventLocation != "" {
		// Geocode address to get lat/lng
		geoResult, err := h.geocodingService.GetLatLng(r.Context(), body.EventLocation)
		if err == nil {
			// Fetch weather forecast
			forecast, err := h.weatherService.GetForecastForDate(r.Context(), geoResult.Lat, geoResult.Lng, eventDate)
			if err == nil && forecast != nil {
				// Determine if event is outdoor (simple heuristic - could be enhanced)
				isOutdoor := strings.Contains(strings.ToLower(body.Occasion), "outdoor") ||
					strings.Contains(strings.ToLower(body.EventLocation), "outdoor") ||
					strings.Contains(strings.ToLower(body.EventLocation), "garden") ||
					strings.Contains(strings.ToLower(body.EventLocation), "patio") ||
					strings.Contains(strings.ToLower(body.EventLocation), "park")

				recommendation := weather.GetWeatherRecommendation(forecast, isOutdoor)
				weatherForecast = &util.WeatherForecastData{
					Temperature:   forecast.Temperature,
					Condition:     forecast.Condition,
					Description:   forecast.Description,
					Recommendation: recommendation,
				}
			} else if err != nil {
				h.logger.Warn("Failed to fetch weather forecast", "error", err)
			}
		} else {
			h.logger.Warn("Failed to geocode address for weather", "error", err)
		}
	}

	// Calculate travel fee based on distance from office
	var travelFeeInfo *util.TravelFeeData
	var travelFeeAmount float64 = 0
	if body.EventLocation != "" && h.geocodingService != nil {
		// Get location info from business config (if available)
		var originLat, originLng, radiusMiles float64
		var originAddress string
		
		// Try to get business config if businessLoader is available
		// Default to "stlpartyhelpers" if no business ID in request
		businessID := "stlpartyhelpers" // Default business ID
		if h.businessLoader != nil {
			var businessConfig *domain.BusinessConfig
			businessConfig, err := h.businessLoader.LoadBusiness(r.Context(), businessID)
			if err == nil && businessConfig != nil {
				// Get location from business config
				lat, lng, radius, err := geo.LocationFromBusinessConfig(r.Context(), &businessConfig.Location, h.geocodingService)
				if err == nil {
					originLat = lat
					originLng = lng
					radiusMiles = radius
					originAddress = businessConfig.Location.OfficeAddress
					if originAddress == "" {
						originAddress = businessConfig.Location.DistanceOrigin
					}
				} else {
					h.logger.Warn("Failed to get location from business config, using defaults", "error", err)
				}
			}
		}
		
		// Fallback to hardcoded defaults if business config not available
		if originLat == 0 && originLng == 0 {
			originLat = geo.OfficeLat
			originLng = geo.OfficeLng
			radiusMiles = geo.ServiceRadiusMiles
			originAddress = geo.OfficeAddress
		}
		
		var distanceMiles float64
		var distanceSource string
		
		// Use Haversine formula for distance calculation (Distance Matrix API not implemented yet)
		if h.geocodingService != nil {
			// Geocode destination
			geoResult, err := h.geocodingService.GetLatLng(r.Context(), body.EventLocation)
			if err == nil {
				distanceMiles = geo.CalculateDistanceFromOrigin(originLat, originLng, geoResult.Lat, geoResult.Lng)
				distanceSource = "Haversine formula (straight-line distance)"
				h.logger.Info("Using Haversine formula for distance calculation",
					"distanceMiles", distanceMiles,
					"origin", originAddress,
					"originLat", originLat,
					"originLng", originLng,
				)
			} else {
				h.logger.Warn("Failed to geocode address for travel fee calculation", "error", err)
			}
		}
		
		// Calculate travel fee if we have a distance
		if distanceMiles > 0 {
			// Calculate travel fee - first check if within service area using config radius
			var travelFeeResult *pricing.TravelFeeResult
			if distanceMiles <= radiusMiles {
				// Within service area - no fee
				travelFeeResult = &pricing.TravelFeeResult{
					IsWithinServiceArea: true,
					DistanceMiles:       distanceMiles,
					TravelFee:           0,
					TravelFeePerHelper:  0,
					TotalTravelFee:      0,
					Message:             "within our service area - no travel fee",
				}
			} else {
				// Outside service area - calculate fee
				// Use standard calculation but adjust for custom radius
				milesOverServiceArea := distanceMiles - radiusMiles
				var feePerHelper float64
				if milesOverServiceArea <= 10.0 {
					feePerHelper = 40.0
				} else {
					milesBeyondFirst10 := milesOverServiceArea - 10.0
					increments := math.Ceil(milesBeyondFirst10 / 10.0)
					feePerHelper = 40.0 + (increments * 10.0)
				}
				feePerHelper = math.Round(feePerHelper)
				totalTravelFee := feePerHelper * float64(body.Helpers)
				
				var message string
				if body.Helpers == 1 {
					message = fmt.Sprintf("outside of our area — $%.0f travel fee", totalTravelFee)
				} else {
					message = fmt.Sprintf("outside of our area — $%.0f travel fee (for %d helpers)", totalTravelFee, body.Helpers)
				}
				
				travelFeeResult = &pricing.TravelFeeResult{
					IsWithinServiceArea: false,
					DistanceMiles:       distanceMiles,
					TravelFee:           totalTravelFee,
					TravelFeePerHelper:  feePerHelper,
					TotalTravelFee:      totalTravelFee,
					Message:             message,
				}
			}
			
			travelFeeAmount = travelFeeResult.TotalTravelFee
			
			// Add travel fee to total cost
			totalCost += travelFeeAmount
			
			// Create travel fee data for email
			travelFeeInfo = &util.TravelFeeData{
				IsWithinServiceArea: travelFeeResult.IsWithinServiceArea,
				DistanceMiles:       travelFeeResult.DistanceMiles,
				TravelFee:           travelFeeResult.TotalTravelFee,
				Message:             travelFeeResult.Message,
			}
			
			h.logger.Info("Travel fee calculated",
				"distanceMiles", distanceMiles,
				"distanceSource", distanceSource,
				"isWithinServiceArea", travelFeeResult.IsWithinServiceArea,
				"travelFee", travelFeeAmount,
				"helpers", body.Helpers,
			)
		}
	}

	// Generate PDF token and link (if PDF service is available)
	pdfDownloadLink := ""
	if h.pdfService != nil {
		pdfToken, err := util.GeneratePDFToken(confirmationNumber, body.To, expirationDate)
		if err == nil {
			// Build PDF download URL - use environment variable for base URL or default
			baseURL := os.Getenv("API_BASE_URL")
			if baseURL == "" {
				baseURL = "https://api.stlpartyhelpers.com"
			}
			pdfDownloadLink = fmt.Sprintf("%s/api/quote/pdf?token=%s", baseURL, pdfToken)
			h.logger.Info("PDF download link generated", "confirmationNumber", confirmationNumber)
		} else {
			h.logger.Warn("Failed to generate PDF token", "error", err)
		}
	}

	// Generate quote email HTML using rates from estimate
	emailData := util.QuoteEmailData{
		ClientName:         body.ClientName,
		EventDate:          body.EventDate,
		EventTime:          body.EventTime,
		EventLocation:      body.EventLocation,
		Occasion:           body.Occasion,
		GuestCount:         body.GuestCount,
		Helpers:            body.Helpers,
		Hours:              body.Hours,
		BaseRate:           estimate.BasePerHelper,
		HourlyRate:         estimate.ExtraPerHourPerHelper,
		TotalCost:          totalCost,
		DepositAmount:      depositAmount,
		RateLabel:          rateLabel,
		ExpirationDate:     expirationFormatted, // Email needs string
		DepositLink:        depositLink,
		ConfirmationNumber: confirmationNumber,
		IsHighDemand:       estimate.IsSpecialDate, // High demand = special date (holiday/surge)
		UrgencyLevel:       urgencyLevel,           // Urgency level based on days until event
		DaysUntilEvent:     daysUntilEvent,         // Number of days until event
		IsReturningClient:  false,                  // TODO: Check if client has booked before (query CRM/calendar)
		WeatherForecast:    weatherForecast,        // Weather forecast (only for events < 10 days)
		TravelFeeInfo:      travelFeeInfo,          // Travel fee information
		PDFDownloadLink:    pdfDownloadLink,        // PDF download link
	}

	htmlBody := util.GenerateQuoteEmailHTML(emailData)

	// Format date with day of week for subject
	eventDateWithDay := formatDateWithDayOfWeek(emailData.EventDate)

	// Shortened subject line
	subject := fmt.Sprintf("%s Quote - %s", emailData.Occasion, eventDateWithDay)

	if body.DryRun {
		subject = "Dry Run - " + subject
	}

	// Generate PDF quote for records (always attached)
	pdfData := util.QuotePDFData{
		ConfirmationNumber: confirmationNumber,
		Occasion:           body.Occasion,
		ClientName:         body.ClientName,
		ClientEmail:        body.To,
		EventDate:          body.EventDate,
		EventTime:          body.EventTime,
		HelpersCount:       body.Helpers,
		Hours:              body.Hours,
		TotalCost:          totalCost,
		DepositAmount:      depositAmount,
		ExpirationDate:     expirationDate,
		DepositLink:        depositLink,
		IssueDate:          time.Now(),
	}

	pdfBytes, err := util.GenerateQuotePDF(pdfData)
	if err != nil {
		h.logger.Error("failed to generate PDF", "error", err)
		// Continue without PDF attachment
	}

	emailReq := &ports.SendEmailRequest{
		To:       body.To,
		Subject:  subject,
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers Team",
	}

	// Attach PDF if generated successfully
	if pdfBytes != nil {
		emailReq.Attachments = []ports.Attachment{
			{
				Filename: fmt.Sprintf("%s-Quote-%s.pdf", body.Occasion, confirmationNumber),
				Content:  pdfBytes,
				MimeType: "application/pdf",
			},
		}
	}

	var emailResult *ports.SendEmailResult
	var emailErr error

	if body.SaveAsDraft {
		if h.gmailSender != nil {
			emailResult, emailErr = h.gmailSender.SendEmailDraft(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, emailErr = h.emailClient.SendEmailDraft(r.Context(), emailReq)
		} else {
			util.WriteError(w, http.StatusInternalServerError, "email service is not configured")
			return
		}
	} else {
		if h.gmailSender != nil {
			emailResult, emailErr = h.gmailSender.SendEmail(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, emailErr = h.emailClient.SendEmail(r.Context(), emailReq)
		} else {
			util.WriteError(w, http.StatusInternalServerError, "email service is not configured")
			return
		}
	}

	if emailErr != nil {
		h.logger.Error("failed to send quote email", "error", emailErr)
		util.WriteError(w, http.StatusInternalServerError, "failed to send quote email: "+emailErr.Error())
		return
	}

	if emailResult == nil {
		util.WriteError(w, http.StatusInternalServerError, "email service returned nil result")
		return
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("quote email sending failed", "error", errorMsg)
		util.WriteError(w, http.StatusInternalServerError, "quote email sending failed: "+errorMsg)
		return
	}

	sent := !body.SaveAsDraft
	draft := body.SaveAsDraft

	h.logger.Info("quote email sent successfully", "messageId", emailResult.MessageID, "to", body.To, "draft", draft)

	// Kick off async PDF generation and storage (if PDF service is available)
	if h.pdfService != nil && !body.SaveAsDraft {
		h.pdfService.GenerateAndStorePDFAsync(r.Context(), emailData, pdfData, expirationDate)
		h.logger.Info("PDF generation task queued", "confirmationNumber", confirmationNumber)
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Quote email sent successfully",
		"email": map[string]interface{}{
			"messageId": emailResult.MessageID,
			"sent":      sent,
			"draft":     draft,
			"error":     "",
		},
	})
}

// parseEventDateFromFormatted parses a formatted date string like "January 2, 2025" to time.Time
func parseEventDateFromFormatted(dateStr string) (time.Time, error) {
	// Try common date formats
	formats := []string{
		"January 2, 2006",
		"Jan 2, 2006",
		"2006-01-02",
		"01/02/2006",
		"1/2/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// formatDateWithDayOfWeek formats a date string to include day of week (e.g., "Mon, June 15, 2026")
func formatDateWithDayOfWeek(dateStr string) string {
	// Try to parse the date
	t, err := parseEventDateFromFormatted(dateStr)
	if err != nil {
		// If parsing fails, return original string
		return dateStr
	}

	// Format as "Mon, January 2, 2006" (short day name, full month)
	return t.Format("Mon, January 2, 2006")
}

// HandleQuoteEmailPreview handles POST /api/email/quote/preview - sends quote email with dummy data
func (h *EmailHandler) HandleQuoteEmailPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		To          string  `json:"to"`
		SaveAsDraft bool    `json:"saveAsDraft"`
		QuoteType   string  `json:"quoteType"` // e.g., "regular", "new_year", "thanksgiving", "surge"
		EventDate   string  `json:"eventDate"` // e.g., "2026-12-24" (YYYY-MM-DD format)
		EventTime   string  `json:"eventTime"` // e.g., "18:00" (HH:MM format)
		Helpers     int     `json:"helpers"`
		Hours       float64 `json:"hours"`
		Occasion    string  `json:"occasion"`
		GuestCount  int     `json:"guestCount"`
		ClientName  string  `json:"clientName"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.To == "" {
		util.WriteError(w, http.StatusBadRequest, "to (recipient email) is required")
		return
	}

	// Parse event date from YYYY-MM-DD format
	var parsedEventDate time.Time
	var parseErr error
	if body.EventDate != "" {
		parsedEventDate, parseErr = parseEventDateFromFormatted(body.EventDate)
		if parseErr != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid eventDate format: %v. Expected format: 'YYYY-MM-DD' or 'January 2, 2006'", parseErr))
			return
		}
	} else {
		// Default to current year if not provided
		parsedEventDate = time.Now()
	}

	// Use defaults if not provided
	helpers := body.Helpers
	if helpers == 0 {
		helpers = 2
	}
	hours := body.Hours
	if hours == 0 {
		hours = 4.0
	}
	occasion := body.Occasion
	if occasion == "" {
		occasion = "Birthday Party"
	}
	clientName := body.ClientName
	if clientName == "" {
		clientName = "John Doe"
	}
	guestCount := body.GuestCount
	if guestCount == 0 {
		guestCount = 50
	}

	// Calculate estimate using REAL pricing logic with the provided date
	estimate, calcErr := pricing.CalculateEstimate(parsedEventDate, hours, helpers)
	if calcErr != nil {
		util.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to calculate estimate: %v", calcErr))
		return
	}

	// Calculate deposit from total cost using proper Stripe deposit calculator
	estimateCents := util.DollarsToCents(estimate.TotalCost)
	depositCalc := stripe.CalculateDepositFromEstimate(estimateCents)
	depositAmount := util.CentsToDollars(depositCalc.Value)

	// Calculate days until event and urgency level
	// Use calendar days (normalize to midnight for accurate day count)
	now := time.Now()
	location := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	eventDay := time.Date(parsedEventDate.Year(), parsedEventDate.Month(), parsedEventDate.Day(), 0, 0, 0, 0, location)
	daysUntilEvent := int(eventDay.Sub(today).Hours() / 24)
	if daysUntilEvent < 0 {
		daysUntilEvent = 0
	}

	// Determine urgency level
	urgencyLevel := util.CalculateUrgencyLevel(daysUntilEvent)

	// Calculate expiration date dynamically based on days until event
	expirationDate, expirationFormatted := util.CalculateExpirationDate(daysUntilEvent)

	// Determine rate label
	rateLabel := "Base Rate"
	if estimate.SpecialLabel != nil {
		rateLabel = *estimate.SpecialLabel
	}

	// For preview, use a placeholder deposit link
	// In production, this would be generated by creating a deposit invoice
	depositLink := "https://invoice.stripe.com/i/acct_placeholder/preview_placeholder"

	// Generate test confirmation number for preview (increments)
	confirmationNumber := util.GenerateTestQuoteID()

	// Format event date for display
	// Format as "Fri, Jan 19, 2026" (day of week, short month, day, year)
	eventDateFormatted := parsedEventDate.Format("Mon, Jan 2, 2006")

	// Weather forecast not available in preview (would require API key and address)
	var weatherForecast *util.WeatherForecastData = nil

	// Travel fee not calculated in preview
	var travelFeeInfo *util.TravelFeeData = nil

	// Parse event time from HH:MM format to readable format (e.g., "18:00" -> "6:00 PM")
	eventTimeFormatted := "6:00 PM" // Default
	if body.EventTime != "" {
		// Parse time from "HH:MM" format
		timeParts := strings.Split(body.EventTime, ":")
		if len(timeParts) == 2 {
			hoursInt, err := strconv.Atoi(timeParts[0])
			if err == nil {
				minutes := timeParts[1]
				ampm := "AM"
				if hoursInt >= 12 {
					ampm = "PM"
					if hoursInt > 12 {
						hoursInt -= 12
					}
				}
				if hoursInt == 0 {
					hoursInt = 12
				}
				eventTimeFormatted = fmt.Sprintf("%d:%s %s", hoursInt, minutes, ampm)
			}
		}
	}

	emailData := util.QuoteEmailData{
		ClientName:         clientName,
		EventDate:          eventDateFormatted,
		EventTime:          eventTimeFormatted,
		EventLocation:      "123 Main St, St. Louis, MO 63110", // Default for preview
		Occasion:           occasion,
		GuestCount:         guestCount,
		Helpers:            helpers,
		Hours:              hours,
		BaseRate:           estimate.BasePerHelper,
		HourlyRate:         estimate.ExtraPerHourPerHelper,
		TotalCost:          estimate.TotalCost,
		DepositAmount:      depositAmount,
		RateLabel:          rateLabel,
		ExpirationDate:     expirationFormatted, // Email needs string
		DepositLink:        depositLink,
		ConfirmationNumber: confirmationNumber,
		IsHighDemand:       estimate.IsSpecialDate, // High demand = special date (holiday/surge)
		UrgencyLevel:       urgencyLevel,           // Urgency level based on days until event
		DaysUntilEvent:     daysUntilEvent,         // Number of days until event
		IsReturningClient:  false,                  // TODO: Check if client has booked before (query CRM/calendar)
		WeatherForecast:    weatherForecast,        // Weather forecast (only for events < 10 days)
		TravelFeeInfo:      travelFeeInfo,          // Travel fee information
	}

	htmlBody := util.GenerateQuoteEmailHTML(emailData)

	// Format date with day of week for subject
	eventDateWithDay := formatDateWithDayOfWeek(emailData.EventDate)

	// Shortened subject line
	subject := fmt.Sprintf("%s Quote - %s", emailData.Occasion, eventDateWithDay)

	// Generate PDF quote for records (always attached for preview too)
	pdfData := util.QuotePDFData{
		ConfirmationNumber: confirmationNumber,
		Occasion:           emailData.Occasion,
		ClientName:         emailData.ClientName,
		ClientEmail:        body.To,
		EventDate:          emailData.EventDate,
		EventTime:          emailData.EventTime,
		HelpersCount:       emailData.Helpers,
		Hours:              emailData.Hours,
		TotalCost:          emailData.TotalCost,
		DepositAmount:      emailData.DepositAmount,
		ExpirationDate:     expirationDate, // PDF needs time.Time
		DepositLink:        depositLink,
		IssueDate:          time.Now(),
	}

	pdfBytes, pdfErr := util.GenerateQuotePDF(pdfData)
	if pdfErr != nil {
		h.logger.Error("failed to generate PDF for preview", "error", pdfErr)
	}

	emailReq := &ports.SendEmailRequest{
		To:       body.To,
		Subject:  "[TEST] " + subject,
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers Team",
	}

	// Attach PDF if generated successfully
	if pdfBytes != nil {
		emailReq.Attachments = []ports.Attachment{
			{
				Filename: fmt.Sprintf("%s-Quote-%s.pdf", emailData.Occasion, confirmationNumber),
				Content:  pdfBytes,
				MimeType: "application/pdf",
			},
		}
	}

	var emailResult *ports.SendEmailResult
	var err error

	if body.SaveAsDraft {
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmailDraft(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmailDraft(r.Context(), emailReq)
		} else {
			util.WriteError(w, http.StatusInternalServerError, "email service is not configured")
			return
		}
	} else {
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmail(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmail(r.Context(), emailReq)
		} else {
			util.WriteError(w, http.StatusInternalServerError, "email service is not configured")
			return
		}
	}

	// #region agent log
	if logFile, logErr := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		errorMsg := ""
		messageID := ""
		success := false
		if err != nil {
			errorMsg = err.Error()
		}
		if emailResult != nil {
			messageID = emailResult.MessageID
			success = emailResult.Success
			if emailResult.Error != nil {
				errorMsg = *emailResult.Error
			}
		}
		json.NewEncoder(logFile).Encode(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A",
			"location":     "email_handler.go:HandleQuoteEmailPreview",
			"message":      "After sending/saving email",
			"data": map[string]interface{}{
				"saveAsDraft": body.SaveAsDraft,
				"to":          body.To,
				"error":       errorMsg,
				"messageID":   messageID,
				"success":     success,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		logFile.Close()
	}
	// #endregion

	if err != nil {
		h.logger.Error("failed to send quote email preview", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to send quote email preview: "+err.Error())
		return
	}

	if emailResult == nil {
		util.WriteError(w, http.StatusInternalServerError, "email service returned nil result")
		return
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("quote email preview sending failed", "error", errorMsg)
		util.WriteError(w, http.StatusInternalServerError, "quote email preview sending failed: "+errorMsg)
		return
	}

	sent := !body.SaveAsDraft
	draft := body.SaveAsDraft

	h.logger.Info("quote email preview sent successfully", "messageId", emailResult.MessageID, "to", body.To, "draft", draft)

	// Try to fetch the sent email from Gmail to confirm it was sent
	var fetchedHTMLBody string
	if !body.SaveAsDraft && h.gmailSender != nil && emailResult.MessageID != "" {
		// Wait 2 seconds to ensure message is available in Gmail
		time.Sleep(2 * time.Second)
		fetched, err := h.gmailSender.GetMessage(r.Context(), emailResult.MessageID)
		if err != nil {
			h.logger.Warn("failed to fetch sent email from Gmail", "error", err, "messageId", emailResult.MessageID)
			// Fall back to generated HTML body
			fetchedHTMLBody = htmlBody
		} else {
			fetchedHTMLBody = fetched
		}
	} else {
		// Use generated HTML body for drafts or if Gmail sender not available
		fetchedHTMLBody = htmlBody
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Quote email preview sent successfully",
		"email": map[string]interface{}{
			"messageId": emailResult.MessageID,
			"sent":      sent,
			"draft":     draft,
			"error":     "",
			"htmlBody":  fetchedHTMLBody, // Include HTML body fetched from Gmail (or generated)
		},
	})
}

// HandleDepositEmailPreview handles POST /api/email/deposit/preview - sends deposit email with dummy data
func (h *EmailHandler) HandleDepositEmailPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		To          string `json:"to"`
		SaveAsDraft bool   `json:"saveAsDraft"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.To == "" {
		util.WriteError(w, http.StatusBadRequest, "to (recipient email) is required")
		return
	}

	// Use dummy data for preview
	dummyName := "John Doe"
	dummyDepositAmount := 50.0
	dummyInvoiceURL := "https://invoice.stripe.com/i/acct_test/live_test?s=ap"

	emailSent, emailError := h.SendDepositEmail(
		r.Context(),
		dummyName,
		body.To,
		dummyDepositAmount,
		dummyInvoiceURL,
		body.SaveAsDraft,
	)

	sent := !body.SaveAsDraft && emailSent
	draft := body.SaveAsDraft

	if emailError != "" {
		h.logger.Error("deposit email preview sending failed", "error", emailError)
		util.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"ok":      true,
			"message": "Deposit email preview sent",
			"email": map[string]interface{}{
				"sent":  sent,
				"draft": draft,
				"error": emailError,
			},
		})
		return
	}

	h.logger.Info("deposit email preview sent successfully", "to", body.To, "draft", draft)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Deposit email preview sent successfully",
		"email": map[string]interface{}{
			"sent":  sent,
			"draft": draft,
			"error": "",
		},
	})
}

// HandleFinalInvoiceEmailPreview handles POST /api/email/final-invoice/preview - sends final invoice email with dummy data
func (h *EmailHandler) HandleFinalInvoiceEmailPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		To          string `json:"to"`
		SaveAsDraft bool   `json:"saveAsDraft"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.To == "" {
		util.WriteError(w, http.StatusBadRequest, "to (recipient email) is required")
		return
	}

	// Use dummy data for preview
	dummyName := "John Doe"
	dummyEventType := "Birthday Party"
	dummyEventDate := "Dec 25, 2025"
	helpersCount := 2
	dummyOriginalQuote := 1000.0
	dummyDepositPaid := 50.0
	dummyRemainingBalance := 950.0
	dummyInvoiceURL := "https://invoice.stripe.com/i/acct_test/live_test?s=ap"
	showGratuity := true

	emailSent, emailError := h.SendFinalInvoiceEmail(
		r.Context(),
		dummyName,
		body.To,
		dummyEventType,
		dummyEventDate,
		&helpersCount,
		dummyOriginalQuote,
		dummyDepositPaid,
		dummyRemainingBalance,
		dummyInvoiceURL,
		showGratuity,
		body.SaveAsDraft,
		"",
	)

	sent := !body.SaveAsDraft && emailSent
	draft := body.SaveAsDraft

	if emailError != "" {
		h.logger.Error("final invoice email preview sending failed", "error", emailError)
		util.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"ok":      true,
			"message": "Final invoice email preview sent",
			"email": map[string]interface{}{
				"sent":  sent,
				"draft": draft,
				"error": emailError,
			},
		})
		return
	}

	h.logger.Info("final invoice email preview sent successfully", "to", body.To, "draft", draft)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Final invoice email preview sent successfully",
		"email": map[string]interface{}{
			"sent":  sent,
			"draft": draft,
			"error": "",
		},
	})
}

// HandleReviewRequestEmailPreview handles POST /api/email/review-request/preview - sends review request email with dummy data
func (h *EmailHandler) HandleReviewRequestEmailPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		To          string `json:"to"`
		SaveAsDraft bool   `json:"saveAsDraft"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.To == "" {
		util.WriteError(w, http.StatusBadRequest, "to (recipient email) is required")
		return
	}

	// Use dummy data for preview
	dummyName := "John Doe"
	dummyReviewURL := "https://g.page/r/test/review"

	templateService := emailService.NewTemplateService()
	htmlBody, err := templateService.GenerateReviewRequestEmail(dummyName, dummyReviewURL)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to generate review request email: "+err.Error())
		return
	}

	emailReq := &ports.SendEmailRequest{
		To:       body.To,
		Subject:  "[TEST] We'd Love Your Feedback - STL Party Helpers",
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers Team",
	}

	var emailResult *ports.SendEmailResult

	if body.SaveAsDraft {
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmailDraft(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmailDraft(r.Context(), emailReq)
		} else {
			util.WriteError(w, http.StatusInternalServerError, "email service is not configured")
			return
		}
	} else {
		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmail(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmail(r.Context(), emailReq)
		} else {
			util.WriteError(w, http.StatusInternalServerError, "email service is not configured")
			return
		}
	}

	if err != nil {
		h.logger.Error("failed to send review request email preview", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to send review request email preview: "+err.Error())
		return
	}

	if emailResult == nil {
		util.WriteError(w, http.StatusInternalServerError, "email service returned nil result")
		return
	}

	if !emailResult.Success {
		errorMsg := "unknown error"
		if emailResult.Error != nil {
			errorMsg = *emailResult.Error
		}
		h.logger.Error("review request email preview sending failed", "error", errorMsg)
		util.WriteError(w, http.StatusInternalServerError, "review request email preview sending failed: "+errorMsg)
		return
	}

	sent := !body.SaveAsDraft
	draft := body.SaveAsDraft

	h.logger.Info("review request email preview sent successfully", "messageId", emailResult.MessageID, "to", body.To, "draft", draft)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Review request email preview sent successfully",
		"email": map[string]interface{}{
			"messageId": emailResult.MessageID,
			"sent":      sent,
			"draft":     draft,
			"error":     "",
		},
	})
}
