package util

import (
	"fmt"
	"strings"
)

// GenerateShortQuoteID generates a 4-character alphanumeric quote ID
// Based on email and event date (matching Apps Script implementation)
func GenerateShortQuoteID(email, eventDate string) string {
	if eventDate == "" || email == "" {
		return "XXXX"
	}

	// Remove dashes from date (YYYY-MM-DD -> YYYYMMDD)
	datePart := strings.ReplaceAll(eventDate, "-", "")

	// Create hash from date + email
	rawString := datePart + email
	hash := 0
	for _, char := range rawString {
		hash += int(char)
	}

	// Convert hash to 4-character alphanumeric code (base36)
	shortID := fmt.Sprintf("%04X", hash%1679616) // 1679616 = 36^4
	if len(shortID) > 4 {
		shortID = shortID[:4]
	}
	// Pad with X if needed
	for len(shortID) < 4 {
		shortID = "X" + shortID
	}

	return shortID
}



