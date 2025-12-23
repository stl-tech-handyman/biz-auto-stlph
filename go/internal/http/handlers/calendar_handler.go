package handlers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bizops360/go-api/internal/infra/calendar"
	"github.com/bizops360/go-api/internal/util"
)

// CalendarHandler handles calendar-related endpoints
type CalendarHandler struct {
	calendarService *calendar.CalendarService
	logger          *slog.Logger
}

// NewCalendarHandler creates a new calendar handler
func NewCalendarHandler(logger *slog.Logger) *CalendarHandler {
	// Initialize calendar service with default calendar ID
	calendarID := os.Getenv("ESTIMATE_SENT_CALENDAR_ID")
	if calendarID == "" {
		calendarID = "c_f8c0098141f20b9bcb25d5e3c05d54c450301eb4f21bff9c75a04b1612138b54@group.calendar.google.com" // Default from Apps Script
	}

	calService, err := calendar.NewCalendarService(calendarID)
	if err != nil {
		logger.Warn("Calendar service not available", "error", err)
		// Create service without credentials (will fail gracefully)
		calService, _ = calendar.NewCalendarService(calendarID)
	} else {
		logger.Info("Calendar service initialized")
	}

	return &CalendarHandler{
		calendarService: calService,
		logger:          logger,
	}
}

// HandleCreate handles POST /api/calendar/create
// Matches the Apps Script createEvent function
func (h *CalendarHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		CalendarID string  `json:"calendarId"` // Optional, overrides default
		ClientName string  `json:"clientName"`
		Occasion   string  `json:"occasion"`
		GuestCount int     `json:"guestCount"`
		EventDate  string  `json:"eventDate"` // Required
		EventTime  string  `json:"eventTime"` // Required
		Phone      string  `json:"phone"`
		Location   string  `json:"location"`
		NumHelpers int     `json:"numHelpers"`
		Duration   float64 `json:"duration"` // Required, in hours
		TotalCost  float64 `json:"totalCost"`
		EmailID    string  `json:"emailId"`
		ThreadID   string  `json:"threadId"`
		DataSource string  `json:"dataSource"`
		Status     string  `json:"status"` // Optional, defaults to "Pending"
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if body.EventDate == "" {
		util.WriteError(w, http.StatusBadRequest, "eventDate is required")
		return
	}
	if body.EventTime == "" {
		util.WriteError(w, http.StatusBadRequest, "eventTime is required")
		return
	}
	if body.Duration < 1 {
		util.WriteError(w, http.StatusBadRequest, "duration must be at least 1 hour")
		return
	}

	// Build request
	req := &calendar.CreateEventRequest{
		CalendarID: body.CalendarID,
		ClientName: body.ClientName,
		Occasion:   body.Occasion,
		GuestCount: body.GuestCount,
		EventDate:  body.EventDate,
		EventTime:  body.EventTime,
		Phone:      body.Phone,
		Location:   body.Location,
		NumHelpers: body.NumHelpers,
		Duration:   body.Duration,
		TotalCost:  body.TotalCost,
		EmailID:    body.EmailID,
		ThreadID:   body.ThreadID,
		DataSource: body.DataSource,
		Status:     body.Status,
	}

	// Create event
	result, err := h.calendarService.CreateEvent(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create calendar event", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to create calendar event: "+err.Error())
		return
	}

	// Check for error in result
	if result.Error != "" {
		h.logger.Warn("Calendar event creation failed", "error", result.Error)
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   result.Error,
		})
		return
	}

	// Success response (matching Apps Script format)
	response := map[string]interface{}{
		"success": true,
		"eventId": result.EventID,
	}

	h.logger.Info("Calendar event created", "eventId", result.EventID)
	util.WriteJSON(w, http.StatusOK, response)
}

