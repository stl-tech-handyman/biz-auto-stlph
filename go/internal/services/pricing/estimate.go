package pricing

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// BaseRate holds base and extra rates for a year
type BaseRate struct {
	BasePerHelper        float64
	ExtraPerHourPerHelper float64
}

// Base rates by year
var baseRateByYear = map[int]BaseRate{
	2025: {BasePerHelper: 275, ExtraPerHourPerHelper: 45},
	2026: {BasePerHelper: 300, ExtraPerHourPerHelper: 50},
	2027: {BasePerHelper: 325, ExtraPerHourPerHelper: 55},
	2028: {BasePerHelper: 400, ExtraPerHourPerHelper: 60},
	2029: {BasePerHelper: 475, ExtraPerHourPerHelper: 65},
	2030: {BasePerHelper: 550, ExtraPerHourPerHelper: 70},
}

// SpecialDateRule represents a special date rule
type SpecialDateRule struct {
	Multiplier   *float64
	FlatIncrease *float64
	Label        string
	Type         string // "holiday" | "surge"
}

// Legacy special date rules (for backward compatibility)
var legacySpecialDateRules = map[string]SpecialDateRule{
	// 2025
	"2025-01-01": {Multiplier: floatPtr(2), Label: "New Year's Day"},
	"2025-11-27": {Multiplier: floatPtr(2), Label: "Thanksgiving"},
	"2025-12-15": {Multiplier: floatPtr(2), Label: "Special Date"},
	"2025-12-23": {Multiplier: floatPtr(2), Label: "Special Date"},
	"2025-12-24": {Multiplier: floatPtr(2), Label: "Christmas Eve"},
	"2025-12-25": {Multiplier: floatPtr(2), Label: "Christmas Day"},
	"2025-12-30": {Multiplier: floatPtr(2), Label: "Special Date"},
	"2025-12-31": {Multiplier: floatPtr(2), Label: "New Year's Eve"},
	// Add more years as needed...
}

// Surge date rules
var surgeDateRules = map[string]SpecialDateRule{
	// Recurring surge dates - January 1st
	"2025-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2026-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2027-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2028-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2029-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2030-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2031-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	"2032-01-01": {Multiplier: floatPtr(1.5), Label: "New Year Surge", Type: "surge"},
	
	// Recurring surge dates - February 5th
	"2025-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2026-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2027-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2028-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2029-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2030-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2031-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	"2032-02-05": {Multiplier: floatPtr(1.5), Label: "February Surge", Type: "surge"},
	
	// Recurring surge dates - May 17th
	"2025-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2026-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2027-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2028-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2029-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2030-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2031-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	"2032-05-17": {Multiplier: floatPtr(1.5), Label: "May Surge", Type: "surge"},
	
	// Recurring surge dates - August 19th
	"2025-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2026-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2027-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2028-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2029-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2030-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2031-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	"2032-08-19": {Multiplier: floatPtr(1.5), Label: "August Surge", Type: "surge"},
	
	// Recurring surge dates - December 27th
	"2025-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2026-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2027-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2028-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2029-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2030-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2031-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
	"2032-12-27": {Multiplier: floatPtr(1.5), Label: "December Surge", Type: "surge"},
}

func floatPtr(f float64) *float64 {
	return &f
}

// GetBaseRatesForYear returns base rates for a given year
func GetBaseRatesForYear(year int) BaseRate {
	if rate, ok := baseRateByYear[year]; ok {
		return rate
	}
	if year < 2025 {
		return baseRateByYear[2025]
	}
	return baseRateByYear[2030]
}

// GetThanksgivingDay calculates Thanksgiving day (4th Thursday of November) for a year
func GetThanksgivingDay(year int) int {
	nov1 := time.Date(year, time.November, 1, 0, 0, 0, 0, time.UTC)
	dayOfWeek := int(nov1.Weekday())
	daysToAdd := (4-dayOfWeek+7)%7 + 21
	return 1 + daysToAdd
}

// GetHolidayDatesForYear returns holiday dates for a year
func GetHolidayDatesForYear(year int) map[string]SpecialDateRule {
	holidays := make(map[string]SpecialDateRule)
	thanksgivingDay := GetThanksgivingDay(year)
	
	holidays[fmt.Sprintf("%d-01-01", year)] = SpecialDateRule{Multiplier: floatPtr(2), Label: "New Year's Day", Type: "holiday"}
	holidays[fmt.Sprintf("%d-11-%02d", year, thanksgivingDay)] = SpecialDateRule{Multiplier: floatPtr(2), Label: "Thanksgiving", Type: "holiday"}
	holidays[fmt.Sprintf("%d-12-24", year)] = SpecialDateRule{Multiplier: floatPtr(2), Label: "Christmas Eve", Type: "holiday"}
	holidays[fmt.Sprintf("%d-12-25", year)] = SpecialDateRule{Multiplier: floatPtr(2), Label: "Christmas Day", Type: "holiday"}
	holidays[fmt.Sprintf("%d-12-31", year)] = SpecialDateRule{Multiplier: floatPtr(2), Label: "New Year's Eve", Type: "holiday"}
	
	return holidays
}

