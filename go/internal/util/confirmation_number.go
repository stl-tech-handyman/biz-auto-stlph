package util

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// GenerateConfirmationNumber generates a unique confirmation number
// Format: Letter (based on year starting from 2025) + 3 digits (from month/date hash)
// Excludes I, l, 1, O, 0 to avoid confusion (uses digits: 2, 3, 4, 5, 6, 7, 8, 9)
// Year mapping: 2025 = A, 2026 = B, 2027 = C, etc. (skips I and O)
func GenerateConfirmationNumber(email, occasion string, eventDate time.Time) string {
	// Get year letter (2025 = A, 2026 = B, etc.)
	// Skip I and O to avoid confusion with 1 and 0
	year := eventDate.Year()
	yearOffset := year - 2025
	
	// Valid letters excluding I and O: A-H, J-N, P-Z (24 letters)
	// Map year offset to valid letters
	validLetters := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	yearLetter := validLetters[yearOffset%len(validLetters)]
	
	// Create hash from month, date, email, and occasion for uniqueness
	month := int(eventDate.Month())
	day := eventDate.Day()
	seed := fmt.Sprintf("%d|%d|%s|%s", month, day, strings.ToLower(email), strings.ToLower(occasion))
	hash := sha256.Sum256([]byte(seed))
	
	// Valid digits excluding 0 and 1: 2, 3, 4, 5, 6, 7, 8, 9
	validDigits := []byte{'2', '3', '4', '5', '6', '7', '8', '9'}
	
	// Generate 3 digits from hash bytes
	digits := make([]byte, 3)
	for i := 0; i < 3; i++ {
		// Use hash bytes to select a digit
		index := int(hash[i]) % len(validDigits)
		digits[i] = validDigits[index]
	}
	
	return string(yearLetter) + string(digits)
}



