package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/infra/calendar"
	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/infra/geo"
	"github.com/bizops360/go-api/internal/services/lead"
	"github.com/bizops360/go-api/internal/util"
)

// BusinessLeadHandler handles business-specific lead processing
type BusinessLeadHandler struct {
	businessLoader *config.BusinessLoader
	leadProcessor  *lead.Processor
	logger         *slog.Logger
}

// NewBusinessLeadHandler creates a new business lead handler
func NewBusinessLeadHandler(
	businessLoader *config.BusinessLoader,
	logger *slog.Logger,
) *BusinessLeadHandler {
	// Initialize services (these could be injected, but for now we'll create them here)
	// In a more advanced setup, these would be in a service container

	// Initialize calendar service
	calendarID := os.Getenv("ESTIMATE_SENT_CALENDAR_ID")
	if calendarID == "" {
		calendarID = "c_f8c0098141f20b9bcb25d5e3c05d54c450301eb4f21bff9c75a04b1612138b54@group.calendar.google.com"
	}
	calendarService, _ := calendar.NewCalendarService(calendarID)

	// Initialize email service
	emailClient := email.NewEmailServiceClient()
	var gmailSender *email.GmailSender
	if emailClient == nil {
		if gs, err := email.NewGmailSender(); err == nil {
			gmailSender = gs
			logger.Info("Using Gmail API for email sending")
		}
	} else {
		logger.Info("Using email service API for email sending")
	}

	// Initialize geocoding service
	geocodingService, _ := geo.NewGeocodingService()

	// Initialize lead processor
	leadProcessor := lead.NewProcessor(
		calendarService,
		emailClient,
		gmailSender,
		geocodingService,
		logger,
		calendarID,
	)

	return &BusinessLeadHandler{
		businessLoader: businessLoader,
		leadProcessor:  leadProcessor,
		logger:         logger,
	}
}

// HandleProcessLead handles POST /api/business/{businessId}/process-lead
func (h *BusinessLeadHandler) HandleProcessLead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract business ID from URL path
	// Using standard library path parsing since we're using http.ServeMux
	businessID := extractBusinessIDFromPath(r.URL.Path)
	if businessID == "" {
		util.WriteError(w, http.StatusBadRequest, "business ID is required")
		return
	}

	// Load business configuration
	ctx := r.Context()
	businessConfig, err := h.businessLoader.LoadBusiness(ctx, businessID)
	if err != nil {
		h.logger.Warn("business not found", "businessId", businessID, "error", err)
		util.WriteError(w, http.StatusNotFound, "business not found")
		return
	}

	h.logger.Debug("processing lead for business",
		"businessId", businessID,
		"displayName", businessConfig.DisplayName,
	)

	// Parse raw payload
	var rawPayload util.RawZapierPayload
	if err := util.ReadJSON(r, &rawPayload); err != nil {
		h.logger.Warn("invalid JSON payload", "error", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Transform payload (all Zapier steps)
	h.logger.Debug("transforming Zapier payload")
	transformedData, err := util.TransformZapierPayload(rawPayload)
	if err != nil {
		h.logger.Warn("failed to transform payload", "error", err)
		util.WriteError(w, http.StatusBadRequest, "failed to transform payload: "+err.Error())
		return
	}

	h.logger.Info("payload transformed successfully",
		"clientName", transformedData.ClientName,
		"email", transformedData.Email,
		"occasion", transformedData.Occasion,
	)

	// Process lead through workflow
	h.logger.Debug("processing lead through workflow")
	result, err := h.leadProcessor.ProcessLead(ctx, transformedData)
	if err != nil {
		h.logger.Error("failed to process lead", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to process lead: "+err.Error())
		return
	}

	// Build response
	response := map[string]interface{}{
		"referenceNumber": result.ReferenceNumber,
		"success":         result.Success,
		"emailSent":       result.EmailSent,
		"estimate":        result.Estimate,
		"calendarCreated": result.CalendarCreated,
	}

	// Add optional fields (only if they have values)
	if result.EmailError != nil {
		response["emailError"] = *result.EmailError
	}
	if result.CalendarError != nil {
		response["calendarError"] = *result.CalendarError
	}
	if result.EventID != nil {
		response["eventId"] = *result.EventID
	}
	if result.Lat != nil {
		response["lat"] = *result.Lat
	}
	if result.Long != nil {
		response["long"] = *result.Long
	}
	if result.FullAddress != nil {
		response["fullAddress"] = *result.FullAddress
	}
	if result.GeoError != nil {
		response["geoError"] = *result.GeoError
	}

	h.logger.Info("lead processed successfully",
		"businessId", businessID,
		"clientName", transformedData.ClientName,
		"emailSent", result.EmailSent,
		"calendarCreated", result.CalendarCreated,
	)

	util.WriteJSON(w, http.StatusOK, response)
}

// extractBusinessIDFromPath extracts business ID from URL path
// Example: "/api/business/stlpartyhelpers/process-lead" -> "stlpartyhelpers"
func extractBusinessIDFromPath(path string) string {
	// Simple path parsing - in production, consider using a proper router like chi or gorilla/mux
	// For now, we'll parse manually
	prefix := "/api/business/"
	suffix := "/process-lead"
	
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return ""
	}
	
	businessID := path[len(prefix) : len(path)-len(suffix)]
	return businessID
}

