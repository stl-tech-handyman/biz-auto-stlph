package pricing

import (
	"testing"
	"time"
)

func TestCalculateEstimate(t *testing.T) {
	tests := []struct {
		name          string
		eventDate     time.Time
		durationHours float64
		numHelpers    int
		expectError   bool
		checkResult   func(*testing.T, *EstimateResult)
	}{
		{
			name:          "basic calculation",
			eventDate:     time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			durationHours: 4,
			numHelpers:    2,
			expectError:   false,
			checkResult: func(t *testing.T, result *EstimateResult) {
				if result.TotalCost <= 0 {
					t.Errorf("expected totalCost > 0, got %f", result.TotalCost)
				}
				if result.NumHelpers != 2 {
					t.Errorf("expected numHelpers=2, got %d", result.NumHelpers)
				}
				if result.DurationHours != 4 {
					t.Errorf("expected durationHours=4, got %f", result.DurationHours)
				}
			},
		},
		{
			name:          "with extra hours",
			eventDate:     time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			durationHours: 6,
			numHelpers:    2,
			expectError:   false,
			checkResult: func(t *testing.T, result *EstimateResult) {
				if result.ExtraSubtotal <= 0 {
					t.Errorf("expected extraSubtotal > 0 for 6 hours, got %f", result.ExtraSubtotal)
				}
			},
		},
		{
			name:          "holiday date",
			eventDate:     time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
			durationHours: 4,
			numHelpers:    2,
			expectError:   false,
			checkResult: func(t *testing.T, result *EstimateResult) {
				if !result.IsSpecialDate {
					t.Error("expected isSpecialDate=true for Christmas")
				}
				if result.SpecialLabel == nil || *result.SpecialLabel != "Christmas Day" {
					t.Errorf("expected specialLabel='Christmas Day', got %v", result.SpecialLabel)
				}
			},
		},
		{
			name:          "invalid duration",
			eventDate:     time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			durationHours: -1,
			numHelpers:    2,
			expectError:   true,
		},
		{
			name:          "invalid helpers",
			eventDate:     time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			durationHours: 4,
			numHelpers:    0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateEstimate(tt.eventDate, tt.durationHours, tt.numHelpers)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestGetHolidayDatesForYear(t *testing.T) {
	holidays := GetHolidayDatesForYear(2025)

	expectedDates := []string{
		"2025-01-01",
		"2025-12-24",
		"2025-12-25",
		"2025-12-31",
	}

	for _, date := range expectedDates {
		if _, ok := holidays[date]; !ok {
			t.Errorf("expected holiday date %s not found", date)
		}
	}

	// Check Thanksgiving is calculated correctly
	thanksgivingDay := GetThanksgivingDay(2025)
	if thanksgivingDay < 22 || thanksgivingDay > 28 {
		t.Errorf("Thanksgiving day should be between 22-28, got %d", thanksgivingDay)
	}
}

func TestGetAllSpecialDates(t *testing.T) {
	result := GetAllSpecialDates(2, nil)

	if len(result) != 2 {
		t.Errorf("expected 2 years, got %d", len(result))
	}

	// Check that each year has holidays
	for year, dates := range result {
		if len(dates.Holidays) == 0 {
			t.Errorf("expected holidays for year %d", year)
		}
	}
}

