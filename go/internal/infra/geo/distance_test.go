package geo

import (
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	officeLat := 38.6255
	officeLng := -90.2456
	
	tests := []struct {
		name     string
		lat1     float64
		lng1     float64
		lat2     float64
		lng2     float64
		wantDist float64 // Approximate expected distance in miles
		tolerance float64 // Tolerance for comparison
	}{
		{
			name:      "same location",
			lat1:      officeLat,
			lng1:      officeLng,
			lat2:      officeLat,
			lng2:      officeLng,
			wantDist:  0.0,
			tolerance: 0.1,
		},
		{
			name:      "office to nearby location (1 mile away)",
			lat1:      officeLat,
			lng1:      officeLng,
			lat2:      officeLat + 0.0145, // ~1 mile north
			lng2:      officeLng,
			wantDist:  1.0,
			tolerance: 0.2,
		},
		{
			name:      "office to downtown St. Louis (~2.5 miles)",
			lat1:      officeLat,
			lng1:      officeLng,
			lat2:      38.6270, // Downtown St. Louis
			lng2:      -90.1994,
			wantDist:  2.5,
			tolerance: 0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDistance(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			
			if math.Abs(got-tt.wantDist) > tt.tolerance {
				t.Errorf("CalculateDistance() = %.2f miles, want %.2f miles (tolerance: %.2f)", 
					got, tt.wantDist, tt.tolerance)
			}
		})
	}
}

func TestCalculateDistanceFromOffice(t *testing.T) {
	officeLat := 38.6255
	officeLng := -90.2456
	
	tests := []struct {
		name     string
		lat      float64
		lng      float64
		wantDist float64
		tolerance float64
	}{
		{
			name:      "office location itself",
			lat:       officeLat,
			lng:       officeLng,
			wantDist:  0.0,
			tolerance: 0.1,
		},
		{
			name:      "location 10 miles away",
			lat:       officeLat + 0.145, // ~10 miles north
			lng:       officeLng,
			wantDist:  10.0,
			tolerance: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDistanceFromOffice(tt.lat, tt.lng)
			
			if math.Abs(got-tt.wantDist) > tt.tolerance {
				t.Errorf("CalculateDistanceFromOffice() = %.2f miles, want %.2f miles (tolerance: %.2f)", 
					got, tt.wantDist, tt.tolerance)
			}
		})
	}
}

func TestIsWithinServiceArea(t *testing.T) {
	officeLat := 38.6255
	officeLng := -90.2456
	
	tests := []struct {
		name     string
		lat      float64
		lng      float64
		want     bool
	}{
		{
			name: "within 15 miles",
			lat:  officeLat + 0.1, // ~7 miles
			lng:  officeLng,
			want: true,
		},
		{
			name: "exactly 15 miles",
			lat:  officeLat + 0.2165, // ~15 miles (adjusted to be exactly at boundary)
			lng:  officeLng,
			want: true,
		},
		{
			name: "outside 15 miles",
			lat:  officeLat + 0.3, // ~20 miles
			lng:  officeLng,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsWithinServiceArea(tt.lat, tt.lng)
			
			if got != tt.want {
				t.Errorf("IsWithinServiceArea() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRadiusConstants(t *testing.T) {
	if ServiceRadiusMiles != 15.0 {
		t.Errorf("ServiceRadiusMiles = %.2f, want 15.0", ServiceRadiusMiles)
	}
	
	expectedAddress := "4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"
	if OfficeAddress != expectedAddress {
		t.Errorf("OfficeAddress = %q, want %q", OfficeAddress, expectedAddress)
	}
}
