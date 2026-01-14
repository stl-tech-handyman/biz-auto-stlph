package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/infra/calendar"
	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/infra/geo"
	"github.com/bizops360/go-api/internal/infra/stripe"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
)

// ZapierHandler handles POST /api/zapier/process-lead
// Replicates the Apps Script sendEstimateAndAddToCalendarFromZapier flow
type ZapierHandler struct {
	geocodingService *geo.GeocodingService
	calendarService  *calendar.CalendarService
	emailClient      *email.EmailServiceClient
	gmailSender      *email.GmailSender
	logger           *slog.Logger
}

// NewZapierHandler creates a new Zapier handler
func NewZapierHandler(logger *slog.Logger) *ZapierHandler {
	handler := &ZapierHandler{
		logger: logger,
	}

	// Initialize geocoding service
	if geoService, err := geo.NewGeocodingService(); err == nil {
		handler.geocodingService = geoService
		logger.Info("Geocoding service initialized")
	} else {
		logger.Warn("Geocoding service not available", "error", err)
	}

	// Initialize calendar service
	calendarID := os.Getenv("ESTIMATE_SENT_CALENDAR_ID")
	if calendarID == "" {
		calendarID = "c_f8c0098141f20b9bcb25d5e3c05d54c450301eb4f21bff9c75a04b1612138b54@group.calendar.google.com" // Default from Apps Script
	}
	if calService, err := calendar.NewCalendarService(calendarID); err == nil {
		handler.calendarService = calService
		logger.Info("Calendar service initialized")
	} else {
		logger.Warn("Calendar service not available", "error", err)
		// Create a service without credentials (will fail gracefully)
		handler.calendarService, _ = calendar.NewCalendarService(calendarID)
	}

	// Initialize email service
	handler.emailClient = email.NewEmailServiceClient()
	if handler.emailClient != nil {
		logger.Info("Using email service API for email sending")
	} else {
		if gmailSender, err := email.NewGmailSender(); err == nil {
			handler.gmailSender = gmailSender
			logger.Info("Using Gmail API for email sending")
		} else {
			logger.Warn("Email service not available", "error", err)
		}
	}

	return handler
}

