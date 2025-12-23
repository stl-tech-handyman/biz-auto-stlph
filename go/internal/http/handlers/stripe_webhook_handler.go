package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/util"
)

// StripeWebhookHandler handles Stripe webhook events
type StripeWebhookHandler struct {
	paymentsProvider ports.PaymentsProvider
	emailClient      *email.EmailServiceClient
	gmailSender      *email.GmailSender
	logger           *slog.Logger
}

// NewStripeWebhookHandler creates a new Stripe webhook handler
func NewStripeWebhookHandler(
	paymentsProvider ports.PaymentsProvider,
	emailClient *email.EmailServiceClient,
	gmailSender *email.GmailSender,
	logger *slog.Logger,
) *StripeWebhookHandler {
	return &StripeWebhookHandler{
		paymentsProvider: paymentsProvider,
		emailClient:      emailClient,
		gmailSender:      gmailSender,
		logger:           logger,
	}
}

// StripeWebhookEvent represents a Stripe webhook event
type StripeWebhookEvent struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Data    StripeWebhookEventData `json:"data"`
	Created int64                  `json:"created"`
}

// StripeWebhookEventData contains the event data
type StripeWebhookEventData struct {
	Object StripeInvoiceObject `json:"object"`
}

// StripeInvoiceObject represents a Stripe invoice in webhook events
type StripeInvoiceObject struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	AmountDue       int64  `json:"amount_due"`
	AmountPaid      int64  `json:"amount_paid"`
	Customer        string `json:"customer"`
	CustomerEmail   string `json:"customer_email"`
	CustomerName    string `json:"customer_name"`
	HostedInvoiceURL string `json:"hosted_invoice_url"`
	InvoicePDF      string `json:"invoice_pdf"`
	Metadata        map[string]string `json:"metadata"`
}

// HandleWebhook handles POST /api/stripe/webhook
func (h *StripeWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var event StripeWebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Warn("failed to decode webhook event", "error", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	h.logger.Info("received Stripe webhook event",
		"event_id", event.ID,
		"event_type", event.Type,
		"invoice_id", event.Data.Object.ID,
	)

	// Handle different event types
	switch event.Type {
	case "invoice.paid":
		h.handleInvoicePaid(r.Context(), &event.Data.Object)
	case "invoice.payment_succeeded":
		h.handleInvoicePaid(r.Context(), &event.Data.Object)
	default:
		h.logger.Debug("unhandled webhook event type", "type", event.Type)
	}

	// Always return 200 to acknowledge receipt
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"message": "Webhook received",
		"eventId": event.ID,
	})
}

// handleInvoicePaid processes when an invoice is paid
func (h *StripeWebhookHandler) handleInvoicePaid(ctx context.Context, invoice *StripeInvoiceObject) {
	h.logger.Info("processing invoice paid event",
		"invoice_id", invoice.ID,
		"customer_email", invoice.CustomerEmail,
		"amount_paid", invoice.AmountPaid,
	)

	// Determine invoice type from metadata
	invoiceType := invoice.Metadata["invoice_type"]
	if invoiceType == "" {
		invoiceType = "unknown"
	}

	switch invoiceType {
	case "deposit", "booking_deposit":
		h.handleBookingDepositPaid(ctx, invoice)
	case "final":
		h.handleFinalInvoicePaid(ctx, invoice)
	default:
		h.logger.Info("invoice paid but type not specified, treating as generic",
			"invoice_id", invoice.ID,
			"invoice_type", invoiceType,
		)
		// Could send a generic confirmation email here
	}
}

// handleBookingDepositPaid handles when a booking deposit is paid
func (h *StripeWebhookHandler) handleBookingDepositPaid(ctx context.Context, invoice *StripeInvoiceObject) {
	h.logger.Info("booking deposit paid",
		"invoice_id", invoice.ID,
		"customer_email", invoice.CustomerEmail,
		"amount_paid", invoice.AmountPaid,
	)

	// Send confirmation email
	if invoice.CustomerEmail != "" {
		customerName := invoice.CustomerName
		if customerName == "" {
			customerName = "Valued Customer"
		}

		emailReq := &ports.SendEmailRequest{
			To:      invoice.CustomerEmail,
			Subject: "Booking Deposit Received - STL Party Helpers",
			HTMLBody: generateBookingDepositConfirmationEmail(customerName, invoice),
			FromName: "STL Party Helpers",
		}

		var result *ports.SendEmailResult
		var err error

		if h.gmailSender != nil {
			result, err = h.gmailSender.SendEmail(ctx, emailReq)
		} else if h.emailClient != nil {
			result, err = h.emailClient.SendEmail(ctx, emailReq)
		}

		if err != nil {
			h.logger.Error("failed to send booking deposit confirmation email",
				"error", err,
				"invoice_id", invoice.ID,
			)
		} else if result != nil && result.Success {
			h.logger.Info("booking deposit confirmation email sent",
				"invoice_id", invoice.ID,
				"message_id", result.MessageID,
			)
		}
	}

	// Here you could also:
	// - Update CRM (Monday.com, etc.)
	// - Send Slack notification
	// - Update calendar event status
	// - Trigger other workflows
}

