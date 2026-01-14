package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bizops360/go-api/internal/services/email_analysis"
	"github.com/bizops360/go-api/internal/util"
)

// EmailAnalysisHandler handles email analysis requests
type EmailAnalysisHandler struct {
	service *email_analysis.Service
	logger  *slog.Logger
}

// NewEmailAnalysisHandler creates a new handler
func NewEmailAnalysisHandler(logger *slog.Logger) (*EmailAnalysisHandler, error) {
	service, err := email_analysis.NewService(logger)
	if err != nil {
		return nil, err
	}

	return &EmailAnalysisHandler{
		service: service,
		logger:  logger,
	}, nil
}

// HandleAnalyze processes email analysis request
func (h *EmailAnalysisHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req email_analysis.AnalyzeEmailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Default values
	if req.MaxEmails == 0 {
		req.MaxEmails = 100 // Default to 100 emails for safety
	}

	ctx := r.Context()
	resp, err := h.service.AnalyzeEmails(ctx, req)
	if err != nil {
		h.logger.Error("failed to analyze emails", "error", err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// HandleStatus returns analysis status
func (h *EmailAnalysisHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	spreadsheetID := r.URL.Query().Get("spreadsheet_id")
	ctx := r.Context()
	state, err := h.service.GetStatus(ctx, spreadsheetID)
	if err != nil {
		h.logger.Error("failed to get status", "error", err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, state)
}
