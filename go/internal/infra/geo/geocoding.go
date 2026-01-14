package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/domain"
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

// OfficeLocation represents the office location
// DEPRECATED: These constants are kept for backward compatibility.
// New code should use LocationConfig from business configuration.
const (
	OfficeLat          = 38.6255  // 4220 Duncan Ave., Ste. 201, St. Louis, MO 63110
	OfficeLng          = -90.2456
	OfficeAddress      = "4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"
	ServiceRadiusMiles = 15.0     // 15 miles service radius
)

// CalculateDistance calculates the distance in miles between two lat/lng coordinates using Haversine formula
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusMiles = 3958.8 // Earth's radius in miles

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlng := lng2Rad - lng1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlng/2)*math.Sin(dlng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusMiles * c
	return distance
}

// CalculateDistanceFromOffice calculates the distance in miles from the office to a given location
// DEPRECATED: Use CalculateDistanceFromOrigin with location from business config instead
func CalculateDistanceFromOffice(lat, lng float64) float64 {
	return CalculateDistance(OfficeLat, OfficeLng, lat, lng)
}

// CalculateDistanceFromOrigin calculates the distance in miles from an origin point to a destination
// originLat, originLng: Origin coordinates
// destLat, destLng: Destination coordinates
func CalculateDistanceFromOrigin(originLat, originLng, destLat, destLng float64) float64 {
	return CalculateDistance(originLat, originLng, destLat, destLng)
}

// IsWithinServiceArea checks if a location is within the service area (15 miles)
// DEPRECATED: Use IsWithinServiceAreaWithRadius with location from business config instead
func IsWithinServiceArea(lat, lng float64) bool {
	distance := CalculateDistanceFromOffice(lat, lng)
	return distance <= ServiceRadiusMiles
}

// IsWithinServiceAreaWithRadius checks if a location is within the service area
// originLat, originLng: Origin coordinates (e.g., from business config)
// destLat, destLng: Destination coordinates
// radiusMiles: Service radius in miles (e.g., from business config)
func IsWithinServiceAreaWithRadius(originLat, originLng, destLat, destLng, radiusMiles float64) bool {
	distance := CalculateDistanceFromOrigin(originLat, originLng, destLat, destLng)
	return distance <= radiusMiles
}

// LocationFromBusinessConfig extracts location coordinates and radius from business config
// Returns lat, lng, radiusMiles, and an error if coordinates cannot be determined
// Priority: 1) Explicit lat/lng from config, 2) Parse from distanceOrigin if coordinates, 3) Geocode distanceOrigin if address
func LocationFromBusinessConfig(ctx context.Context, config *domain.LocationConfig, geocodingService *GeocodingService) (lat, lng, radiusMiles float64, err error) {
	// Default radius
	radiusMiles = 15.0
	if config.ServiceRadiusMiles > 0 {
		radiusMiles = config.ServiceRadiusMiles
	}

	// Priority 1: Use explicit lat/lng if provided
	if config.Lat != 0 && config.Lng != 0 {
		return config.Lat, config.Lng, radiusMiles, nil
	}

	// Priority 2: Try to parse coordinates from distanceOrigin string
	if config.DistanceOrigin != "" {
		// Check if it's in "lat,lng" format
		if lat, lng, err := parseCoordinates(config.DistanceOrigin); err == nil {
			return lat, lng, radiusMiles, nil
		}

		// Priority 3: Geocode the address if geocoding service is available
		if geocodingService != nil {
			geoResult, err := geocodingService.GetLatLng(ctx, config.DistanceOrigin)
			if err == nil {
				return geoResult.Lat, geoResult.Lng, radiusMiles, nil
			}
		}
	}

	// Fallback: Try to geocode officeAddress
	if config.OfficeAddress != "" && geocodingService != nil {
		geoResult, err := geocodingService.GetLatLng(ctx, config.OfficeAddress)
		if err == nil {
			return geoResult.Lat, geoResult.Lng, radiusMiles, nil
		}
	}

	return 0, 0, radiusMiles, fmt.Errorf("unable to determine location coordinates from config")
}

// parseCoordinates parses coordinates from a string in "lat,lng" format
// Returns lat, lng, and error if parsing fails
func parseCoordinates(coordsStr string) (float64, float64, error) {
	coordsStr = strings.TrimSpace(coordsStr)
	parts := strings.Split(coordsStr, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate format, expected 'lat,lng'")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %w", err)
	}

	return lat, lng, nil
}