// handleFinalInvoicePaid handles when a final invoice is paid
func (h *StripeWebhookHandler) handleFinalInvoicePaid(ctx context.Context, invoice *StripeInvoiceObject) {
	h.logger.Info("final invoice paid",
		"invoice_id", invoice.ID,
		"customer_email", invoice.CustomerEmail,
		"amount_paid", invoice.AmountPaid,
	)

	// Send thank you email
	if invoice.CustomerEmail != "" {
		customerName := invoice.CustomerName
		if customerName == "" {
			customerName = "Valued Customer"
		}

		emailReq := &ports.SendEmailRequest{
			To:      invoice.CustomerEmail,
			Subject: "Thank You - Final Payment Received - STL Party Helpers",
			HTMLBody: generateFinalInvoicePaidEmail(customerName, invoice),
			FromName: "STL Party Helpers",
		}

		var result *ports.SendEmailResult
		var err error

		if h.gmailSender != nil {
			result, err = h.gmailSender.SendEmail(ctx, emailReq)
		} else if h.emailClient != nil {
			result, err = h.emailClient.SendEmail(ctx, emailReq)
		}

		if err != nil {
			h.logger.Error("failed to send final invoice paid email",
				"error", err,
				"invoice_id", invoice.ID,
			)
		} else if result != nil && result.Success {
			h.logger.Info("final invoice paid email sent",
				"invoice_id", invoice.ID,
				"message_id", result.MessageID,
			)
		}
	}

	// Here you could also:
	// - Update CRM status to "Paid in Full"
	// - Send receipt
	// - Archive event
	// - Trigger follow-up workflows
}

// generateBookingDepositConfirmationEmail generates HTML for booking deposit confirmation
func generateBookingDepositConfirmationEmail(name string, invoice *StripeInvoiceObject) string {
	amountPaid := float64(invoice.AmountPaid) / 100
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Booking Deposit Received</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #0047ab;">Hello ` + name + `!</h1>
        
        <p>Thank you for your booking deposit payment!</p>
        
        <div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h2 style="margin-top: 0;">Payment Confirmation</h2>
            <p><strong>Invoice ID:</strong> ` + invoice.ID + `</p>
            <p><strong>Amount Paid:</strong> $` + formatCurrency(amountPaid) + `</p>
            <p><strong>Status:</strong> <span style="color: green; font-weight: bold;">Paid</span></p>
        </div>
        
        <p>Your event reservation is now confirmed. We'll be in touch soon with more details about your event.</p>
        
        <p style="font-size: 0.9em; color: #666;">
            If you have any questions, please don't hesitate to contact us.
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
</html>`
}

// generateFinalInvoicePaidEmail generates HTML for final invoice paid confirmation
func generateFinalInvoicePaidEmail(name string, invoice *StripeInvoiceObject) string {
	amountPaid := float64(invoice.AmountPaid) / 100
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Final Payment Received</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #0047ab;">Hello ` + name + `!</h1>
        
        <p>Thank you for your final payment!</p>
        
        <div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h2 style="margin-top: 0;">Payment Confirmation</h2>
            <p><strong>Invoice ID:</strong> ` + invoice.ID + `</p>
            <p><strong>Amount Paid:</strong> $` + formatCurrency(amountPaid) + `</p>
            <p><strong>Status:</strong> <span style="color: green; font-weight: bold;">Paid in Full</span></p>
        </div>
        
        <p>We hope you enjoyed your event with STL Party Helpers! We appreciate your business and look forward to serving you again in the future.</p>
        
        <p style="font-size: 0.9em; color: #666;">
            If you have any questions or feedback, please don't hesitate to contact us.
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
</html>`
}

// formatCurrency formats a float as currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

