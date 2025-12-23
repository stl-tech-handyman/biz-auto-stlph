package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
)

// StripeHandler handles Stripe-related endpoints
type StripeHandler struct {
	paymentsProvider ports.PaymentsProvider
}

// NewStripeHandler creates a new Stripe handler
func NewStripeHandler(paymentsProvider ports.PaymentsProvider) *StripeHandler {
	return &StripeHandler{
		paymentsProvider: paymentsProvider,
	}
}

// HandleDeposit handles POST /api/stripe/deposit
func (h *StripeHandler) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Email          string  `json:"email"`
		Name           string  `json:"name"`
		EstimatedTotal *float64 `json:"estimatedTotal"`
		DepositValue   *float64 `json:"depositValue"`
		Deposit        *float64 `json:"deposit"`
		HelpersCount   *int    `json:"helpersCount"`
		Hours          *float64 `json:"hours"`
		UseTest        bool    `json:"useTest"`
		DryRun         bool    `json:"dryRun"`
		MockStripe     bool    `json:"mockStripe"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Convert dollars to cents if needed
	var depositCents *int64
	if body.DepositValue != nil {
		cents := int64(*body.DepositValue * 100)
		depositCents = &cents
	} else if body.Deposit != nil {
		cents := int64(*body.Deposit * 100)
		depositCents = &cents
	}

	var estimateCents *int64
	if body.EstimatedTotal != nil {
		cents := int64(*body.EstimatedTotal * 100)
		estimateCents = &cents
	}

	// For now, return a simple response (full implementation would call Stripe)
	response := map[string]interface{}{
		"ok": true,
		"message": "Deposit calculation (stub - full Stripe integration pending)",
	}

	if depositCents != nil {
		response["deposit"] = map[string]interface{}{
			"value": *depositCents,
			"valueDollars": float64(*depositCents) / 100,
		}
	}

	if estimateCents != nil {
		deposit, _ := h.paymentsProvider.CalculateDeposit(r.Context(), *estimateCents)
		response["deposit"] = map[string]interface{}{
			"value": deposit.AmountCents,
			"valueDollars": deposit.AmountDollars,
			"percentage": deposit.Percentage,
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleDepositCalculate handles GET /api/stripe/deposit/calculate
func (h *StripeHandler) HandleDepositCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	estimateStr := r.URL.Query().Get("estimate")
	depositStr := r.URL.Query().Get("deposit")
	_ = r.URL.Query().Get("show_table") == "true" // Reserved for future use

	var estimateCents *int64
	var depositCents *int64

	if depositStr != "" {
		depositDollars, err := strconv.ParseFloat(depositStr, 64)
		if err == nil {
			cents := int64(depositDollars * 100)
			depositCents = &cents
		}
	}

	if estimateStr != "" {
		estimateDollars, err := strconv.ParseFloat(estimateStr, 64)
		if err == nil {
			cents := int64(estimateDollars * 100)
			estimateCents = &cents
		}
	}

	if depositCents == nil && estimateCents == nil {
		util.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"ok": true,
			"message": "No estimate or deposit provided.",
			"usage": "Add ?estimate=100 to calculate deposit, or ?deposit=50 to set manually (both in dollars)",
		})
		return
	}

	response := map[string]interface{}{
		"ok": true,
	}

	if depositCents != nil {
		response["deposit"] = float64(*depositCents) / 100
		response["depositCents"] = *depositCents
		response["pickedBy"] = "manual"
		response["isManualOverride"] = true
	}

	if estimateCents != nil {
		calc, _ := h.paymentsProvider.CalculateDeposit(r.Context(), *estimateCents)
		if depositCents == nil {
			response["deposit"] = calc.AmountDollars
			response["depositCents"] = calc.AmountCents
			response["pickedBy"] = "calculated"
			response["isManualOverride"] = false
		}
		response["requested_estimate"] = float64(*estimateCents) / 100
		response["calculation"] = map[string]interface{}{
			"deposit": calc.AmountDollars,
			"percentage": calc.Percentage,
			"min_range": float64(calc.AmountCents - 1000) / 100, // Simplified
			"max_range": float64(calc.AmountCents + 1000) / 100, // Simplified
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleDepositWithEmail handles POST /api/stripe/deposit/with-email
func (h *StripeHandler) HandleDepositWithEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Name              string   `json:"name"`
		Email             string   `json:"email"`
		EventType         string   `json:"eventType"`
		EventDateTimeLocal string  `json:"eventDateTimeLocal"`
		EventDate         string   `json:"eventDate"`
		HelpersCount      *int     `json:"helpersCount"`
		Hours             *float64 `json:"hours"`
		Duration          *float64 `json:"duration"`
		Estimate          *float64 `json:"estimate"`
		EstimatedTotal    *float64 `json:"estimatedTotal"`
		DepositValue      *float64 `json:"depositValue"`
		UseTest           bool     `json:"useTest"`
		DryRun            bool     `json:"dryRun"`
		SaveAsDraft       bool     `json:"saveAsDraft"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.Name == "" {
		util.WriteError(w, http.StatusBadRequest, "Missing required field: name")
		return
	}

	// Calculate estimate if event details provided
	var estimateResult *pricing.EstimateResult
	if body.EventDateTimeLocal != "" || body.EventDate != "" {
		eventDateStr := body.EventDateTimeLocal
		if eventDateStr == "" {
			eventDateStr = body.EventDate
		}
		
		eventDate, err := time.Parse("2006-01-02", eventDateStr[:10])
		if err == nil {
			durationHours := 4.0
			if body.Hours != nil {
				durationHours = *body.Hours
			} else if body.Duration != nil {
				durationHours = *body.Duration
			}
			
			numHelpers := 2
			if body.HelpersCount != nil {
				numHelpers = *body.HelpersCount
			}
			
			estimateResult, _ = pricing.CalculateEstimate(eventDate, durationHours, numHelpers)
		}
	}

	// Determine deposit
	var depositCents int64
	if body.DepositValue != nil {
		depositCents = int64(*body.DepositValue * 100)
	} else if estimateResult != nil {
		deposit, _ := h.paymentsProvider.CalculateDeposit(r.Context(), int64(estimateResult.TotalCost * 100))
		depositCents = deposit.AmountCents
	} else if body.Estimate != nil {
		deposit, _ := h.paymentsProvider.CalculateDeposit(r.Context(), int64(*body.Estimate * 100))
		depositCents = deposit.AmountCents
	} else if body.EstimatedTotal != nil {
		deposit, _ := h.paymentsProvider.CalculateDeposit(r.Context(), int64(*body.EstimatedTotal * 100))
		depositCents = deposit.AmountCents
	}

	response := map[string]interface{}{
		"ok": true,
		"message": "Invoice generated (stub - full Stripe integration pending)",
		"dryRun": body.DryRun,
		"saveAsDraft": body.SaveAsDraft,
		"generatedInvoice": map[string]interface{}{
			"id": "stub_invoice_id",
			"url": "https://stripe.com/invoice/stub",
			"amount": float64(depositCents) / 100,
		},
	}

	if estimateResult != nil {
		response["basedOnEstimate"] = estimateResult
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleTest handles POST /api/stripe/test
func (h *StripeHandler) HandleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get query parameters (can override with body)
	useTest := r.URL.Query().Get("useTest") == "true" || r.URL.Query().Get("use_test") == "true"
	sendEmail := r.URL.Query().Get("sendEmail") == "true" || r.URL.Query().Get("send_email") == "true"

	var body map[string]interface{}
	if err := util.ReadJSON(r, &body); err != nil {
		// Body is optional, continue with defaults
		body = make(map[string]interface{})
	}

	// Helper to get param from query or body
	getParam := func(queryKey, bodyKey string, defaultValue interface{}) interface{} {
		if val := r.URL.Query().Get(queryKey); val != "" {
			return val
		}
		if val, ok := body[bodyKey]; ok {
			return val
		}
		return defaultValue
	}

	// Parse useTest and sendEmail from body if present
	if val, ok := body["useTest"].(bool); ok {
		useTest = val
	}
	if val, ok := body["sendEmail"].(bool); ok {
		sendEmail = val
	}

	// Default demo data - convert to cents
	estimatedTotal := 10.0 // $10.00 default
	if val := getParam("estimate", "estimatedTotal", nil); val != nil {
		if f, ok := val.(float64); ok {
			estimatedTotal = f
		} else if s, ok := val.(string); ok {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				estimatedTotal = f
			}
		}
	}

	// Convert dollars to cents if < 10000
	if estimatedTotal < 10000 {
		estimatedTotal = estimatedTotal * 100
	}

	depositValue := getParam("depositValue", "depositValue", nil)
	var depositValueFloat *float64
	if depositValue != nil {
		if f, ok := depositValue.(float64); ok {
			depositValueFloat = &f
			if *depositValueFloat < 10000 {
				*depositValueFloat = *depositValueFloat * 100
			}
		} else if s, ok := depositValue.(string); ok {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				depositValueFloat = &f
				if *depositValueFloat < 10000 {
					*depositValueFloat = *depositValueFloat * 100
				}
			}
		}
	}

	// Create invoice request
	invoiceReq := &ports.CreateInvoiceRequest{
		CustomerEmail: func() string {
			if email := getParam("email", "email", "alexey@shevelyov.com"); email != nil {
				if s, ok := email.(string); ok {
					return s
				}
			}
			return "alexey@shevelyov.com"
		}(),
		CustomerName: func() string {
			if name := getParam("name", "name", "Test Customer"); name != nil {
				if s, ok := name.(string); ok {
					return s
				}
			}
			return "Test Customer"
		}(),
		AmountCents: int64(estimatedTotal),
		Currency:    "usd",
		Description: "Test Booking Deposit Invoice",
	}

	if depositValueFloat != nil {
		invoiceReq.AmountCents = int64(*depositValueFloat)
	}

	// Generate invoice
	invoiceResult, err := h.paymentsProvider.CreateInvoice(r.Context(), invoiceReq)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invoice generation failed: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"ok":      true,
		"message": "Test invoice generated successfully",
		"test": map[string]interface{}{
			"useTest":   useTest,
			"sendEmail": sendEmail,
		},
		"generatedInvoice": map[string]interface{}{
			"id":     invoiceResult.InvoiceID,
			"url":    invoiceResult.HostedInvoiceURL,
			"amount": float64(invoiceResult.AmountDue) / 100,
		},
		"demoData": map[string]interface{}{
			"email":          invoiceReq.CustomerEmail,
			"name":           invoiceReq.CustomerName,
			"estimatedTotal": float64(invoiceReq.AmountCents) / 100,
			"depositValue": func() interface{} {
				if depositValueFloat != nil {
					return *depositValueFloat / 100
				}
				return nil
			}(),
			"useTest": useTest,
		},
	}

	// Optionally send email (stub for now)
	if sendEmail {
		response["email"] = map[string]interface{}{
			"sent":    false,
			"message": "Email sending not yet implemented in Go API",
			"note":    "Use /api/stripe/deposit/with-email for full email integration",
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleFinalInvoice handles POST /api/stripe/final-invoice
func (h *StripeHandler) HandleFinalInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Email            string            `json:"email"`
		Name             string            `json:"name"`
		TotalAmountCents *int64            `json:"totalAmountCents"`
		TotalAmount      *float64          `json:"totalAmount"` // in dollars
		DepositPaidCents *int64            `json:"depositPaidCents"`
		DepositPaid      *float64          `json:"depositPaid"` // in dollars
		Currency         string            `json:"currency"`
		Description      string            `json:"description"`
		Metadata         map[string]string `json:"metadata"`
		UseTest          bool              `json:"useTest"`
		SendEmail        bool              `json:"sendEmail"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.Email == "" {
		util.WriteError(w, http.StatusBadRequest, "email is required")
		return
	}

	if body.Name == "" {
		util.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	// Convert dollars to cents if provided
	var totalCents int64
	if body.TotalAmountCents != nil {
		totalCents = *body.TotalAmountCents
	} else if body.TotalAmount != nil {
		totalCents = int64(*body.TotalAmount * 100)
	} else {
		util.WriteError(w, http.StatusBadRequest, "totalAmount or totalAmountCents is required")
		return
	}

	var depositPaidCents int64
	if body.DepositPaidCents != nil {
		depositPaidCents = *body.DepositPaidCents
	} else if body.DepositPaid != nil {
		depositPaidCents = int64(*body.DepositPaid * 100)
	} else {
		util.WriteError(w, http.StatusBadRequest, "depositPaid or depositPaidCents is required")
		return
	}

	if body.Currency == "" {
		body.Currency = "usd"
	}

	// Create final invoice request
	req := &ports.CreateFinalInvoiceRequest{
		CustomerEmail:     body.Email,
		CustomerName:      body.Name,
		TotalAmountCents:  totalCents,
		DepositPaidCents:  depositPaidCents,
		Currency:          body.Currency,
		Description:       body.Description,
		Metadata:          body.Metadata,
	}

	// Create final invoice
	invoiceResult, err := h.paymentsProvider.CreateFinalInvoice(r.Context(), req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Failed to create final invoice: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"ok": true,
		"message": "Final invoice created successfully",
		"invoice": map[string]interface{}{
			"id":     invoiceResult.InvoiceID,
			"url":    invoiceResult.HostedInvoiceURL,
			"amount": float64(invoiceResult.AmountDue) / 100,
			"status": invoiceResult.Status,
			"pdf":    invoiceResult.InvoicePDF,
		},
		"details": map[string]interface{}{
			"totalAmount":    float64(totalCents) / 100,
			"depositPaid":    float64(depositPaidCents) / 100,
			"remainingBalance": float64(invoiceResult.AmountDue) / 100,
		},
	}

	// Optionally send invoice via Stripe
	if body.SendEmail {
		if err := h.paymentsProvider.SendInvoice(r.Context(), invoiceResult.InvoiceID, body.UseTest); err != nil {
			response["emailWarning"] = "Invoice created but email sending failed: " + err.Error()
		} else {
			response["emailSent"] = true
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleGetDepositAmount handles POST /api/stripe/deposit/amount
// This is equivalent to STRIPE_GET_BOOKING_DEPOSIT_AMOUNT
func (h *StripeHandler) HandleGetDepositAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		EstimatedTotal *float64 `json:"estimatedTotal"`
		Estimate       *float64 `json:"estimate"`
		DepositValue   *float64 `json:"depositValue"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	var estimateCents *int64
	if body.EstimatedTotal != nil {
		cents := int64(*body.EstimatedTotal * 100)
		estimateCents = &cents
	} else if body.Estimate != nil {
		cents := int64(*body.Estimate * 100)
		estimateCents = &cents
	}

	if estimateCents == nil {
		util.WriteError(w, http.StatusBadRequest, "estimatedTotal or estimate is required")
		return
	}

	// Calculate deposit
	deposit, err := h.paymentsProvider.CalculateDeposit(r.Context(), *estimateCents)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to calculate deposit: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"ok": true,
		"deposit": map[string]interface{}{
			"amountCents":   deposit.AmountCents,
			"amountDollars": deposit.AmountDollars,
			"percentage":   deposit.Percentage,
		},
		"estimate": map[string]interface{}{
			"totalCents":   deposit.EstimateTotalCents,
			"totalDollars": float64(deposit.EstimateTotalCents) / 100,
		},
	}

	if body.DepositValue != nil {
		response["manualOverride"] = true
		response["requestedDeposit"] = *body.DepositValue
	}

	util.WriteJSON(w, http.StatusOK, response)
}

