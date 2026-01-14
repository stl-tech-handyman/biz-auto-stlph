package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/infra/stripe"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
)

// TestScenario represents a single test scenario
type TestScenario struct {
	Name            string   `json:"name"`
	EventDate       string   `json:"eventDate"` // YYYY-MM-DD format
	EventTime       string   `json:"eventTime"` // HH:MM format
	Helpers         int      `json:"helpers"`
	Hours           float64  `json:"hours"`
	IsHighDemand    bool     `json:"isHighDemand"`
	ExpectedUrgency string   `json:"expectedUrgency"` // critical, urgent, high, moderate, normal
	ExpectedDays    int      `json:"expectedDays"`
	ExpectedBanner  string   `json:"expectedBanner"` // Expected banner text (if any)
	ExpectedHTML    []string `json:"expectedHTML"`   // HTML elements that must be present
	MustNotContain  []string `json:"mustNotContain"` // HTML elements that must NOT be present
}

// TestResult represents the result of a single test scenario
type TestResult struct {
	Scenario      TestScenario `json:"scenario"`
	Passed        bool         `json:"passed"`
	HTMLGenerated string       `json:"htmlGenerated"`
	Errors        []string     `json:"errors"`
	Warnings      []string     `json:"warnings"`
}

// HandleQuoteEmailTestAll handles GET /api/email/quote/test-all - runs all test scenarios
func (h *EmailHandler) HandleQuoteEmailTestAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Define all test scenarios
	now := time.Now()
	scenarios := []TestScenario{
		// Critical urgency (≤3 days)
		{
			Name:            "Critical: 1 day ahead",
			EventDate:       now.Add(1 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "18:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "critical",
			ExpectedDays:    1,
			ExpectedBanner:  "Only a few spots left — secure your date today",
			ExpectedHTML: []string{
				"Only a few spots left",
				"secure your date today",
				"#b91c1c", // Dark red background
			},
		},
		{
			Name:            "Critical: 3 days ahead",
			EventDate:       now.Add(3 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "19:00",
			Helpers:         3,
			Hours:           5.0,
			IsHighDemand:    false,
			ExpectedUrgency: "critical",
			ExpectedDays:    3,
			ExpectedBanner:  "Only a few spots left",
			ExpectedHTML: []string{
				"Only a few spots left",
				"#b91c1c",
			},
		},
		// Urgent urgency (4-7 days)
		{
			Name:            "Urgent: 4 days ahead",
			EventDate:       now.Add(4 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "17:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "urgent",
			ExpectedDays:    4,
			ExpectedBanner:  "Dates fill up fast",
			ExpectedHTML: []string{
				"Dates fill up fast",
				"secure your spot now",
				"#b91c1c",
			},
		},
		{
			Name:            "Urgent: 7 days ahead",
			EventDate:       now.Add(7 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "20:00",
			Helpers:         4,
			Hours:           6.0,
			IsHighDemand:    false,
			ExpectedUrgency: "urgent",
			ExpectedDays:    7,
			ExpectedBanner:  "Dates fill up fast",
			ExpectedHTML: []string{
				"Dates fill up fast",
				"#b91c1c",
			},
		},
		// High urgency (8-14 days)
		{
			Name:            "High: 8 days ahead",
			EventDate:       now.Add(8 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "18:30",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "high",
			ExpectedDays:    8,
			ExpectedBanner:  "Popular time period",
			ExpectedHTML: []string{
				"Popular time period",
				"secure your date soon",
				"#dc6a1c", // Orange warning
			},
		},
		{
			Name:            "High: 14 days ahead",
			EventDate:       now.Add(14 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "19:00",
			Helpers:         3,
			Hours:           5.0,
			IsHighDemand:    false,
			ExpectedUrgency: "high",
			ExpectedDays:    14,
			ExpectedBanner:  "Popular time period",
			ExpectedHTML: []string{
				"Popular time period",
				"#dc6a1c",
			},
		},
		// Moderate urgency (15-30 days)
		{
			Name:            "Moderate: 15 days ahead",
			EventDate:       now.Add(15 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "17:30",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "moderate",
			ExpectedDays:    15,
			ExpectedBanner:  "Spots are filling up",
			ExpectedHTML: []string{
				"Spots are filling up",
				"#f59e0b", // Less intense orange
			},
		},
		{
			Name:            "Moderate: 30 days ahead",
			EventDate:       now.Add(30 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "18:00",
			Helpers:         3,
			Hours:           4.5,
			IsHighDemand:    false,
			ExpectedUrgency: "moderate",
			ExpectedDays:    30,
			ExpectedBanner:  "Spots are filling up",
			ExpectedHTML: []string{
				"Spots are filling up",
				"#f59e0b",
			},
		},
		// Normal urgency (>30 days)
		{
			Name:            "Normal: 31 days ahead",
			EventDate:       now.Add(31 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "19:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "normal",
			ExpectedDays:    31,
			ExpectedBanner:  "", // No urgency banner for normal
			ExpectedHTML: []string{
				"Birthday Party Quote",
				"Event Details",
			},
			MustNotContain: []string{
				"Only a few spots left",
				"Dates fill up fast",
				"Popular time period",
				"Spots are filling up",
			},
		},
		{
			Name:            "Normal: 60 days ahead",
			EventDate:       now.Add(60 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "20:00",
			Helpers:         4,
			Hours:           6.0,
			IsHighDemand:    false,
			ExpectedUrgency: "normal",
			ExpectedDays:    60,
			ExpectedBanner:  "",
			ExpectedHTML: []string{
				"Birthday Party Quote",
			},
			MustNotContain: []string{
				"Only a few spots left",
				"Dates fill up fast",
			},
		},
		// High demand dates (special dates)
		{
			Name:            "High Demand: New Year's Eve (normal urgency)",
			EventDate:       "2026-12-31",
			EventTime:       "20:00",
			Helpers:         3,
			Hours:           5.0,
			IsHighDemand:    true,
			ExpectedUrgency: "normal", // >30 days from now
			ExpectedDays:    int(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Sub(now).Hours() / 24),
			ExpectedBanner:  "Popular Date",
			ExpectedHTML: []string{
				"Popular Date",
				"Dates Fill Up Fast",
			},
		},
		// Edge cases
		{
			Name:            "Edge: Exactly 3 days (critical)",
			EventDate:       now.Add(3 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "18:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "critical",
			ExpectedDays:    3,
			ExpectedBanner:  "Only a few spots left",
			ExpectedHTML: []string{
				"Only a few spots left",
				"#b91c1c",
			},
		},
		{
			Name:            "Edge: Exactly 7 days (urgent)",
			EventDate:       now.Add(7 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "19:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "urgent",
			ExpectedDays:    7,
			ExpectedBanner:  "Dates fill up fast",
			ExpectedHTML: []string{
				"Dates fill up fast",
				"#b91c1c",
			},
		},
		{
			Name:            "Edge: Exactly 14 days (high)",
			EventDate:       now.Add(14 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "18:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "high",
			ExpectedDays:    14,
			ExpectedBanner:  "Popular time period",
			ExpectedHTML: []string{
				"Popular time period",
				"#dc6a1c",
			},
		},
		{
			Name:            "Edge: Exactly 30 days (moderate)",
			EventDate:       now.Add(30 * 24 * time.Hour).Format("2006-01-02"),
			EventTime:       "19:00",
			Helpers:         2,
			Hours:           4.0,
			IsHighDemand:    false,
			ExpectedUrgency: "moderate",
			ExpectedDays:    30,
			ExpectedBanner:  "Spots are filling up",
			ExpectedHTML: []string{
				"Spots are filling up",
				"#f59e0b",
			},
		},
	}

	// Run all test scenarios
	results := []TestResult{}
	for _, scenario := range scenarios {
		result := h.runTestScenario(r.Context(), scenario)
		results = append(results, result)
	}

	// Calculate summary
	passedCount := 0
	failedCount := 0
	for _, result := range results {
		if result.Passed {
			passedCount++
		} else {
			failedCount++
		}
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ok": true,
		"summary": map[string]interface{}{
			"total":  len(results),
			"passed": passedCount,
			"failed": failedCount,
		},
		"results": results,
	})
}

// runTestScenario runs a single test scenario and returns the result
func (h *EmailHandler) runTestScenario(ctx context.Context, scenario TestScenario) TestResult {
	result := TestResult{
		Scenario: scenario,
		Passed:   true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Parse event date
	parsedEventDate, err := parseEventDateFromFormatted(scenario.EventDate)
	if err != nil {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to parse event date: %v", err))
		return result
	}

	// Calculate estimate
	estimate, err := pricing.CalculateEstimate(parsedEventDate, scenario.Hours, scenario.Helpers)
	if err != nil {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to calculate estimate: %v", err))
		return result
	}

	// Calculate deposit
	estimateCents := util.DollarsToCents(estimate.TotalCost)
	depositCalc := stripe.CalculateDepositFromEstimate(estimateCents)
	depositAmount := util.CentsToDollars(depositCalc.Value)

	// Calculate days until event and urgency
	now := time.Now()
	daysUntilEvent := int(parsedEventDate.Sub(now).Hours() / 24)
	if daysUntilEvent < 0 {
		daysUntilEvent = 0
	}
	// Calculate urgency level manually (same logic as util.CalculateUrgencyLevel)
	var urgencyLevel string
	if daysUntilEvent <= 3 {
		urgencyLevel = "critical"
	} else if daysUntilEvent <= 7 {
		urgencyLevel = "urgent"
	} else if daysUntilEvent <= 14 {
		urgencyLevel = "high"
	} else if daysUntilEvent <= 30 {
		urgencyLevel = "moderate"
	} else {
		urgencyLevel = "normal"
	}

	// Validate urgency level
	if urgencyLevel != scenario.ExpectedUrgency {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("Expected urgency '%s', got '%s'", scenario.ExpectedUrgency, urgencyLevel))
	}

	// Validate days until event (allow ±1 day tolerance for edge cases)
	if scenario.ExpectedDays > 0 {
		daysDiff := daysUntilEvent - scenario.ExpectedDays
		if daysDiff < -1 || daysDiff > 1 {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("Expected %d days until event, got %d", scenario.ExpectedDays, daysUntilEvent))
		}
	}

	// Determine rate label
	rateLabel := "Base Rate"
	if estimate.SpecialLabel != nil {
		rateLabel = *estimate.SpecialLabel
	}

	// Format event date and time
	// Format as "Fri, Jan 19, 2026" (day of week, short month, day, year)
	eventDateFormatted := parsedEventDate.Format("Mon, Jan 2, 2006")
	eventTimeFormatted := formatTimeFromHHMM(scenario.EventTime)

	// Calculate expiration date based on urgency
	var expirationDate time.Time
	var expirationFormatted string
	switch urgencyLevel {
	case "critical":
		location, _ := time.LoadLocation("America/Chicago")
		today := time.Now().In(location)
		expirationDate = time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, location)
		expirationFormatted = expirationDate.Format("January 2, 2006 at 11:59 PM")
	case "urgent":
		expirationDate = now.Add(48 * time.Hour)
		expirationFormatted = expirationDate.Format("January 2, 2006 at 3:04 PM")
	default:
		expirationDate = now.Add(72 * time.Hour)
		expirationFormatted = expirationDate.Format("January 2, 2006 at 3:04 PM")
	}

	// Generate email HTML
	emailData := util.QuoteEmailData{
		ClientName:         "Test Client",
		EventDate:          eventDateFormatted,
		EventTime:          eventTimeFormatted,
		EventLocation:      "123 Test St, St. Louis, MO 63110",
		Occasion:           "Birthday Party",
		GuestCount:         50,
		Helpers:            scenario.Helpers,
		Hours:              scenario.Hours,
		BaseRate:           estimate.BasePerHelper,
		HourlyRate:         estimate.ExtraPerHourPerHelper,
		TotalCost:          estimate.TotalCost,
		DepositAmount:      depositAmount,
		RateLabel:          rateLabel,
		ExpirationDate:     expirationFormatted,
		DepositLink:        "https://invoice.stripe.com/i/test",
		ConfirmationNumber: "TEST",
		IsHighDemand:       scenario.IsHighDemand || estimate.IsSpecialDate,
		DaysUntilEvent:     daysUntilEvent,
		UrgencyLevel:       urgencyLevel,
	}

	htmlBody := util.GenerateQuoteEmailHTML(emailData)
	result.HTMLGenerated = htmlBody

	// Validate HTML contains expected elements
	for _, expected := range scenario.ExpectedHTML {
		if !strings.Contains(htmlBody, expected) {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("Expected HTML to contain '%s'", expected))
		}
	}

	// Validate HTML does NOT contain forbidden elements
	for _, forbidden := range scenario.MustNotContain {
		if strings.Contains(htmlBody, forbidden) {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("HTML should NOT contain '%s'", forbidden))
		}
	}

	// Validate banner if expected
	if scenario.ExpectedBanner != "" {
		if !strings.Contains(htmlBody, scenario.ExpectedBanner) {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("Expected banner text '%s' not found", scenario.ExpectedBanner))
		}
	} else {
		// If no banner expected, check that urgency banners are not present (unless high demand)
		if !scenario.IsHighDemand {
			urgencyBanners := []string{
				"Only a few spots left",
				"Dates fill up fast",
				"Popular time period",
				"Spots are filling up",
			}
			for _, banner := range urgencyBanners {
				if strings.Contains(htmlBody, banner) {
					result.Warnings = append(result.Warnings, fmt.Sprintf("Unexpected urgency banner found: '%s'", banner))
				}
			}
		}
	}

	// Validate deposit amount ends in .00 or .50
	depositCents := int(depositAmount * 100)
	if depositCents%50 != 0 {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("Deposit amount $%.2f does not end in .00 or .50", depositAmount))
	}

	return result
}

// formatTimeFromHHMM formats time from "HH:MM" to "H:MM PM" format
func formatTimeFromHHMM(timeStr string) string {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return "6:00 PM" // Default
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return "6:00 PM"
	}

	minutes := parts[1]
	ampm := "AM"
	if hours >= 12 {
		ampm = "PM"
		if hours > 12 {
			hours -= 12
		}
	}
	if hours == 0 {
		hours = 12
	}

	return fmt.Sprintf("%d:%s %s", hours, minutes, ampm)
}
