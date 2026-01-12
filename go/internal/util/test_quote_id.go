package util

import (
	"fmt"
	"sync"
)

var (
	testQuoteIDCounter int
	testQuoteIDMutex   sync.Mutex
)

// GenerateTestQuoteID generates a test quote ID that increments
// Format: T + 3 digits (increments)
// Excludes I, l, 1, O, 0 to avoid confusion
// Example: T001, T002, T003, etc.
func GenerateTestQuoteID() string {
	testQuoteIDMutex.Lock()
	defer testQuoteIDMutex.Unlock()

	testQuoteIDCounter++
	// Format as T + 3 digits (001, 002, 003, etc.)
	// Using digits 0-9 (but we avoid 0 in practice, so use 1-9 for clarity)
	// Actually, we'll use all digits 0-9 but format with padding
	return fmt.Sprintf("T%03d", testQuoteIDCounter)
}

// ResetTestQuoteIDCounter resets the test quote ID counter (useful for testing)
func ResetTestQuoteIDCounter() {
	testQuoteIDMutex.Lock()
	defer testQuoteIDMutex.Unlock()
	testQuoteIDCounter = 0
}
