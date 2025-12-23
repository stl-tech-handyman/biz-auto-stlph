package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/infra/stripe"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
)

// EstimateHandler handles estimate-related endpoints
type EstimateHandler struct {
	paymentsProvider ports.PaymentsProvider
}

// NewEstimateHandler creates a new estimate handler
func NewEstimateHandler(paymentsProvider ports.PaymentsProvider) *EstimateHandler {
	return &EstimateHandler{
		paymentsProvider: paymentsProvider,
	}
}

// HandleCalculate handles POST /api/estimate
func (h *EstimateHandler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		EventDate      string  `json:"eventDate"`
		DurationHours  float64 `json:"durationHours"`
		NumHelpers     int     `json:"numHelpers"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.EventDate == "" {
		util.WriteError(w, http.StatusBadRequest, "eventDate is required")
		return
	}

	if body.DurationHours <= 0 {
		util.WriteError(w, http.StatusBadRequest, "durationHours must be a positive number")
		return
	}

	if body.NumHelpers <= 0 {
		util.WriteError(w, http.StatusBadRequest, "numHelpers must be a positive integer")
		return
	}

	eventDate, err := time.Parse("2006-01-02", body.EventDate)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid eventDate format: "+err.Error())
		return
	}

	result, err := pricing.CalculateEstimate(eventDate, body.DurationHours, body.NumHelpers)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Calculate deposit with full details
	estimateCents := int64(result.TotalCost * 100)
	deposit, _ := h.paymentsProvider.CalculateDeposit(r.Context(), estimateCents)
	
	// Build full deposit structure matching JavaScript API format
	depositSections := buildDepositSections(estimateCents, deposit)

	response := map[string]interface{}{
		"ok": true,
		"data": map[string]interface{}{
			"year": result.Year,
			"eventDate": result.EventDate,
			"dateKey": result.DateKey,
			"numHelpers": result.NumHelpers,
			"durationHours": result.DurationHours,
			"basePerHelper": result.BasePerHelper,
			"extraPerHourPerHelper": result.ExtraPerHourPerHelper,
			"baseSubtotal": result.BaseSubtotal,
			"extraSubtotal": result.ExtraSubtotal,
			"subtotalBeforeAdjustments": result.SubtotalBeforeAdjustments,
			"isSpecialDate": result.IsSpecialDate,
			"specialLabel": result.SpecialLabel,
			"rateType": result.RateType,
			"specialDateMultiplier": result.SpecialDateMultiplier,
			"specialDateFlatIncrease": result.SpecialDateFlatIncrease,
			"totalCost": result.TotalCost,
			"currency": result.Currency,
			"breakdown": result.Breakdown,
			"calculationSummary": result.CalculationSummary,
			"deposit": depositSections,
		},
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// buildDepositSections builds the full deposit structure matching JavaScript API format
func buildDepositSections(estimateCents int64, deposit *domain.Deposit) map[string]interface{} {
	const (
		minPercent = 0.25 // 25%
		maxPercent = 0.40 // 40%
	)
	
	// Calculate deposit details using the same logic as JavaScript
	calc := stripe.CalculateDepositFromEstimate(estimateCents)
	
	// Round amounts
	roundedMin := int64(math.Round(float64(calc.MinAmount)))
	roundedMax := int64(math.Round(float64(calc.MaxAmount)))
	roundedTarget := int64(math.Round(float64(calc.TargetAmount)))
	roundedFlooredTarget := int64(math.Round(float64(calc.FlooredAmount)))
	
	centsToDollars := func(cents int64) float64 {
		return math.Round(float64(cents)) / 100
	}
	
	return map[string]interface{}{
		"recommended": map[string]interface{}{
			"amountCents": deposit.AmountCents,
			"amount":      deposit.AmountDollars,
			"percentage":  deposit.Percentage,
			"pickedBy":    calc.PickedBy,
			"isManualOverride": false,
			"estimateSource": "provided",
		},
		"range": map[string]interface{}{
			"minPercent":    minPercent * 100,
			"maxPercent":    maxPercent * 100,
			"minAmountCents": roundedMin,
			"maxAmountCents": roundedMax,
			"minAmount":     centsToDollars(roundedMin),
			"maxAmount":     centsToDollars(roundedMax),
			"description":   "Target booking deposits stay between 25% and 40% of the estimate.",
		},
		"calculation": map[string]interface{}{
			"estimateCents":       estimateCents,
			"targetAmountCents":   roundedTarget,
			"targetAmount":        centsToDollars(roundedTarget),
			"flooredTargetCents": roundedFlooredTarget,
			"flooredTargetAmount": centsToDollars(roundedFlooredTarget),
			"summary": fmt.Sprintf("Recommended deposit is $%.2f (~%.1f%% of the total), rounded up to the nearest professional amount.",
				deposit.AmountDollars, deposit.Percentage),
		},
	}
}

// HandleSpecialDates handles GET /api/estimate/special-dates
func (h *EstimateHandler) HandleSpecialDates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	yearsAhead := 5
	if yearsStr := r.URL.Query().Get("years"); yearsStr != "" {
		if y, err := strconv.Atoi(yearsStr); err == nil && y >= 1 && y <= 20 {
			yearsAhead = y
		}
	}

	var startYear *int
	if startYearStr := r.URL.Query().Get("startYear"); startYearStr != "" {
		if y, err := strconv.Atoi(startYearStr); err == nil && y >= 2020 && y <= 2100 {
			startYear = &y
		}
	}

	result := pricing.GetAllSpecialDates(yearsAhead, startYear)

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok": true,
		"yearsAhead": yearsAhead,
		"startYear": func() int {
			if startYear != nil {
				return *startYear
			}
			return time.Now().Year()
		}(),
		"data": result,
	})
}

