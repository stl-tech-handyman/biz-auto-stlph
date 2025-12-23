package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/ports"
	emailService "github.com/bizops360/go-api/internal/services/email"
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
	emailClient *email.EmailServiceClient
	gmailSender *email.GmailSender
	logger      *slog.Logger
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(logger *slog.Logger) *EmailHandler {
	handler := &EmailHandler{
		logger: logger,
	}

	// Try to use email service client first (if EMAIL_SERVICE_URL is set)
	handler.emailClient = email.NewEmailServiceClient()
	if handler.emailClient != nil {
		logger.Info("Using email service API for email sending")
		return handler
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

	return handler
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
	htmlBody, err := templateService.GenerateFinalInvoiceEmail(
		body.Name, "Event", "", nil, body.TotalAmount, body.DepositPaid, body.RemainingBalance, body.InvoiceURL, true)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to generate email template: "+err.Error())
		return
	}

	emailReq := &ports.SendEmailRequest{
		To:       body.Email,
		Subject:  "Final Invoice - STL Party Helpers",
		HTMLBody: htmlBody,
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
func (h *EmailHandler) SendFinalInvoiceEmail(ctx context.Context, name, email, eventType, eventDate string, helpersCount *int, originalQuote, depositPaid, remainingBalance float64, invoiceURL string, showGratuity bool) (bool, string) {
	if name == "" || email == "" || invoiceURL == "" {
		return false, "name, email, and invoiceUrl are required"
	}

	templateService := emailService.NewTemplateService()
	htmlBody, err := templateService.GenerateFinalInvoiceEmail(name, eventType, eventDate, helpersCount, originalQuote, depositPaid, remainingBalance, invoiceURL, showGratuity)
	if err != nil {
		return false, fmt.Sprintf("failed to generate email template: %v", err)
	}

	emailReq := &ports.SendEmailRequest{
		To:       email,
		Subject:  "Final Invoice - STL Party Helpers",
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers",
	}

	var emailResult *ports.SendEmailResult

	if h.gmailSender != nil {
		emailResult, err = h.gmailSender.SendEmail(ctx, emailReq)
	} else if h.emailClient != nil {
		emailResult, err = h.emailClient.SendEmail(ctx, emailReq)
	} else {
		return false, "email service is not configured"
	}

	if err != nil {
		h.logger.Error("failed to send final invoice email", "error", err)
		return false, err.Error()
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
func (h *EmailHandler) SendDepositEmail(ctx context.Context, name, email string, depositAmount float64, invoiceURL string) (bool, string) {
	if name == "" || email == "" || invoiceURL == "" {
		return false, "name, email, and invoiceUrl are required"
	}

	templateService := emailService.NewTemplateService()
	htmlBody, err := templateService.GenerateDepositEmail(name, depositAmount, invoiceURL)
	if err != nil {
		return false, fmt.Sprintf("failed to generate email template: %v", err)
	}

	emailReq := &ports.SendEmailRequest{
		To:       email,
		Subject:  "Booking Deposit - STL Party Helpers",
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers",
	}

	var emailResult *ports.SendEmailResult

	if h.gmailSender != nil {
		emailResult, err = h.gmailSender.SendEmail(ctx, emailReq)
	} else if h.emailClient != nil {
		emailResult, err = h.emailClient.SendEmail(ctx, emailReq)
	} else {
		return false, "email service is not configured"
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