// ValidateSurgeMultiplier validates surge multiplier (1.25-3.0)
func ValidateSurgeMultiplier(multiplier float64) bool {
	return multiplier >= 1.25 && multiplier <= 3.0
}

// ToDateKey normalizes a date to "YYYY-MM-DD"
func ToDateKey(eventDate time.Time) string {
	return eventDate.Format("2006-01-02")
}

// EstimateResult represents the result of an estimate calculation
type EstimateResult struct {
	Year                      int     `json:"year"`
	EventDate                 string  `json:"eventDate"`
	DateKey                   string  `json:"dateKey"`
	NumHelpers               int     `json:"numHelpers"`
	DurationHours            float64 `json:"durationHours"`
	BasePerHelper            float64 `json:"basePerHelper"`
	ExtraPerHourPerHelper    float64 `json:"extraPerHourPerHelper"`
	BaseSubtotal             float64 `json:"baseSubtotal"`
	ExtraSubtotal            float64 `json:"extraSubtotal"`
	SubtotalBeforeAdjustments float64 `json:"subtotalBeforeAdjustments"`
	IsSpecialDate            bool    `json:"isSpecialDate"`
	SpecialLabel             *string `json:"specialLabel,omitempty"`
	RateType                 *string `json:"rateType,omitempty"`
	SpecialDateMultiplier    *float64 `json:"specialDateMultiplier,omitempty"`
	SpecialDateFlatIncrease  *float64 `json:"specialDateFlatIncrease,omitempty"`
	TotalCost                float64 `json:"totalCost"`
	Currency                 string  `json:"currency"`
	Breakdown                map[string]interface{} `json:"breakdown"`
	CalculationSummary       string  `json:"calculationSummary"`
}

