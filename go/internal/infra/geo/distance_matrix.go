package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// DistanceMatrixService handles distance calculations via Google Distance Matrix API
type DistanceMatrixService struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewDistanceMatrixService creates a new distance matrix service
func NewDistanceMatrixService() (*DistanceMatrixService, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		// Try to get from Secret Manager reference
		apiKey = os.Getenv("MAPS_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_MAPS_API_KEY or MAPS_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("GOOGLE_MAPS_DISTANCE_MATRIX_URL")
	if baseURL == "" {
		baseURL = "https://maps.googleapis.com/maps/api/distancematrix/json"
	}

	return &DistanceMatrixService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}, nil
}

// DistanceMatrixResult represents the result from Distance Matrix API
type DistanceMatrixResult struct {
	DistanceMiles float64 // Distance in miles
	DurationMins  float64 // Duration in minutes (optional, for future use)
	Status        string   // API status (OK, NOT_FOUND, etc.)
}

// GetDrivingDistance calculates the driving distance in miles between two addresses
// origin and destination can be either addresses (strings) or coordinates (lat,lng format)
// Returns distance in miles, or error if API call fails
func (d *DistanceMatrixService) GetDrivingDistance(ctx context.Context, origin, destination string) (*DistanceMatrixResult, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	// Build request URL
	params := url.Values{}
	params.Set("origins", origin)
	params.Set("destinations", destination)
	params.Set("units", "imperial") // Get results in miles
	params.Set("key", d.apiKey)

	reqURL := fmt.Sprintf("%s?%s", d.baseURL, params.Encode())

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call distance matrix API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("distance matrix API returned status %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Status            string   `json:"status"`
		OriginAddresses   []string `json:"origin_addresses"`
		DestinationAddresses []string `json:"destination_addresses"`
		Rows              []struct {
			Elements []struct {
				Status string `json:"status"`
				Distance struct {
					Text  string  `json:"text"`
					Value float64 `json:"value"` // Value in meters
				} `json:"distance"`
				Duration struct {
					Text  string  `json:"text"`
					Value float64 `json:"value"` // Value in seconds
				} `json:"duration"`
			} `json:"elements"`
		} `json:"rows"`
		ErrorMessage string `json:"error_message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse distance matrix response: %w", err)
	}

	if result.Status != "OK" {
		errorMsg := result.ErrorMessage
		if errorMsg == "" {
			errorMsg = result.Status
		}
		return nil, fmt.Errorf("distance matrix API failed: %s", errorMsg)
	}

	if len(result.Rows) == 0 || len(result.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("no results found in distance matrix response")
	}

	element := result.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("distance matrix element status: %s", element.Status)
	}

	// Convert meters to miles (1 meter = 0.000621371 miles)
	distanceMiles := element.Distance.Value * 0.000621371
	
	// Convert seconds to minutes
	durationMins := element.Duration.Value / 60.0

	return &DistanceMatrixResult{
		DistanceMiles: distanceMiles,
		DurationMins:  durationMins,
		Status:        element.Status,
	}, nil
}

// GetDrivingDistanceFromOffice calculates the driving distance from the office to a destination address
// DEPRECATED: Use GetDrivingDistanceFromOrigin with location from business config instead
func (d *DistanceMatrixService) GetDrivingDistanceFromOffice(ctx context.Context, destination string) (*DistanceMatrixResult, error) {
	return d.GetDrivingDistance(ctx, OfficeAddress, destination)
}

// GetDrivingDistanceFromOrigin calculates the driving distance from an origin to a destination
// origin: Origin address or coordinates (lat,lng format)
// destination: Destination address or coordinates (lat,lng format)
func (d *DistanceMatrixService) GetDrivingDistanceFromOrigin(ctx context.Context, origin, destination string) (*DistanceMatrixResult, error) {
	return d.GetDrivingDistance(ctx, origin, destination)
}

// GetDrivingDistanceFromOriginCoords calculates the driving distance from origin coordinates to destination coordinates
// originLat, originLng: Origin coordinates
// destLat, destLng: Destination coordinates
func (d *DistanceMatrixService) GetDrivingDistanceFromOriginCoords(ctx context.Context, originLat, originLng, destLat, destLng float64) (*DistanceMatrixResult, error) {
	destCoords := fmt.Sprintf("%f,%f", destLat, destLng)
	originCoords := fmt.Sprintf("%f,%f", originLat, originLng)
	return d.GetDrivingDistance(ctx, originCoords, destCoords)
}

// GetDrivingDistanceFromOfficeCoords calculates the driving distance from the office to destination coordinates
// DEPRECATED: Use GetDrivingDistanceFromOriginCoords with location from business config instead
func (d *DistanceMatrixService) GetDrivingDistanceFromOfficeCoords(ctx context.Context, destLat, destLng float64) (*DistanceMatrixResult, error) {
	destCoords := fmt.Sprintf("%f,%f", destLat, destLng)
	originCoords := fmt.Sprintf("%f,%f", OfficeLat, OfficeLng)
	return d.GetDrivingDistance(ctx, originCoords, destCoords)
}
