package pricing

import (
	"math"
	"testing"
)

func TestCalculateTravelFee_WithinServiceArea(t *testing.T) {
	tests := []struct {
		name         string
		distanceMiles float64
		numHelpers   int
		wantFee      float64
		wantMessage  string
	}{
		{
			name:         "exactly 15 miles - within service area",
			distanceMiles: 15.0,
			numHelpers:   2,
			wantFee:      0,
			wantMessage:  "within our service area - no travel fee",
		},
		{
			name:         "10 miles - within service area",
			distanceMiles: 10.0,
			numHelpers:   1,
			wantFee:      0,
			wantMessage:  "within our service area - no travel fee",
		},
		{
			name:         "5 miles - within service area",
			distanceMiles: 5.0,
			numHelpers:   3,
			wantFee:      0,
			wantMessage:  "within our service area - no travel fee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTravelFee(tt.distanceMiles, tt.numHelpers)
			
			if result.IsWithinServiceArea != true {
				t.Errorf("IsWithinServiceArea = %v, want true", result.IsWithinServiceArea)
			}
			if result.TravelFee != tt.wantFee {
				t.Errorf("TravelFee = %.2f, want %.2f", result.TravelFee, tt.wantFee)
			}
			if result.TotalTravelFee != tt.wantFee {
				t.Errorf("TotalTravelFee = %.2f, want %.2f", result.TotalTravelFee, tt.wantFee)
			}
			if result.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", result.Message, tt.wantMessage)
			}
		})
	}
}

func TestCalculateTravelFee_OutsideServiceArea(t *testing.T) {
	tests := []struct {
		name         string
		distanceMiles float64
		numHelpers   int
		wantFeePerHelper float64
		wantTotalFee float64
	}{
		{
			name:         "16 miles - just outside (1 mile over), 1 helper",
			distanceMiles: 16.0,
			numHelpers:   1,
			wantFeePerHelper: 40.0, // $40 minimum for 1-10 miles over
			wantTotalFee: 40.0,
		},
		{
			name:         "16 miles - just outside (1 mile over), 2 helpers",
			distanceMiles: 16.0,
			numHelpers:   2,
			wantFeePerHelper: 40.0,
			wantTotalFee: 80.0, // $40 * 2
		},
		{
			name:         "25 miles - 10 miles over, 1 helper",
			distanceMiles: 25.0,
			numHelpers:   1,
			wantFeePerHelper: 40.0, // $40 for 1-10 miles over
			wantTotalFee: 40.0,
		},
		{
			name:         "25 miles - 10 miles over, 3 helpers",
			distanceMiles: 25.0,
			numHelpers:   3,
			wantFeePerHelper: 40.0,
			wantTotalFee: 120.0, // $40 * 3
		},
		{
			name:         "35 miles - 20 miles over, 2 helpers",
			distanceMiles: 35.0,
			numHelpers:   2,
			wantFeePerHelper: 50.0, // $40 + $10 (1 increment beyond first 10 miles)
			wantTotalFee: 100.0, // $50 * 2
		},
		{
			name:         "45 miles - 30 miles over, 1 helper",
			distanceMiles: 45.0,
			numHelpers:   1,
			wantFeePerHelper: 60.0, // $40 + $20 (2 increments beyond first 10 miles)
			wantTotalFee: 60.0,
		},
		{
			name:         "50 miles - 35 miles over, 4 helpers",
			distanceMiles: 50.0,
			numHelpers:   4,
			wantFeePerHelper: 70.0, // $40 + $30 (3 increments beyond first 10 miles)
			wantTotalFee: 280.0, // $70 * 4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTravelFee(tt.distanceMiles, tt.numHelpers)
			
			if result.IsWithinServiceArea != false {
				t.Errorf("IsWithinServiceArea = %v, want false", result.IsWithinServiceArea)
			}
			if math.Abs(result.TravelFeePerHelper-tt.wantFeePerHelper) > 0.01 {
				t.Errorf("TravelFeePerHelper = %.2f, want %.2f", result.TravelFeePerHelper, tt.wantFeePerHelper)
			}
			if math.Abs(result.TotalTravelFee-tt.wantTotalFee) > 0.01 {
				t.Errorf("TotalTravelFee = %.2f, want %.2f", result.TotalTravelFee, tt.wantTotalFee)
			}
			if result.TravelFee != result.TotalTravelFee {
				t.Errorf("TravelFee = %.2f, should equal TotalTravelFee = %.2f", result.TravelFee, result.TotalTravelFee)
			}
		})
	}
}

func TestCalculateTravelFee_MessageFormat(t *testing.T) {
	tests := []struct {
		name         string
		distanceMiles float64
		numHelpers   int
		wantContains string
	}{
		{
			name:         "single helper message",
			distanceMiles: 20.0,
			numHelpers:   1,
			wantContains: "outside of our area",
		},
		{
			name:         "multiple helpers message",
			distanceMiles: 20.0,
			numHelpers:   2,
			wantContains: "for 2 helpers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTravelFee(tt.distanceMiles, tt.numHelpers)
			
			if result.Message == "" {
				t.Error("Message should not be empty")
			}
			if !contains(result.Message, tt.wantContains) {
				t.Errorf("Message = %q, should contain %q", result.Message, tt.wantContains)
			}
		})
	}
}

func TestCalculateTravelFee_DistanceRounding(t *testing.T) {
	result := CalculateTravelFee(16.789, 1)
	
	// Distance should be rounded to 1 decimal place
	expectedDistance := 16.8
	if math.Abs(result.DistanceMiles-expectedDistance) > 0.01 {
		t.Errorf("DistanceMiles = %.2f, want %.2f", result.DistanceMiles, expectedDistance)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