// CalculateEstimate calculates event estimate
func CalculateEstimate(eventDate time.Time, durationHours float64, numHelpers int) (*EstimateResult, error) {
	if eventDate.IsZero() {
		return nil, fmt.Errorf("eventDate is required")
	}
	if durationHours <= 0 {
		return nil, fmt.Errorf("durationHours must be a positive number")
	}
	if numHelpers <= 0 {
		return nil, fmt.Errorf("numHelpers must be a positive integer")
	}

	year := eventDate.Year()
	rates := GetBaseRatesForYear(year)

	// Base block covers up to the first 4 hours
	billedBaseBlock := 1.0
	if durationHours <= 0 {
		billedBaseBlock = 0
	}
	extraHours := math.Max(durationHours-4, 0)

	baseSubtotal := rates.BasePerHelper * float64(numHelpers) * billedBaseBlock
	extraSubtotal := rates.ExtraPerHourPerHelper * float64(numHelpers) * extraHours
	subtotalBeforeAdjustments := baseSubtotal + extraSubtotal

	subtotal := subtotalBeforeAdjustments

	dateKey := ToDateKey(eventDate)

	// Check for holiday first
	holidayDates := GetHolidayDatesForYear(year)
	holidayRule, isHoliday := holidayDates[dateKey]

	// Check for surge date only if not a holiday
	var surgeRule SpecialDateRule
	var isSurge bool
	if !isHoliday {
		surgeRule, isSurge = surgeDateRules[dateKey]
	}

	// Check legacy rules
	var legacyRule SpecialDateRule
	var isLegacy bool
	if !isHoliday && !isSurge {
		legacyRule, isLegacy = legacySpecialDateRules[dateKey]
	}

	var specialRule SpecialDateRule
	var isSpecialDate bool
	var specialLabel *string
	var rateType *string

	if isHoliday {
		specialRule = holidayRule
		isSpecialDate = true
		specialLabel = &holidayRule.Label
		t := "holiday"
		rateType = &t
	} else if isSurge {
		specialRule = surgeRule
		isSpecialDate = true
		specialLabel = &surgeRule.Label
		t := "surge"
		rateType = &t
	} else if isLegacy {
		specialRule = legacyRule
		isSpecialDate = true
		specialLabel = &legacyRule.Label
	}

	if isSpecialDate {
		if specialRule.Multiplier != nil {
			if rateType != nil && *rateType == "surge" && !ValidateSurgeMultiplier(*specialRule.Multiplier) {
				return nil, fmt.Errorf("invalid surge multiplier: %.2f. Must be between 1.25 and 3.0", *specialRule.Multiplier)
			}
			subtotal *= *specialRule.Multiplier
		}
		if specialRule.FlatIncrease != nil {
			subtotal += *specialRule.FlatIncrease
		}
	}

	totalCost := math.Round(subtotal*100) / 100

	// Build breakdown
	breakdown := make(map[string]interface{})
	breakdown["baseBlock"] = fmt.Sprintf("%d helpers × $%.2f (first 4 hours) = $%.2f", numHelpers, rates.BasePerHelper, baseSubtotal)
	if extraHours > 0 {
		breakdown["extraHours"] = fmt.Sprintf("%d helpers × %.1f hours × $%.2f/hour = $%.2f", numHelpers, extraHours, rates.ExtraPerHourPerHelper, extraSubtotal)
	} else {
		breakdown["extraHours"] = nil
	}
	breakdown["subtotal"] = fmt.Sprintf("$%.2f", subtotalBeforeAdjustments)
	if isSpecialDate {
		adj := ""
		if specialRule.Multiplier != nil {
			adj += fmt.Sprintf("×%.2f", *specialRule.Multiplier)
		}
		if specialRule.FlatIncrease != nil {
			if adj != "" {
				adj += " "
			}
			adj += fmt.Sprintf("+ $%.2f", *specialRule.FlatIncrease)
		}
		if specialLabel != nil {
			breakdown["specialDateAdjustment"] = fmt.Sprintf("%s: %s", *specialLabel, adj)
		}
	} else {
		breakdown["specialDateAdjustment"] = nil
	}
	breakdown["total"] = fmt.Sprintf("$%.2f", totalCost)

	// Build calculation summary
	summary := fmt.Sprintf("%d helpers, %.1f hours, %d rates ($%.2f base + $%.2f/hour extra)", numHelpers, durationHours, year, rates.BasePerHelper, rates.ExtraPerHourPerHelper)
	if isSpecialDate && specialLabel != nil {
		adj := ""
		if rateType != nil {
			adj += *rateType + " "
		}
		if specialRule.Multiplier != nil {
			adj += fmt.Sprintf("×%.2f", *specialRule.Multiplier)
		}
		if specialRule.FlatIncrease != nil {
			if adj != "" {
				adj += " "
			}
			adj += fmt.Sprintf("+ $%.2f", *specialRule.FlatIncrease)
		}
		summary += fmt.Sprintf(", %s (%s) = $%.2f", *specialLabel, adj, totalCost)
	} else {
		summary += fmt.Sprintf(" = $%.2f", totalCost)
	}

	result := &EstimateResult{
		Year:                      year,
		EventDate:                 eventDate.Format(time.RFC3339),
		DateKey:                   dateKey,
		NumHelpers:               numHelpers,
		DurationHours:            durationHours,
		BasePerHelper:            rates.BasePerHelper,
		ExtraPerHourPerHelper:    rates.ExtraPerHourPerHelper,
		BaseSubtotal:             baseSubtotal,
		ExtraSubtotal:            extraSubtotal,
		SubtotalBeforeAdjustments: subtotalBeforeAdjustments,
		IsSpecialDate:            isSpecialDate,
		SpecialLabel:             specialLabel,
		RateType:                 rateType,
		SpecialDateMultiplier:    specialRule.Multiplier,
		SpecialDateFlatIncrease:  specialRule.FlatIncrease,
		TotalCost:                totalCost,
		Currency:                 "USD",
		Breakdown:                breakdown,
		CalculationSummary:       summary,
	}

	return result, nil
}

// SpecialDate represents a special date
type SpecialDate struct {
	Date       string   `json:"date"`
	Multiplier *float64 `json:"multiplier,omitempty"`
	FlatIncrease *float64 `json:"flatIncrease,omitempty"`
	Label      string   `json:"label"`
	Type       string   `json:"type"`
}

// YearSpecialDates represents special dates for a year
type YearSpecialDates struct {
	Holidays   []SpecialDate `json:"holidays"`
	SurgeDates []SpecialDate `json:"surgeDates"`
	LegacyDates []SpecialDate `json:"legacyDates"`
	AllDates   []SpecialDate `json:"allDates"`
}

