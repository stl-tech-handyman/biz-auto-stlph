package util

import (
	"fmt"
	"strconv"
)

// DollarsToCents converts dollars to cents
func DollarsToCents(dollars float64) int64 {
	return int64(dollars * 100)
}

// CentsToDollars converts cents to dollars
func CentsToDollars(cents int64) float64 {
	return float64(cents) / 100
}

// ParseDollarAmount parses a dollar amount from string or float64, converting to cents
// If the value is less than 10000, it's assumed to be in dollars and will be converted
func ParseDollarAmount(value interface{}) (int64, error) {
	switch v := value.(type) {
	case float64:
		// If less than 10000, assume it's in dollars
		if v < 10000 {
			return DollarsToCents(v), nil
		}
		// Otherwise assume it's already in cents
		return int64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		if f < 10000 {
			return DollarsToCents(f), nil
		}
		return int64(f), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type for dollar amount: %T", value)
	}
}

