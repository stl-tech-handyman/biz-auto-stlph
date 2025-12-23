package handlers

import (
	"net/http"

	"github.com/bizops360/go-api/internal/app"
	"github.com/bizops360/go-api/internal/util"
)

// FormEventsHandler handles POST /v1/form-events
type FormEventsHandler struct {
	service *app.FormEventsService
}

// NewFormEventsHandler creates a new form events handler
func NewFormEventsHandler(service *app.FormEventsService) *FormEventsHandler {
	return &FormEventsHandler{service: service}
}

// ServeHTTP handles the HTTP request
func (h *FormEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := r.Context()
	requestID := util.GetRequestID(ctx)

	// Parse request body
	var body struct {
		BusinessID  string         `json:"businessId"`
		PipelineKey string         `json:"pipelineKey"`
		DryRun      bool           `json:"dryRun"`
		Options     map[string]any `json:"options"`
		Fields      map[string]any `json:"fields"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Extract headers (header overrides body)
	businessID := r.Header.Get("X-Business-Id")
	if businessID == "" {
		businessID = body.BusinessID
	}
	if businessID == "" {
		util.WriteError(w, http.StatusBadRequest, "businessId is required")
		return
	}

	pipelineKey := r.Header.Get("X-Pipeline-Key")
	if pipelineKey == "" {
		pipelineKey = body.PipelineKey
	}

	source := r.Header.Get("X-Source")
	if source == "" {
		source = "form"
	}

	dryRun := body.DryRun
	if r.Header.Get("X-Dry-Run") == "true" {
		dryRun = true
	}

	// Build request
	req := &app.FormEventsRequest{
		BusinessID:  businessID,
		PipelineKey: pipelineKey,
		Source:      source,
		DryRun:      dryRun,
		Fields:      body.Fields,
		Options:     body.Options,
		RequestID:   requestID,
	}

	// Execute pipeline
	result, err := h.service.Run(ctx, req)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return result
	util.WriteJSON(w, http.StatusOK, result)
}