// GetAllSpecialDates gets all special dates for the next N years
func GetAllSpecialDates(yearsAhead int, startYear *int) map[int]YearSpecialDates {
	currentYear := time.Now().Year()
	if startYear != nil {
		currentYear = *startYear
	}

	result := make(map[int]YearSpecialDates)

	for i := 0; i < yearsAhead; i++ {
		year := currentYear + i
		holidays := GetHolidayDatesForYear(year)

		// Convert holidays
		holidayList := make([]SpecialDate, 0)
		for dateKey, rule := range holidays {
			holidayList = append(holidayList, SpecialDate{
				Date:       dateKey,
				Multiplier: rule.Multiplier,
				FlatIncrease: rule.FlatIncrease,
				Label:      rule.Label,
				Type:       rule.Type,
			})
		}
		sort.Slice(holidayList, func(i, j int) bool {
			return holidayList[i].Date < holidayList[j].Date
		})

		// Get surge dates for this year
		surgeList := make([]SpecialDate, 0)
		for dateKey, rule := range surgeDateRules {
			if len(dateKey) >= 4 && dateKey[:4] == strconv.Itoa(year) {
				surgeList = append(surgeList, SpecialDate{
					Date:       dateKey,
					Multiplier: rule.Multiplier,
					FlatIncrease: rule.FlatIncrease,
					Label:      rule.Label,
					Type:       rule.Type,
				})
			}
		}
		sort.Slice(surgeList, func(i, j int) bool {
			return surgeList[i].Date < surgeList[j].Date
		})

		// Get legacy dates for this year
		legacyList := make([]SpecialDate, 0)
		for dateKey, rule := range legacySpecialDateRules {
			if len(dateKey) >= 4 && dateKey[:4] == strconv.Itoa(year) {
				// Only include if not already a holiday
				if _, isHoliday := holidays[dateKey]; !isHoliday {
					legacyList = append(legacyList, SpecialDate{
						Date:       dateKey,
						Multiplier: rule.Multiplier,
						FlatIncrease: rule.FlatIncrease,
						Label:      rule.Label,
						Type:       rule.Type,
					})
				}
			}
		}
		sort.Slice(legacyList, func(i, j int) bool {
			return legacyList[i].Date < legacyList[j].Date
		})

		// Combine all dates
		allDates := make([]SpecialDate, 0)
		allDates = append(allDates, holidayList...)
		allDates = append(allDates, surgeList...)
		allDates = append(allDates, legacyList...)
		sort.Slice(allDates, func(i, j int) bool {
			return allDates[i].Date < allDates[j].Date
		})

		result[year] = YearSpecialDates{
			Holidays:   holidayList,
			SurgeDates: surgeList,
			LegacyDates: legacyList,
			AllDates:   allDates,
		}
	}

	return result
}


// TravelFeeResult contains travel fee calculation details
type TravelFeeResult struct {
	IsWithinServiceArea bool    // True if within 15 miles
	DistanceMiles       float64 // Distance from office in miles
	TravelFee           float64 // Calculated travel fee (0 if within service area)
	TravelFeePerHelper  float64 // Travel fee per helper
	TotalTravelFee      float64 // Total travel fee (per helper * num helpers)
	Message             string  // Message to display (e.g., "within our service area - no travel fee")
}

// CalculateTravelFee calculates travel fee based on distance from office
// Rules:
// - Within 15 miles: $0 (no travel fee)
// - Outside 15 miles: Minimum $40 per helper, increases in $10 increments
// - Fee is per helper, so total = fee per helper * num helpers
// - Always rounds to avoid dealing with dimes and nickels
func CalculateTravelFee(distanceMiles float64, numHelpers int) *TravelFeeResult {
	result := &TravelFeeResult{
		DistanceMiles: distanceMiles,
	}

	// Round distance to 1 decimal place for display
	result.DistanceMiles = math.Round(distanceMiles*10) / 10

	// Check if within service area (15 miles)
	if distanceMiles <= 15.0 {
		result.IsWithinServiceArea = true
		result.TravelFee = 0
		result.TravelFeePerHelper = 0
		result.TotalTravelFee = 0
		result.Message = "within our service area - no travel fee"
		return result
	}

	// Outside service area - calculate fee
	// Minimum $40 per helper, increases in $10 increments based on distance
	// For every 10 miles beyond 15, add $10
	// Formula: $40 base + ($10 per 10 miles over 15)
	milesOverServiceArea := distanceMiles - 15.0
	
	// Calculate base fee per helper: $40 minimum
	// Add $10 for each 10-mile increment beyond 15 miles
	// Round up to nearest 10-mile increment
	var feePerHelper float64
	if milesOverServiceArea <= 10.0 {
		feePerHelper = 40.0
	} else {
		milesBeyondFirst10 := milesOverServiceArea - 10.0
		increments := math.Ceil(milesBeyondFirst10 / 10.0)
		feePerHelper = 40.0 + (increments * 10.0)
	}
	
	// Round to nearest dollar (no cents)
	feePerHelper = math.Round(feePerHelper)

	result.IsWithinServiceArea = false
	result.TravelFeePerHelper = feePerHelper
	result.TotalTravelFee = feePerHelper * float64(numHelpers)
	result.TravelFee = result.TotalTravelFee

	// Build message
	if numHelpers == 1 {
		result.Message = fmt.Sprintf("outside of our area — $%.0f travel fee", result.TotalTravelFee)
	} else {
		result.Message = fmt.Sprintf("outside of our area — $%.0f travel fee (for %d helpers)", result.TotalTravelFee, numHelpers)
	}

	return result
}
