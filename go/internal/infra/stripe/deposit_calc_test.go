package stripe

import (
	"testing"
)

func TestCalculateDepositFromEstimate(t *testing.T) {
	tests := []struct {
		name            string
		estimateCents   int64
		expectedMin     int64
		expectedMax     int64
		checkPercentage func(*testing.T, float64)
	}{
		{
			name:          "small estimate",
			estimateCents: 10000, // $100
			expectedMin:   2500,  // 25% = $25
			expectedMax:   4000,  // 40% = $40
			checkPercentage: func(t *testing.T, pct float64) {
				// For small estimates, rounding to professional amounts may exceed 40%
				// This is acceptable behavior
				if pct < 25 {
					t.Errorf("expected percentage >= 25%%, got %.1f%%", pct)
				}
			},
		},
		{
			name:          "medium estimate",
			estimateCents: 100000, // $1000
			expectedMin:   25000,  // 25% = $250
			expectedMax:   40000,  // 40% = $400
			checkPercentage: func(t *testing.T, pct float64) {
				if pct < 25 || pct > 40 {
					t.Errorf("expected percentage between 25-40%%, got %.1f%%", pct)
				}
			},
		},
		{
			name:          "large estimate",
			estimateCents: 500000, // $5000
			expectedMin:   125000, // 25% = $1250
			expectedMax:   200000, // 40% = $2000
			checkPercentage: func(t *testing.T, pct float64) {
				if pct < 25 || pct > 40 {
					t.Errorf("expected percentage between 25-40%%, got %.1f%%", pct)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := CalculateDepositFromEstimate(tt.estimateCents)

			// For small estimates, value may exceed max due to professional rounding
			if calc.Value < tt.expectedMin {
				t.Errorf("expected deposit >= %d cents, got %d", tt.expectedMin, calc.Value)
			}

			if calc.MinAmount != tt.expectedMin {
				t.Errorf("expected minAmount=%d, got %d", tt.expectedMin, calc.MinAmount)
			}

			if calc.MaxAmount != tt.expectedMax {
				t.Errorf("expected maxAmount=%d, got %d", tt.expectedMax, calc.MaxAmount)
			}

			if tt.checkPercentage != nil {
				tt.checkPercentage(t, calc.Percentage)
			}

			// Check that value is a professional amount (multiple of $50)
			if calc.Value%5000 != 0 {
				t.Errorf("expected deposit to be multiple of $50 (5000 cents), got %d", calc.Value)
			}
		})
	}
}

func TestRoundUpToProfessionalAmount(t *testing.T) {
	candidates := []int64{5000, 10000, 15000, 20000} // $50, $100, $150, $200

	tests := []struct {
		name     string
		amount   int64
		expected int64
	}{
		{
			name:     "below minimum",
			amount:   1000,
			expected: 5000,
		},
		{
			name:     "exact match",
			amount:   10000,
			expected: 10000,
		},
		{
			name:     "between values",
			amount:   7500,
			expected: 10000,
		},
		{
			name:     "above maximum",
			amount:   50000,
			expected: 20000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundUpToProfessionalAmount(tt.amount, candidates)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

