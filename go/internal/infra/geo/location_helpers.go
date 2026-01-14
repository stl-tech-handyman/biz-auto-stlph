package geo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bizops360/go-api/internal/domain"
)

// LocationInfo holds location information for distance calculations
type LocationInfo struct {
	OriginAddress      string  // Address or coordinates for distance origin
	OriginLat          float64 // Origin latitude (if coordinates provided)
	OriginLng          float64 // Origin longitude (if coordinates provided)
	ServiceRadiusMiles float64 // Service radius in miles
}

// GetLocationInfo extracts location information from business config
// Returns location info with defaults if not configured
func GetLocationInfo(businessConfig *domain.BusinessConfig) LocationInfo {
	info := LocationInfo{
		ServiceRadiusMiles: 15.0, // Default service radius
	}

	if businessConfig == nil {
		// Fallback to constants if no business config
		info.OriginAddress = OfficeAddress
		info.OriginLat = OfficeLat
		info.OriginLng = OfficeLng
		return info
	}

	// Get office address
	info.OriginAddress = businessConfig.Location.OfficeAddress
	if info.OriginAddress == "" {
		info.OriginAddress = OfficeAddress // Fallback to constant
	}

	// PRIORITY 1: Check for explicit lat/lng fields (takes precedence - recommended for accuracy)
	if businessConfig.Location.Lat != 0 && businessConfig.Location.Lng != 0 {
		info.OriginLat = businessConfig.Location.Lat
		info.OriginLng = businessConfig.Location.Lng
		info.OriginAddress = fmt.Sprintf("%f,%f", info.OriginLat, info.OriginLng)
		if businessConfig.Location.ServiceRadiusMiles > 0 {
			info.ServiceRadiusMiles = businessConfig.Location.ServiceRadiusMiles
		}
		return info
	}

	// PRIORITY 2: Get distance origin (can be different from office address)
	distanceOrigin := businessConfig.Location.DistanceOrigin
	if distanceOrigin == "" {
		// Default to office address if not specified
		distanceOrigin = info.OriginAddress
	}

	// Check if distance origin is coordinates (lat,lng format)
	if strings.Contains(distanceOrigin, ",") {
		parts := strings.Split(distanceOrigin, ",")
		if len(parts) == 2 {
			if lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64); err == nil {
				if lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					info.OriginLat = lat
					info.OriginLng = lng
					info.OriginAddress = distanceOrigin // Use coordinates as address
					if businessConfig.Location.ServiceRadiusMiles > 0 {
						info.ServiceRadiusMiles = businessConfig.Location.ServiceRadiusMiles
					}
					return info
				}
			}
		}
	}

	// PRIORITY 3: Use distance origin as address (will be geocoded later if needed)
	info.OriginAddress = distanceOrigin

	// Get service radius
	if businessConfig.Location.ServiceRadiusMiles > 0 {
		info.ServiceRadiusMiles = businessConfig.Location.ServiceRadiusMiles
	}

	return info
}

// GetOriginCoordinates gets origin coordinates from location info
// If location info has coordinates, returns them directly
// Otherwise, geocodes the address using the geocoding service
func GetOriginCoordinates(ctx context.Context, locationInfo LocationInfo, geocodingService *GeocodingService) (lat, lng float64, err error) {
	// If we already have coordinates, use them
	if locationInfo.OriginLat != 0 && locationInfo.OriginLng != 0 {
		return locationInfo.OriginLat, locationInfo.OriginLng, nil
	}

	// Otherwise, geocode the address
	if geocodingService == nil {
		// Fallback to constants if no geocoding service
		return OfficeLat, OfficeLng, nil
	}

	result, err := geocodingService.GetLatLng(ctx, locationInfo.OriginAddress)
	if err != nil {
		// Fallback to constants on error
		return OfficeLat, OfficeLng, fmt.Errorf("failed to geocode origin address, using defaults: %w", err)
	}

	return result.Lat, result.Lng, nil
}
