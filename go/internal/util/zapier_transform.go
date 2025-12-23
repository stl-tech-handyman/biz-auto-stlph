package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RawZapierPayload represents the raw input from Zapier webhook
type RawZapierPayload struct {
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	EmailAddress           string `json:"email_address"`
	PhoneNumber            string `json:"phone_number"`
	EventDate              string `json:"event_date"`
	EventTime              string `json:"event_time"`
	EventLocation          string `json:"event_location"`
	HelpersRequested       string `json:"helpers_requested"`
	ForHowManyHours        string `json:"for_how_many_hours"`
	Occasion               string `json:"occasion"`
	OccasionAsYouSeeIt     string `json:"occasion_as_you_see_it"`
	GuestsExpected         string `json:"guests_expected"`
	EventRole              string `json:"event_role"`
	EventRoleAsYouSeeIt    string `json:"event_role_as_you_see_it"`
	ScheduleCall           string `json:"schedule_call"`
	DryRun                 bool   `json:"dryRun"`
}

// TransformedLeadData represents the cleaned and normalized lead data
type TransformedLeadData struct {
	ClientName    string
	Email         string
	Phone         string
	EventDate     time.Time
	EventDateStr  string // Original format for display
	EventTime     string
	EventLocation string
	NumHelpers    int
	Duration      float64
	Occasion      string
	GuestCount    int
	Role          string
	ScheduleCall  bool
	DryRun        bool
}

// ConsolidateWithFallback consolidates a primary value with a fallback value
// If primary contains the trigger word (case-insensitive), returns fallback if available, otherwise default
// Otherwise returns primary if available, otherwise default
func ConsolidateWithFallback(primary, fallback, trigger, defaultValue string) string {
	if primary != "" && strings.Contains(strings.ToLower(primary), strings.ToLower(trigger)) {
		if fallback != "" {
			return fallback
		}
		return defaultValue
	}
	if primary != "" {
		return primary
	}
	return defaultValue
}

// ExtractFirstInteger extracts the first integer from a string
// Returns 0 if no integer is found
func ExtractFirstInteger(s string) int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindString(s)
	if matches == "" {
		return 0
	}
	val, _ := strconv.Atoi(matches)
	return val
}

// ExtractFirstFloat extracts the first floating-point number from a string
// Returns 0.0 if no number is found
func ExtractFirstFloat(s string) float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	matches := re.FindString(s)
	if matches == "" {
		return 0.0
	}
	val, _ := strconv.ParseFloat(matches, 64)
	return val
}

// ParseBooleanFromText checks if text contains a positive indicator (case-insensitive, whole word)
// Returns true if text contains any of the positive indicators as whole words, false otherwise
func ParseBooleanFromText(text string, positiveIndicators ...string) bool {
	if len(positiveIndicators) == 0 {
		positiveIndicators = []string{"yes", "true", "1"}
	}
	lower := strings.ToLower(strings.TrimSpace(text))
	for _, indicator := range positiveIndicators {
		indicatorLower := strings.ToLower(indicator)
		// Check for whole word match (word boundary or exact match)
		if lower == indicatorLower {
			return true
		}
		// Check if it's at the start followed by non-word character
		if strings.HasPrefix(lower, indicatorLower+" ") || strings.HasPrefix(lower, indicatorLower+",") {
			return true
		}
		// Check if it's in the middle with word boundaries
		if strings.Contains(lower, " "+indicatorLower+" ") || strings.Contains(lower, ","+indicatorLower+" ") {
			return true
		}
	}
	return false
}

// ExtractTimeComponent extracts and formats the time component from a date-time string
// Returns time formatted as "3:04 PM" or error if extraction fails
func ExtractTimeComponent(dateTimeStr string) (string, error) {
	formats := []string{
		"January 2, 2006 3:04 PM",
		"Jan 2, 2006 3:04 PM",
		"2006-01-02 3:04 PM",
		"2006-01-02 15:04",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateTimeStr); err == nil {
			return t.Format("3:04 PM"), nil
		}
	}

	// Fallback: try to extract time pattern with regex
	timeRe := regexp.MustCompile(`(\d{1,2}:\d{2}\s*(?:AM|PM|am|pm))`)
	matches := timeRe.FindString(dateTimeStr)
	if matches != "" {
		return matches, nil
	}

	return "", fmt.Errorf("unable to extract time from: %s", dateTimeStr)
}

// ParseDate parses a date string using multiple common formats
func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
		"2006/01/02",
		"January 2, 2006 3:04 PM",
		"Jan 2, 2006 3:04 PM",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}


// TransformZapierPayload transforms raw Zapier payload into normalized lead data
// Uses generalized utility functions to transform the data
func TransformZapierPayload(payload RawZapierPayload) (*TransformedLeadData, error) {
	// Consolidate occasion (if contains "other", use fallback)
	occasion := ConsolidateWithFallback(payload.Occasion, payload.OccasionAsYouSeeIt, "other", "Unspecified")

	// Extract duration (first float from string)
	duration := ExtractFirstFloat(payload.ForHowManyHours)
	if duration <= 0 {
		return nil, fmt.Errorf("invalid duration: %s", payload.ForHowManyHours)
	}

	// Extract helpers count (first integer from string)
	numHelpers := ExtractFirstInteger(payload.HelpersRequested)
	if numHelpers <= 0 {
		return nil, fmt.Errorf("invalid helpers requested: %s", payload.HelpersRequested)
	}

	// Consolidate role (if contains "other", use fallback)
	role := ConsolidateWithFallback(payload.EventRole, payload.EventRoleAsYouSeeIt, "other", "Unspecified")

	// Parse schedule call (check if contains "yes")
	scheduleCall := ParseBooleanFromText(payload.ScheduleCall, "yes")

	// Extract time from date-time (if event_time is not provided)
	eventTime := payload.EventTime
	if eventTime == "" && payload.EventDate != "" {
		extractedTime, err := ExtractTimeComponent(payload.EventDate)
		if err == nil {
			eventTime = extractedTime
		}
	}

	// Parse event date
	eventDate, err := ParseDate(payload.EventDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event date: %w", err)
	}

	// Extract guest count (first integer from string)
	guestCount := ExtractFirstInteger(payload.GuestsExpected)

	// Build client name
	clientName := strings.TrimSpace(fmt.Sprintf("%s %s", payload.FirstName, payload.LastName))
	if clientName == "" {
		return nil, fmt.Errorf("first_name and last_name are required")
	}

	// Validate email
	email := strings.TrimSpace(payload.EmailAddress)
	if email == "" {
		return nil, fmt.Errorf("email_address is required")
	}

	return &TransformedLeadData{
		ClientName:    clientName,
		Email:         email,
		Phone:         strings.TrimSpace(payload.PhoneNumber),
		EventDate:     eventDate,
		EventDateStr:  payload.EventDate, // Keep original for display
		EventTime:     eventTime,
		EventLocation: strings.TrimSpace(payload.EventLocation),
		NumHelpers:    numHelpers,
		Duration:      duration,
		Occasion:      occasion,
		GuestCount:    guestCount,
		Role:          role,
		ScheduleCall:  scheduleCall,
		DryRun:        payload.DryRun,
	}, nil
}

