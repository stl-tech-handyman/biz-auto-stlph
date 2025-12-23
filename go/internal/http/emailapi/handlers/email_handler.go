package handlers

import (
	"log/slog"
	"net/http"

	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/util"
)

// EmailHandler handles email-related endpoints for the email API
type EmailHandler struct {
	gmailSender *email.GmailSender
	logger      *slog.Logger
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(gmailSender *email.GmailSender, logger *slog.Logger) *EmailHandler {
	return &EmailHandler{
		gmailSender: gmailSender,
		logger:      logger,
	}
}

// HandleRoot handles GET /
func (h *EmailHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service":       "bizops360-email-api",
		"version":       "1.0.0",
		"status":        "running",
		"documentation": "/api/health",
	})
}

// HandleHealth handles GET /api/health
func (h *EmailHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "bizops360-email-api",
	})
}

// HandleSend handles POST /api/email/send
func (h *EmailHandler) HandleSend(w http.ResponseWriter, r *http.Request) {
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

	if body.Subject == "" {
		util.WriteError(w, http.StatusBadRequest, "subject is required")
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

	result, err := h.gmailSender.SendEmail(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to send email", "error", err, "to", req.To)
		util.WriteError(w, http.StatusInternalServerError, "failed to send email: "+err.Error())
		return
	}

	if !result.Success {
		errorMsg := "unknown error"
		if result.Error != nil {
			errorMsg = *result.Error
		}
		h.logger.Error("email sending failed", "error", errorMsg, "to", req.To)
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

// HandleTest handles POST /api/email/test (test endpoint for compatibility)
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

	if body.Subject == "" {
		body.Subject = "Test Email"
	}

	req := &ports.SendEmailRequest{
		To:       body.To,
		Subject:  body.Subject,
		HTMLBody: body.HTML,
		TextBody: body.Text,
		From:     body.From,
		FromName: body.FromName,
	}

	result, err := h.gmailSender.SendEmail(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to send email", "error", err, "to", req.To)
		util.WriteError(w, http.StatusInternalServerError, "failed to send email: "+err.Error())
		return
	}

	if !result.Success {
		errorMsg := "unknown error"
		if result.Error != nil {
			errorMsg = *result.Error
		}
		h.logger.Error("email sending failed", "error", errorMsg, "to", req.To)
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

// HandleDraft handles POST /api/email/draft
func (h *EmailHandler) HandleDraft(w http.ResponseWriter, r *http.Request) {
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

	if body.Subject == "" {
		util.WriteError(w, http.StatusBadRequest, "subject is required")
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

	result, err := h.gmailSender.SendEmailDraft(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to create email draft", "error", err, "to", req.To)
		util.WriteError(w, http.StatusInternalServerError, "failed to create email draft: "+err.Error())
		return
	}

	if !result.Success {
		errorMsg := "unknown error"
		if result.Error != nil {
			errorMsg = *result.Error
		}
		h.logger.Error("email draft creation failed", "error", errorMsg, "to", req.To)
		util.WriteError(w, http.StatusInternalServerError, "email draft creation failed: "+errorMsg)
		return
	}

	h.logger.Info("email draft created successfully", "messageId", result.MessageID, "to", req.To)
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Email draft created successfully",
		"result": map[string]interface{}{
			"messageId": result.MessageID,
			"success":   result.Success,
		},
	})
}

