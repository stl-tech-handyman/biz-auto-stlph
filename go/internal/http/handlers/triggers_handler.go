package handlers

import (
	"net/http"

	"github.com/bizops360/go-api/internal/app"
	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/util"
)

// TriggersHandler handles POST /v1/triggers
type TriggersHandler struct {
	service *app.TriggersService
}

// NewTriggersHandler creates a new triggers handler
func NewTriggersHandler(service *app.TriggersService) *TriggersHandler {
	return &TriggersHandler{service: service}
}

// ServeHTTP handles the HTTP request
func (h *TriggersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := r.Context()
	requestID := util.GetRequestID(ctx)

	// Parse request body
	var body struct {
		Source      string                 `json:"source"`
		BusinessID  string                 `json:"businessId"`
		TriggerKey  string                 `json:"triggerKey"`
		PipelineKey string                 `json:"pipelineKey"`
		Resource    *domain.ResourceContext `json:"resource"`
		Payload     map[string]any         `json:"payload"`
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

	triggerKey := r.Header.Get("X-Trigger-Key")
	if triggerKey == "" {
		triggerKey = body.TriggerKey
	}

	pipelineKey := r.Header.Get("X-Pipeline-Key")
	if pipelineKey == "" {
		pipelineKey = body.PipelineKey
	}

	source := r.Header.Get("X-Source")
	if source == "" {
		source = body.Source
	}
	if source == "" {
		source = "trigger"
	}

	dryRun := r.Header.Get("X-Dry-Run") == "true"

	// Build request
	req := &app.TriggerRequest{
		BusinessID:  businessID,
		TriggerKey:  triggerKey,
		PipelineKey: pipelineKey,
		Source:      source,
		Resource:    body.Resource,
		Payload:     body.Payload,
		DryRun:      dryRun,
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

