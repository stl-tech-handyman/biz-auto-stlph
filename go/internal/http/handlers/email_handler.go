package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/ports"
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
	if gmailSender, err := email.NewGmailSender(); err == nil {
		handler.gmailSender = gmailSender
		logger.Info("Using Gmail API for email sending")
	} else {
		logger.Warn("Gmail API not available", "error", err)
		logger.Warn("Email functionality requires EMAIL_SERVICE_URL or GMAIL_CREDENTIALS_JSON to be configured")
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
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Final Invoice - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #0047ab;">Hello %s!</h1>
        
        <p>Thank you for your business with STL Party Helpers!</p>
        
        <p>Your event has been completed. Please find your final invoice below for the remaining balance.</p>
        
        <div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h2 style="margin-top: 0;">Invoice Details</h2>
            <p><strong>Total Event Cost:</strong> $%.2f</p>
            <p><strong>Deposit Paid:</strong> $%.2f</p>
            <p><strong>Remaining Balance:</strong> <strong style="color: #0047ab; font-size: 1.2em;">$%.2f</strong></p>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="%s" 
               style="display: inline-block; background-color: #0047ab; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;">
                Pay Final Invoice
            </a>
        </p>
        
        <p style="font-size: 0.9em; color: #666;">
            If you have any questions about this invoice, please don't hesitate to contact us.
        </p>
        
        <hr style="border: none; border-top: 1px solid #ddd; margin: 30px 0;">
        
        <p style="font-size: 0.85em; color: #666; text-align: center;">
            STL Party Helpers<br>
            4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br>
            <a href="tel:+13147145514" style="color: #0047ab;">(314) 714-5514</a><br>
            <a href="https://stlpartyhelpers.com" style="color: #0047ab;">stlpartyhelpers.com</a>
        </p>
    </div>
</body>
</html>`, body.Name, body.TotalAmount, body.DepositPaid, body.RemainingBalance, body.InvoiceURL)

	emailReq := &ports.SendEmailRequest{
		To:       body.Email,
		Subject:  "Final Invoice - STL Party Helpers",
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers",
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


