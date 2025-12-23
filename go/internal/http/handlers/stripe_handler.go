package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bizops360/go-api/internal/http/handlers/dto"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/services/email"
	"github.com/bizops360/go-api/internal/services/pricing"
	stripeService "github.com/bizops360/go-api/internal/services/stripe"
	"github.com/bizops360/go-api/internal/util"
)

// StripeHandler handles Stripe-related endpoints
type StripeHandler struct {
	invoiceService  *stripeService.InvoiceService
	emailHandler    *EmailHandler
	templateService *email.TemplateService
}

// NewStripeHandler creates a new Stripe handler
func NewStripeHandler(paymentsProvider ports.PaymentsProvider) *StripeHandler {
	return &StripeHandler{
		invoiceService:  stripeService.NewInvoiceService(paymentsProvider),
		templateService: email.NewTemplateService(),
	}
}

// SetEmailHandler sets the email handler for orchestrated endpoints
func (h *StripeHandler) SetEmailHandler(emailHandler *EmailHandler) {
	h.emailHandler = emailHandler
}

// HandleDeposit handles POST /api/stripe/deposit
func (h *StripeHandler) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodPost, w) {
		return
	}

	var req dto.DepositRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"ok":      true,
		"message": "Deposit calculation (stub - full Stripe integration pending)",
	}

	// Handle manual deposit value
	if req.DepositValue != nil {
		depositCents := util.DollarsToCents(*req.DepositValue)
		response["deposit"] = map[string]interface{}{
			"value":        depositCents,
			"valueDollars": util.CentsToDollars(depositCents),
		}
	} else if req.Deposit != nil {
		depositCents := util.DollarsToCents(*req.Deposit)
		response["deposit"] = map[string]interface{}{
			"value":        depositCents,
			"valueDollars": util.CentsToDollars(depositCents),
		}
	}

	// Handle calculated deposit from estimate
	if req.EstimatedTotal != nil {
		estimateCents := util.DollarsToCents(*req.EstimatedTotal)
		deposit, err := h.invoiceService.CalculateDepositFromEstimate(r.Context(), estimateCents)
		if err == nil {
			response["deposit"] = map[string]interface{}{
				"value":        deposit.AmountCents,
				"valueDollars": deposit.AmountDollars,
				"percentage":   deposit.Percentage,
			}
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleDepositCalculate handles GET /api/stripe/deposit/calculate
func (h *StripeHandler) HandleDepositCalculate(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodGet, w) {
		return
	}

	estimateStr := r.URL.Query().Get("estimate")
	depositStr := r.URL.Query().Get("deposit")

	var estimateCents *int64
	var depositCents *int64

	if depositStr != "" {
		depositDollars, err := strconv.ParseFloat(depositStr, 64)
		if err == nil {
			cents := util.DollarsToCents(depositDollars)
			depositCents = &cents
		}
	}

	if estimateStr != "" {
		estimateDollars, err := strconv.ParseFloat(estimateStr, 64)
		if err == nil {
			cents := util.DollarsToCents(estimateDollars)
			estimateCents = &cents
		}
	}

	if depositCents == nil && estimateCents == nil {
		util.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"ok":      true,
			"message": "No estimate or deposit provided.",
			"usage":   "Add ?estimate=100 to calculate deposit, or ?deposit=50 to set manually (both in dollars)",
		})
		return
	}

	response := map[string]interface{}{
		"ok": true,
	}

	if depositCents != nil {
		response["deposit"] = util.CentsToDollars(*depositCents)
		response["depositCents"] = *depositCents
		response["pickedBy"] = "manual"
		response["isManualOverride"] = true
	}

	if estimateCents != nil {
		calc, err := h.invoiceService.CalculateDepositFromEstimate(r.Context(), *estimateCents)
		if err == nil {
			if depositCents == nil {
				response["deposit"] = calc.AmountDollars
				response["depositCents"] = calc.AmountCents
				response["pickedBy"] = "calculated"
				response["isManualOverride"] = false
			}
			response["requested_estimate"] = util.CentsToDollars(*estimateCents)
			response["calculation"] = map[string]interface{}{
				"deposit":    calc.AmountDollars,
				"percentage": calc.Percentage,
				"min_range":  util.CentsToDollars(calc.AmountCents - 1000),
				"max_range":  util.CentsToDollars(calc.AmountCents + 1000),
			}
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleDepositWithEmail handles POST /api/stripe/deposit/with-email
func (h *StripeHandler) HandleDepositWithEmail(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodPost, w) {
		return
	}

	var req dto.DepositWithEmailRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if !ValidateRequiredString(req.Name, "name", w) {
		return
	}

	// Calculate estimate if event details provided
	var estimateResult *pricing.EstimateResult
	var depositCents int64

	if req.EventDateTimeLocal != "" || req.EventDate != "" {
		eventDateStr := req.EventDateTimeLocal
		if eventDateStr == "" {
			eventDateStr = req.EventDate
		}

		durationHours := 4.0
		if req.Hours != nil {
			durationHours = *req.Hours
		} else if req.Duration != nil {
			durationHours = *req.Duration
		}

		numHelpers := 2
		if req.HelpersCount != nil {
			numHelpers = *req.HelpersCount
		}

		deposit, result, err := h.invoiceService.CalculateDepositFromEventDetails(
			r.Context(), eventDateStr, durationHours, numHelpers)
		if err == nil {
			depositCents = deposit.AmountCents
			estimateResult = result
		}
	}

	// Determine deposit from various sources
	if req.DepositValue != nil {
		depositCents = util.DollarsToCents(*req.DepositValue)
	} else if estimateResult != nil {
		// Already calculated above
	} else if req.Estimate != nil {
		deposit, err := h.invoiceService.CalculateDepositFromEstimate(
			r.Context(), util.DollarsToCents(*req.Estimate))
		if err == nil {
			depositCents = deposit.AmountCents
		}
	} else if req.EstimatedTotal != nil {
		deposit, err := h.invoiceService.CalculateDepositFromEstimate(
			r.Context(), util.DollarsToCents(*req.EstimatedTotal))
		if err == nil {
			depositCents = deposit.AmountCents
		}
	}

	response := map[string]interface{}{
		"ok":          true,
		"message":     "Invoice generated (stub - full Stripe integration pending)",
		"dryRun":      req.DryRun,
		"saveAsDraft": req.SaveAsDraft,
		"generatedInvoice": map[string]interface{}{
			"id":     "stub_invoice_id",
			"url":    "https://stripe.com/invoice/stub",
			"amount": util.CentsToDollars(depositCents),
		},
	}

	if estimateResult != nil {
		response["basedOnEstimate"] = estimateResult
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleTest handles POST /api/stripe/test
func (h *StripeHandler) HandleTest(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodPost, w) {
		return
	}

	// Get query parameters
	useTest := r.URL.Query().Get("useTest") == "true" || r.URL.Query().Get("use_test") == "true"
	sendEmail := r.URL.Query().Get("sendEmail") == "true" || r.URL.Query().Get("send_email") == "true"

	var body map[string]interface{}
	if err := util.ReadJSON(r, &body); err != nil {
		body = make(map[string]interface{})
	}

	// Parse useTest and sendEmail from body if present
	if val, ok := body["useTest"].(bool); ok {
		useTest = val
	}
	if val, ok := body["sendEmail"].(bool); ok {
		sendEmail = val
	}

	// Parse estimated total
	estimatedTotal := 10.0
	if val, err := ParseFloatFromMap(body, "estimatedTotal"); err == nil && val != nil {
		estimatedTotal = *val
	} else if val, err := ParseFloatFromQuery(r, "estimate"); err == nil && val != nil {
		estimatedTotal = *val
	}

	// Convert to cents if needed
	estimatedTotalCents := util.DollarsToCents(estimatedTotal)
	if estimatedTotal < 10000 {
		estimatedTotalCents = util.DollarsToCents(estimatedTotal)
	}

	// Parse deposit value
	var depositValueCents *int64
	if val, err := ParseFloatFromMap(body, "depositValue"); err == nil && val != nil {
		cents := util.DollarsToCents(*val)
		if *val < 10000 {
			cents = util.DollarsToCents(*val)
		}
		depositValueCents = &cents
	} else if val, err := ParseFloatFromQuery(r, "depositValue"); err == nil && val != nil {
		cents := util.DollarsToCents(*val)
		if *val < 10000 {
			cents = util.DollarsToCents(*val)
		}
		depositValueCents = &cents
	}

	// Get customer info
	customerEmail := GetStringFromMap(body, "email")
	if customerEmail == "" {
		customerEmail = r.URL.Query().Get("email")
		if customerEmail == "" {
			customerEmail = "bizops-dev-alexey-at-shevelyov-dot-com@shevelyov.com"
		}
	}

	customerName := GetStringFromMap(body, "name")
	if customerName == "" {
		customerName = r.URL.Query().Get("name")
		if customerName == "" {
			customerName = "Test Customer"
		}
	}

	// Create invoice request
	amountCents := estimatedTotalCents
	if depositValueCents != nil {
		amountCents = *depositValueCents
	}

	invoiceReq := &ports.CreateInvoiceRequest{
		CustomerEmail: customerEmail,
		CustomerName:  customerName,
		AmountCents:   amountCents,
		Currency:      "usd",
		Description:   "Test Booking Deposit Invoice",
	}

	// Generate invoice
	invoiceResult, err := h.invoiceService.CreateDepositInvoice(r.Context(), &stripeService.CreateDepositInvoiceRequest{
		CustomerEmail:     invoiceReq.CustomerEmail,
		CustomerName:      invoiceReq.CustomerName,
		DepositValueCents: &amountCents,
		Description:       invoiceReq.Description,
	})
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
			"amount": util.CentsToDollars(invoiceResult.AmountDue),
		},
		"demoData": map[string]interface{}{
			"email":          customerEmail,
			"name":           customerName,
			"estimatedTotal": util.CentsToDollars(estimatedTotalCents),
			"depositValue": func() interface{} {
				if depositValueCents != nil {
					return util.CentsToDollars(*depositValueCents)
				}
				return nil
			}(),
			"useTest": useTest,
		},
	}

	if sendEmail {
		response["email"] = map[string]interface{}{
			"sent":    false,
			"message": "Email sending not yet implemented in Go API",
			"note":    "Use /api/stripe/deposit/with-email for full email integration",
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// createFinalInvoiceCommon contains common logic for creating final invoices
func (h *StripeHandler) createFinalInvoiceCommon(ctx context.Context, req dto.FinalInvoiceRequest) (*ports.InvoiceResult, error) {
	// Validate required fields
	if req.Email == "" || req.Name == "" {
		return nil, fmt.Errorf("email and name are required")
	}

	// Create final invoice using service
	// Convert DTO CustomFields to ports CustomFields (only non-empty fields)
	customFields := make([]ports.CustomField, 0, len(req.CustomFields))
	for _, cf := range req.CustomFields {
		// Only include non-empty fields (both name and value must be non-empty)
		if strings.TrimSpace(cf.Name) != "" && strings.TrimSpace(cf.Value) != "" {
			customFields = append(customFields, ports.CustomField{
				Name:  strings.TrimSpace(cf.Name),
				Value: strings.TrimSpace(cf.Value),
			})
		}
	}

	// Convert ports CustomFields to service layer CustomFields
	serviceCustomFields := make([]ports.CustomField, len(customFields))
	copy(serviceCustomFields, customFields)

	invoiceReq := &stripeService.CreateFinalInvoiceRequest{
		CustomerEmail:    req.Email,
		CustomerName:     req.Name,
		TotalAmountCents: req.TotalAmountCents,
		TotalAmount:      req.TotalAmount,
		DepositPaidCents: req.DepositPaidCents,
		DepositPaid:      req.DepositPaid,
		Currency:         req.Currency,
		Description:      req.Description,
		Metadata:         req.Metadata,
		CustomFields:     serviceCustomFields,
	}

	return h.invoiceService.CreateFinalInvoice(ctx, invoiceReq)
}

// HandleFinalInvoice handles POST /api/stripe/final-invoice
func (h *StripeHandler) HandleFinalInvoice(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodPost, w) {
		return
	}

	var req dto.FinalInvoiceRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if !ValidateRequiredString(req.Email, "email", w) || !ValidateRequiredString(req.Name, "name", w) {
		return
	}

	// Create final invoice
	invoiceResult, err := h.createFinalInvoiceCommon(r.Context(), req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Failed to create final invoice: "+err.Error())
		return
	}

	// Calculate amounts for response
	totalCents := int64(0)
	if req.TotalAmountCents != nil {
		totalCents = *req.TotalAmountCents
	} else if req.TotalAmount != nil {
		totalCents = util.DollarsToCents(*req.TotalAmount)
	}

	depositPaidCents := int64(0)
	if req.DepositPaidCents != nil {
		depositPaidCents = *req.DepositPaidCents
	} else if req.DepositPaid != nil {
		depositPaidCents = util.DollarsToCents(*req.DepositPaid)
	}

	response := map[string]interface{}{
		"ok":      true,
		"message": "Final invoice created successfully",
		"invoice": map[string]interface{}{
			"id":     invoiceResult.InvoiceID,
			"url":    invoiceResult.HostedInvoiceURL,
			"amount": util.CentsToDollars(invoiceResult.AmountDue),
			"status": invoiceResult.Status,
			"pdf":    invoiceResult.InvoicePDF,
		},
		"details": map[string]interface{}{
			"totalAmount":      util.CentsToDollars(totalCents),
			"depositPaid":      util.CentsToDollars(depositPaidCents),
			"remainingBalance": util.CentsToDollars(invoiceResult.AmountDue),
		},
	}

	// Optionally send invoice via Stripe
	if req.SendEmail {
		if err := h.invoiceService.SendInvoice(r.Context(), invoiceResult.InvoiceID, req.UseTest); err != nil {
			response["emailWarning"] = "Invoice created but email sending failed: " + err.Error()
		} else {
			response["emailSent"] = true
		}
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleFinalInvoiceWithEmail handles POST /api/stripe/final-invoice/with-email
func (h *StripeHandler) HandleFinalInvoiceWithEmail(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodPost, w) {
		return
	}

	var req dto.FinalInvoiceRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if !ValidateRequiredString(req.Email, "email", w) || !ValidateRequiredString(req.Name, "name", w) {
		return
	}

	// Create final invoice
	invoiceResult, err := h.createFinalInvoiceCommon(r.Context(), req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Failed to create final invoice: "+err.Error())
		return
	}

	// Calculate amounts for email
	totalCents := int64(0)
	if req.TotalAmountCents != nil {
		totalCents = *req.TotalAmountCents
	} else if req.TotalAmount != nil {
		totalCents = util.DollarsToCents(*req.TotalAmount)
	}

	depositPaidCents := int64(0)
	if req.DepositPaidCents != nil {
		depositPaidCents = *req.DepositPaidCents
	} else if req.DepositPaid != nil {
		depositPaidCents = util.DollarsToCents(*req.DepositPaid)
	}

	remainingBalance := util.CentsToDollars(invoiceResult.AmountDue)
	totalAmount := util.CentsToDollars(totalCents)
	depositPaid := util.CentsToDollars(depositPaidCents)

	// Send custom email
	var emailSent bool
	var emailError string
	if h.emailHandler != nil {
		emailSent, emailError = h.emailHandler.SendFinalInvoiceEmail(
			r.Context(),
			req.Name,
			req.Email,
			totalAmount,
			depositPaid,
			remainingBalance,
			invoiceResult.HostedInvoiceURL,
		)
	} else {
		emailError = "email handler is not configured"
	}

	response := map[string]interface{}{
		"ok":      true,
		"message": "Final invoice created and email sent",
		"invoice": map[string]interface{}{
			"id":     invoiceResult.InvoiceID,
			"url":    invoiceResult.HostedInvoiceURL,
			"amount": remainingBalance,
			"status": invoiceResult.Status,
			"pdf":    invoiceResult.InvoicePDF,
		},
		"details": map[string]interface{}{
			"totalAmount":      totalAmount,
			"depositPaid":      depositPaid,
			"remainingBalance": remainingBalance,
		},
		"email": map[string]interface{}{
			"sent":  emailSent,
			"error": emailError,
		},
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// HandleGetDepositAmount handles POST /api/stripe/deposit/amount
func (h *StripeHandler) HandleGetDepositAmount(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(r, http.MethodPost, w) {
		return
	}

	var req dto.DepositAmountRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	var estimateCents *int64
	if req.EstimatedTotal != nil {
		cents := util.DollarsToCents(*req.EstimatedTotal)
		estimateCents = &cents
	} else if req.Estimate != nil {
		cents := util.DollarsToCents(*req.Estimate)
		estimateCents = &cents
	}

	if estimateCents == nil {
		util.WriteError(w, http.StatusBadRequest, "estimatedTotal or estimate is required")
		return
	}

	// Calculate deposit
	deposit, err := h.invoiceService.CalculateDepositFromEstimate(r.Context(), *estimateCents)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to calculate deposit: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"ok": true,
		"deposit": map[string]interface{}{
			"amountCents":   deposit.AmountCents,
			"amountDollars": deposit.AmountDollars,
			"percentage":    deposit.Percentage,
		},
		"estimate": map[string]interface{}{
			"totalCents":   deposit.EstimateTotalCents,
			"totalDollars": util.CentsToDollars(deposit.EstimateTotalCents),
		},
	}

	if req.DepositValue != nil {
		response["manualOverride"] = true
		response["requestedDeposit"] = *req.DepositValue
	}

	util.WriteJSON(w, http.StatusOK, response)
}
