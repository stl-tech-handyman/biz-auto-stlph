# Location Configuration Guide

## Overview

The quote engine uses configurable location settings to calculate travel fees and determine service area boundaries. These settings are stored in business configuration YAML files and can be customized per business.

## Configuration File Location

Business location settings are stored in:
```
config/businesses/{businessId}.yaml
```

For example:
```
config/businesses/stlpartyhelpers.yaml
```

## Configuration Structure

Location settings are defined in the `location` section of the business configuration:

```yaml
location:
  # Office address (for display purposes)
  officeAddress: "4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"
  
  # Distance origin - the starting point for all distance calculations
  # Can be an address string or coordinates in "lat,lng" format
  # If not specified, defaults to officeAddress
  distanceOrigin: "4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"
  
  # Latitude and Longitude (RECOMMENDED for accuracy)
  # If provided, these take precedence over distanceOrigin string
  # Avoids geocoding API calls and ensures accurate calculations
  lat: 38.6255   # Latitude of the distance origin
  lng: -90.2456  # Longitude of the distance origin
  
  # Service radius in miles
  # Locations within this radius have no travel fee
  serviceRadiusMiles: 15.0
```

## Field Descriptions

### `officeAddress` (string, optional)
- **Purpose**: Display address for the business office
- **Usage**: Shown in emails and UI for reference
- **Example**: `"4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"`

### `distanceOrigin` (string, optional)
- **Purpose**: Starting point for distance calculations
- **Format**: Can be either:
  - An address string: `"4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"`
  - Coordinates: `"38.6255,-90.2456"` (lat,lng format)
- **Usage**: Used if `lat`/`lng` are not provided
- **Fallback**: Defaults to `officeAddress` if not specified

### `lat` (float, optional but recommended)
- **Purpose**: Latitude of the distance origin point
- **Precedence**: Takes priority over `distanceOrigin` if provided
- **Benefits**: 
  - Avoids geocoding API calls
  - Ensures accurate calculations
  - Faster processing
- **Example**: `38.6255`

### `lng` (float, optional but recommended)
- **Purpose**: Longitude of the distance origin point
- **Precedence**: Takes priority over `distanceOrigin` if provided
- **Benefits**: Same as `lat`
- **Example**: `-90.2456`

### `serviceRadiusMiles` (float, optional)
- **Purpose**: Service area radius in miles
- **Default**: `15.0` miles if not specified
- **Usage**: Locations within this radius have no travel fee
- **Example**: `15.0`

## Priority Order

The system determines location coordinates in the following priority order:

1. **Explicit `lat`/`lng`** (highest priority - recommended)
   - If both `lat` and `lng` are provided, these are used directly
   - No API calls required
   - Most accurate and fastest

2. **Parse `distanceOrigin` as coordinates**
   - If `distanceOrigin` is in "lat,lng" format, it's parsed
   - Example: `"38.6255,-90.2456"`

3. **Geocode `distanceOrigin` address**
   - If `distanceOrigin` is an address string, it's geocoded using Google Maps API
   - Requires API key and network call

4. **Geocode `officeAddress`** (fallback)
   - If all else fails, `officeAddress` is geocoded
   - Last resort fallback

## Code Implementation

### Domain Model
Location configuration is defined in:
```
go/internal/domain/business.go
```

The `LocationConfig` struct:
```go
type LocationConfig struct {
    OfficeAddress      string  `yaml:"officeAddress" json:"officeAddress"`
    DistanceOrigin     string  `yaml:"distanceOrigin" json:"distanceOrigin"`
    Lat                float64 `yaml:"lat" json:"lat"`
    Lng                float64 `yaml:"lng" json:"lng"`
    ServiceRadiusMiles float64 `yaml:"serviceRadiusMiles" json:"serviceRadiusMiles"`
}
```

### Helper Functions
Location extraction functions are in:
```
go/internal/infra/geo/geocoding.go
```

Key function:
```go
func LocationFromBusinessConfig(
    ctx context.Context,
    config *domain.LocationConfig,
    geocodingService *GeocodingService,
) (lat, lng, radiusMiles float64, err error)
```

### Usage in Email Handler
Location configuration is used in:
```
go/internal/http/handlers/email_handler.go
```

The email handler:
1. Loads business config using `BusinessLoader`
2. Extracts location using `LocationFromBusinessConfig`
3. Calculates distance from origin to event location
4. Determines travel fee based on service radius

## Travel Fee Calculation

Travel fees are calculated based on:
- Distance from origin (configured location) to event location
- Service radius (from `serviceRadiusMiles`)
- Number of helpers

### Rules:
- **Within service radius**: No travel fee
- **Outside service radius**: 
  - Minimum $40 per helper
  - Increases in $10 increments for every 10 miles beyond the first 10 miles outside the radius
  - Total = fee per helper × number of helpers

### Example:
- Service radius: 15 miles
- Event location: 25 miles away
- Helpers: 2
- Calculation:
  - Distance outside: 25 - 15 = 10 miles
  - Fee per helper: $40 (within first 10 miles)
  - Total: $40 × 2 = $80

## Best Practices

1. **Always provide `lat`/`lng`**
   - Most accurate
   - Avoids API calls
   - Faster processing

2. **Keep `officeAddress` for display**
   - Used in emails and UI
   - Human-readable format

3. **Set appropriate `serviceRadiusMiles`**
   - Based on your service area
   - Consider travel time and costs

4. **Use `distanceOrigin` for flexibility**
   - Can be different from office (e.g., warehouse, depot)
   - Useful if you calculate from a different location

## Testing

Location configuration can be tested using:
- Test Dashboard: `/test-dashboard.html`
- Quote Preview: `/quote-preview.html`
- API endpoints: `/api/estimate`, `/api/email/quote/preview`

## Troubleshooting

### Issue: Travel fees not calculating correctly
- **Check**: Ensure `lat`/`lng` are provided or `distanceOrigin` is valid
- **Check**: Verify `serviceRadiusMiles` is set correctly
- **Check**: Ensure Google Maps API key is configured (if using geocoding)

### Issue: Location not found
- **Check**: Verify coordinates are valid (lat: -90 to 90, lng: -180 to 180)
- **Check**: Ensure address format is correct if using address string
- **Check**: API key is valid and has geocoding permissions

### Issue: Changes not taking effect
- **Solution**: Restart the server after configuration changes
- **Note**: Configuration is cached for performance

## Related Documentation

- [Business Configuration Guide](./BUSINESS_CONFIG.md)
- [Travel Fee Calculation](./TRAVEL_FEE_CALCULATION.md)
- [API Documentation](./API.md)

## Migration from Hardcoded Values

If you're migrating from hardcoded location constants:

1. **Old approach** (deprecated):
   ```go
   // In go/internal/infra/geo/geocoding.go
   const OfficeLat = 38.6255
   const OfficeLng = -90.2456
   const ServiceRadiusMiles = 15.0
   ```

2. **New approach** (recommended):
   ```yaml
   # In config/businesses/{businessId}.yaml
   location:
     lat: 38.6255
     lng: -90.2456
     serviceRadiusMiles: 15.0
   ```

3. **Benefits**:
   - Per-business configuration
   - No code changes needed
   - Easy to update
   - Supports multiple businesses

## Support

For questions or issues:
- Check the test dashboard for configuration info
- Review logs for geocoding errors
- Verify API keys are configured correctly
