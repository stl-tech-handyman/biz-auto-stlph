package stripe

import (
	"testing"
)

// TestCalculateDepositFromEstimateRange tests deposit calculation across wide range
func TestCalculateDepositFromEstimateRange(t *testing.T) {
	testCases := []struct {
		name         string
		estimateCents int64
		minPercent   float64
		maxPercent   float64
	}{
		{"Small estimate ($100)", 10000, 0.15, 0.30},
		{"Medium estimate ($500)", 50000, 0.15, 0.30},
		{"Large estimate ($1000)", 100000, 0.15, 0.30},
		{"Very large estimate ($5000)", 500000, 0.15, 0.30},
		{"Edge case ($50)", 5000, 0.15, 0.30},
		{"Edge case ($10000)", 1000000, 0.15, 0.30},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calc := CalculateDepositFromEstimate(tc.estimateCents)
			
			// Verify deposit is within range
			depositPercent := float64(calc.Value) / float64(tc.estimateCents)
			if depositPercent < tc.minPercent || depositPercent > tc.maxPercent {
				t.Errorf("Deposit percentage %.2f%% is outside expected range (%.0f%%-%.0f%%)",
					depositPercent*100, tc.minPercent*100, tc.maxPercent*100)
			}
			
			// Verify deposit is a professional amount (multiple of $50)
			if calc.Value%5000 != 0 {
				t.Errorf("Deposit %d cents is not a multiple of $50 (5000 cents)", calc.Value)
			}
			
			// Verify deposit is not zero
			if calc.Value == 0 {
				t.Error("Deposit amount is zero")
			}
		})
	}
}

// TestDepositCalculationConsistencyRange tests that same estimate always gives same deposit
func TestDepositCalculationConsistencyRange(t *testing.T) {
	estimateCents := int64(100000) // $1000
	
	firstCalc := CalculateDepositFromEstimate(estimateCents)
	
	// Run calculation 10 times
	for i := 0; i < 10; i++ {
		calc := CalculateDepositFromEstimate(estimateCents)
		if calc.Value != firstCalc.Value {
			t.Errorf("Inconsistent deposit calculation: first=%d, iteration %d=%d", firstCalc.Value, i, calc.Value)
		}
		if calc.Percentage != firstCalc.Percentage {
			t.Errorf("Inconsistent deposit percentage: first=%.2f%%, iteration %d=%.2f%%",
				firstCalc.Percentage*100, i, calc.Percentage*100)
		}
	}
}
