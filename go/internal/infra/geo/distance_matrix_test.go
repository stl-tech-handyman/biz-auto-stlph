package geo

import (
	"context"
	"os"
	"testing"
)

func TestDistanceMatrixService_GetDrivingDistance(t *testing.T) {
	// Skip if API key is not set (for CI/CD environments)
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MAPS_API_KEY")
	}
	if apiKey == "" {
		t.Skip("Skipping test: GOOGLE_MAPS_API_KEY or MAPS_API_KEY not set")
	}

	service, err := NewDistanceMatrixService()
	if err != nil {
		t.Fatalf("Failed to create DistanceMatrixService: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		origin      string
		destination string
		wantErr     bool
		checkDistance func(float64) bool // Check if distance is reasonable
	}{
		{
			name:        "office to nearby location",
			origin:      OfficeAddress,
			destination: "2300 Hitzert Ct, Fenton, MO 63026",
			wantErr:     false,
			checkDistance: func(d float64) bool {
				// Should be roughly 10-20 miles (driving distance is longer than straight-line)
				return d > 5 && d < 30
			},
		},
		{
			name:        "office to downtown St. Louis",
			origin:      OfficeAddress,
			destination: "Gateway Arch, St. Louis, MO",
			wantErr:     false,
			checkDistance: func(d float64) bool {
				// Should be roughly 2-5 miles
				return d > 1 && d < 10
			},
		},
		{
			name:        "using coordinates",
			origin:      "38.6255,-90.2456", // Office coordinates
			destination: "38.6275,-90.2000", // Nearby coordinates
			wantErr:     false,
			checkDistance: func(d float64) bool {
				// Should be a few miles
				return d > 0 && d < 10
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetDrivingDistance(ctx, tt.origin, tt.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDrivingDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result == nil {
					t.Fatal("GetDrivingDistance() returned nil result")
				}
				if result.DistanceMiles <= 0 {
					t.Errorf("GetDrivingDistance() distance = %v, want > 0", result.DistanceMiles)
				}
				if !tt.checkDistance(result.DistanceMiles) {
					t.Errorf("GetDrivingDistance() distance = %v miles, failed distance check", result.DistanceMiles)
				}
				if result.DurationMins <= 0 {
					t.Errorf("GetDrivingDistance() duration = %v, want > 0", result.DurationMins)
				}
				if result.Status != "OK" {
					t.Errorf("GetDrivingDistance() status = %v, want OK", result.Status)
				}
			}
		})
	}
}

func TestDistanceMatrixService_GetDrivingDistanceFromOffice(t *testing.T) {
	// Skip if API key is not set
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MAPS_API_KEY")
	}
	if apiKey == "" {
		t.Skip("Skipping test: GOOGLE_MAPS_API_KEY or MAPS_API_KEY not set")
	}

	service, err := NewDistanceMatrixService()
	if err != nil {
		t.Fatalf("Failed to create DistanceMatrixService: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		destination string
		wantErr     bool
	}{
		{
			name:        "valid destination address",
			destination: "2300 Hitzert Ct, Fenton, MO 63026",
			wantErr:     false,
		},
		{
			name:        "downtown St. Louis",
			destination: "Gateway Arch, St. Louis, MO",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetDrivingDistanceFromOffice(ctx, tt.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDrivingDistanceFromOffice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result == nil {
					t.Fatal("GetDrivingDistanceFromOffice() returned nil result")
				}
				if result.DistanceMiles <= 0 {
					t.Errorf("GetDrivingDistanceFromOffice() distance = %v, want > 0", result.DistanceMiles)
				}
			}
		})
	}
}

func TestDistanceMatrixService_GetDrivingDistanceFromOfficeCoords(t *testing.T) {
	// Skip if API key is not set
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MAPS_API_KEY")
	}
	if apiKey == "" {
		t.Skip("Skipping test: GOOGLE_MAPS_API_KEY or MAPS_API_KEY not set")
	}

	service, err := NewDistanceMatrixService()
	if err != nil {
		t.Fatalf("Failed to create DistanceMatrixService: %v", err)
	}

	ctx := context.Background()

	// Test with coordinates near office
	result, err := service.GetDrivingDistanceFromOfficeCoords(ctx, 38.6275, -90.2000)
	if err != nil {
		t.Fatalf("GetDrivingDistanceFromOfficeCoords() error = %v", err)
	}

	if result.DistanceMiles <= 0 {
		t.Errorf("GetDrivingDistanceFromOfficeCoords() distance = %v, want > 0", result.DistanceMiles)
	}
}

func TestDistanceMatrixService_ErrorHandling(t *testing.T) {
	// Skip if API key is not set
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MAPS_API_KEY")
	}
	if apiKey == "" {
		t.Skip("Skipping test: GOOGLE_MAPS_API_KEY or MAPS_API_KEY not set")
	}

	service, err := NewDistanceMatrixService()
	if err != nil {
		t.Fatalf("Failed to create DistanceMatrixService: %v", err)
	}

	ctx := context.Background()

	// Test with invalid address
	_, err = service.GetDrivingDistance(ctx, "", "some destination")
	if err == nil {
		t.Error("GetDrivingDistance() with empty origin should return error")
	}

	_, err = service.GetDrivingDistance(ctx, "some origin", "")
	if err == nil {
		t.Error("GetDrivingDistance() with empty destination should return error")
	}
}