// HandleProcessLead handles POST /api/zapier/process-lead
// Matches the Apps Script processNewLeadFromZapier function
func (h *ZapierHandler) HandleProcessLead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var payload struct {
		// Zapier form fields (matching Apps Script)
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		EmailAddress     string `json:"email_address"`
		PhoneNumber      string `json:"phone_number"`
		EventDate        string `json:"event_date"`
		EventTime        string `json:"event_time"`
		EventLocation    string `json:"event_location"`
		HelpersRequested string `json:"helpers_requested"`
		ForHowManyHours  string `json:"for_how_many_hours"`
		Occasion         string `json:"occasion"`
		GuestsExpected   string `json:"guests_expected"`
		DryRun           bool   `json:"dryRun"`
	}

	if err := util.ReadJSON(r, &payload); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Parse and validate input (matching Apps Script logic)
	clientName := fmt.Sprintf("%s %s", strings.TrimSpace(payload.FirstName), strings.TrimSpace(payload.LastName))
	if clientName == "" {
		util.WriteError(w, http.StatusBadRequest, "first_name and last_name are required")
		return
	}

	email := strings.TrimSpace(payload.EmailAddress)
	if email == "" {
		util.WriteError(w, http.StatusBadRequest, "email_address is required")
		return
	}

	// Parse helpers (e.g., "I Need 2 Helpers" -> 2)
	numHelpers := parseHelpers(payload.HelpersRequested)
	if numHelpers <= 0 {
		util.WriteError(w, http.StatusBadRequest, "helpers_requested must contain a valid number")
		return
	}

	// Parse duration (e.g., "for 5 Hours" -> 5.0)
	duration := parseDuration(payload.ForHowManyHours)
	if duration <= 0 {
		util.WriteError(w, http.StatusBadRequest, "for_how_many_hours must contain a valid number")
		return
	}

	// Parse guest count (e.g., "300 - 400 Guests" -> 300)
	guestCount := parseGuestCount(payload.GuestsExpected)

	// Parse event date
	eventDate, err := parseEventDate(payload.EventDate)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid event_date: %v", err))
		return
	}

	// Calculate total cost (matching calculateTotalCost_v2)
	estimate, err := pricing.CalculateEstimate(eventDate, duration, numHelpers)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("failed to calculate estimate: %v", err))
		return
	}

	// Generate quote ID (matching generateShortQuoteID)
	dateKey := eventDate.Format("2006-01-02")
	referenceNumber := util.GenerateShortQuoteID(email, dateKey)

	// Create calendar event
	var eventID string
	var calendarError string
	if h.calendarService != nil {
		calendarReq := &calendar.CreateEventRequest{
			ClientName: clientName,
			Occasion:   payload.Occasion,
			GuestCount: guestCount,
			EventDate:  payload.EventDate, // Use original format, not just dateKey
			EventTime:  payload.EventTime,
			Phone:      payload.PhoneNumber,
			Location:   payload.EventLocation,
			NumHelpers: numHelpers,
			Duration:   duration,
			TotalCost:  estimate.TotalCost,
			EmailID:    email,
			ThreadID:   "",
			DataSource: "zapier",
			Status:     "Pending",
		}

		calendarResult, err := h.calendarService.CreateEvent(r.Context(), calendarReq)
		if err != nil || calendarResult.Error != "" {
			calendarError = err.Error()
			if calendarError == "" {
				calendarError = calendarResult.Error
			}
			h.logger.Warn("Failed to create calendar event", "error", calendarError)
		} else {
			eventID = calendarResult.EventID
		}
	}

	// Send quote email (matching sendQuoteEmailOnly_v2)
	var emailSent bool
	var emailError string
	if h.emailClient != nil || h.gmailSender != nil {
		// Determine rate label
		rateLabel := "Base Rate"
		if estimate.SpecialLabel != nil {
			rateLabel = *estimate.SpecialLabel
		}

		// Calculate deposit from total cost
		estimateCents := util.DollarsToCents(estimate.TotalCost)
		depositCalc := stripe.CalculateDepositFromEstimate(estimateCents)
		depositAmount := util.CentsToDollars(depositCalc.Value)

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
		_, expirationFormatted := util.CalculateExpirationDate(daysUntilEvent)

		// Generate confirmation number
		confirmationNumber := util.GenerateConfirmationNumber(email, payload.Occasion, eventDate)

		// Generate email HTML
		emailData := util.QuoteEmailData{
			ClientName:         clientName,
			EventDate:          formatDateForEmail(eventDate),
			EventTime:          payload.EventTime,
			EventLocation:      payload.EventLocation,
			Occasion:           payload.Occasion,
			GuestCount:         guestCount,
			Helpers:            numHelpers,
			Hours:              duration,
			BaseRate:           estimate.BasePerHelper,
			HourlyRate:         estimate.ExtraPerHourPerHelper,
			TotalCost:          estimate.TotalCost,
			DepositAmount:      depositAmount,
			RateLabel:          rateLabel,
			ExpirationDate:     expirationFormatted,
			DepositLink:        "", // Will be generated when deposit invoice is created
			ConfirmationNumber: confirmationNumber,
			IsHighDemand:       estimate.IsSpecialDate, // High demand = special date (holiday/surge)
			UrgencyLevel:       urgencyLevel,           // Urgency level based on days until event
			DaysUntilEvent:     daysUntilEvent,         // Number of days until event
			IsReturningClient:  false,                  // TODO: Check if client has booked before (query CRM/calendar)
		}

		htmlBody := util.GenerateQuoteEmailHTML(emailData)
		subject := fmt.Sprintf("Party Helpers for %s - %s - Estimate & Details for %s",
			payload.Occasion, formatDateForEmail(eventDate), clientName)

		if payload.DryRun {
			subject = "Dry Run - " + subject
		}

		emailReq := &ports.SendEmailRequest{
			To:       email,
			Subject:  subject,
			HTMLBody: htmlBody,
			FromName: "STL Party Helpers Team",
		}

		// Add CC for non-dry-run
		if !payload.DryRun {
			// Note: Gmail API doesn't support CC in the same way, so we'd need to handle this differently
			// For now, we'll send to the main recipient
		}

		var emailResult *ports.SendEmailResult
		var err error

		if h.gmailSender != nil {
			emailResult, err = h.gmailSender.SendEmail(r.Context(), emailReq)
		} else if h.emailClient != nil {
			emailResult, err = h.emailClient.SendEmail(r.Context(), emailReq)
		}

		if err != nil {
			emailError = err.Error()
			h.logger.Error("Failed to send quote email", "error", emailError)
		} else if !emailResult.Success {
			if emailResult.Error != nil {
				emailError = *emailResult.Error
			} else {
				emailError = "unknown error"
			}
			h.logger.Error("Quote email sending failed", "error", emailError)
		} else {
			emailSent = true
			h.logger.Info("Quote email sent successfully", "to", email)
		}
	} else {
		emailError = "email service not configured"
		h.logger.Warn("Email service not available")
	}

	// Geocode address (matching getLatLng)
	var geoData *geo.GeocodeResult
	if h.geocodingService != nil && payload.EventLocation != "" {
		result, err := h.geocodingService.GetLatLng(r.Context(), payload.EventLocation)
		if err != nil {
			h.logger.Warn("Geocoding failed", "error", err, "address", payload.EventLocation)
			// Don't fail the request if geocoding fails
		} else {
			geoData = result
		}
	}

	// Build response (matching Apps Script response format)
	response := map[string]interface{}{
		"referenceNumber": referenceNumber,
		"success":         true,
		"emailSent":       emailSent,
		"estimate":        estimate.TotalCost,
		"calendarCreated": calendarError == "",
		"calendarError": func() interface{} {
			if calendarError != "" {
				return calendarError
			}
			return nil
		}(),
	}

	if geoData != nil {
		response["lat"] = geoData.Lat
		response["long"] = geoData.Lng
		response["fullAddress"] = geoData.FullAddress
	} else {
		response["lat"] = nil
		response["long"] = nil
		response["fullAddress"] = nil
	}

	if eventID != "" {
		response["eventId"] = eventID
	} else {
		response["eventId"] = nil
	}

	if emailError != "" {
		response["emailError"] = emailError
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// Helper functions to parse Zapier form fields

func parseHelpers(helpersStr string) int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindString(helpersStr)
	if matches == "" {
		return 0
	}
	val, _ := strconv.Atoi(matches)
	return val
}

func parseDuration(durationStr string) float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	matches := re.FindString(durationStr)
	if matches == "" {
		return 0
	}
	val, _ := strconv.ParseFloat(matches, 64)
	return val
}

func parseGuestCount(guestsStr string) int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(guestsStr, -1)
	if len(matches) == 0 {
		return 0
	}
	// Take the first number found
	val, _ := strconv.Atoi(matches[0])
	return val
}

func parseEventDate(dateStr string) (time.Time, error) {
	// Try multiple date formats
	formats := []string{
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
		"2006/01/02",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func formatDateForEmail(date time.Time) string {
	// Format as "Fri, Jan 19, 2026" (day of week, short month, day, year)
	return date.Format("Mon, Jan 2, 2006")
}
