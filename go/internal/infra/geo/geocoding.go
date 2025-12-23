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

// GeocodingService handles address geocoding via Google Maps API
type GeocodingService struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewGeocodingService creates a new geocoding service
func NewGeocodingService() (*GeocodingService, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		// Try to get from Secret Manager reference
		apiKey = os.Getenv("MAPS_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_MAPS_API_KEY or MAPS_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("GOOGLE_MAPS_GEOCODE_URL")
	if baseURL == "" {
		baseURL = "https://maps.googleapis.com/maps/api/geocode/json"
	}

	return &GeocodingService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}, nil
}

// GeocodeResult represents geocoding result
type GeocodeResult struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	FullAddress string  `json:"fullAddress"`
}

// GetLatLng geocodes an address and returns lat/lng coordinates
// Matches the Apps Script getLatLng function
func (g *GeocodingService) GetLatLng(ctx context.Context, address string) (*GeocodeResult, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Build request URL
	params := url.Values{}
	params.Set("address", address)
	params.Set("key", g.apiKey)

	reqURL := fmt.Sprintf("%s?%s", g.baseURL, params.Encode())

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode address: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Status  string `json:"status"`
		Results []struct {
			FormattedAddress string `json:"formatted_address"`
			Geometry         struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		ErrorMessage string `json:"error_message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if result.Status != "OK" {
		errorMsg := result.ErrorMessage
		if errorMsg == "" {
			errorMsg = result.Status
		}
		return nil, fmt.Errorf("geocoding failed: %s", errorMsg)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no results found for address")
	}

	location := result.Results[0].Geometry.Location

	return &GeocodeResult{
		Lat:         location.Lat,
		Lng:         location.Lng,
		FullAddress: result.Results[0].FormattedAddress,
	}, nil
}



