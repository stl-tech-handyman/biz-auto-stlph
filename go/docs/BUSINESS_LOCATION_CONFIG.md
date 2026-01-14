# Business Location Configuration

## Overview

Each business can configure its location settings for distance calculations. This allows you to:
- Set the business office address (for display)
- Configure a separate origin point for distance calculations (e.g., warehouse vs office)
- Set a custom service radius for travel fee calculations

## Configuration

Add a `location` section to your business YAML configuration file:

```yaml
location:
  # Office address (for display purposes)
  officeAddress: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
  
  # Distance origin - the starting point for all distance calculations
  # This can be different from officeAddress if you want to calculate from
  # a warehouse, depot, or other location. Can be:
  # - An address: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
  # - Coordinates: "38.6255,-90.2456" (lat,lng format)
  # If not specified, defaults to officeAddress
  # NOTE: If lat/lng are provided below, they take precedence over distanceOrigin
  distanceOrigin: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
  
  # Latitude and Longitude of the distance origin (RECOMMENDED for accuracy)
  # If provided, these take precedence over distanceOrigin string
  # These coordinates are used for all distance calculations (travel fees, service area checks)
  # Using coordinates avoids geocoding API calls and ensures accuracy
  lat: 38.6255   # Latitude of your distance origin
  lng: -90.2456  # Longitude of your distance origin
  
  # Service radius in miles - locations within this radius have no travel fee
  # Default: 15.0 miles
  serviceRadiusMiles: 15.0
```

## Fields Explained

### `officeAddress` (Required)

The business office address. This is used for display purposes and as a fallback if `distanceOrigin` is not specified.

**Example:**
```yaml
officeAddress: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
```

### `lat` and `lng` (Recommended)

**The most accurate way to specify your distance origin.** These are separate latitude and longitude fields that take **highest precedence** over all other location settings.

**Why use lat/lng:**
- ✅ **Most accurate** - No geocoding API calls needed
- ✅ **Faster** - Avoids API latency
- ✅ **Cost-effective** - Saves geocoding API quota
- ✅ **Reliable** - No dependency on address parsing

**Priority:** `lat`/`lng` > `distanceOrigin` > `officeAddress`

**Example:**
```yaml
lat: 38.6255   # Latitude
lng: -90.2456  # Longitude
```

**How to find coordinates:**
- Use Google Maps: Right-click on location → "What's here?" → Copy coordinates
- Use geocoding API once to get coordinates, then hardcode them
- Use online tools like latlong.net

### `distanceOrigin` (Optional)

The starting point for all distance calculations. This is a **separate variable** from `officeAddress` to allow flexibility.

**Priority:** Only used if `lat`/`lng` are not provided.

**Use cases:**
- Calculate from a warehouse instead of the office
- Calculate from a central depot
- Use coordinates as a string (less preferred than lat/lng fields)

**Formats:**
- **Address string:** `"4220 Duncan Ave Suite 201, St. Louis, MO 63110"`
- **Coordinates string:** `"38.6255,-90.2456"` (latitude,longitude)

**If not specified:** Defaults to `officeAddress`

**Example - Different warehouse:**
```yaml
officeAddress: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
distanceOrigin: "1234 Warehouse Blvd, St. Louis, MO 63120"  # Calculate from warehouse
```

**Example - Using coordinates string (less preferred):**
```yaml
officeAddress: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
distanceOrigin: "38.6255,-90.2456"  # Specific coordinates
```

### `serviceRadiusMiles` (Optional)

The service area radius in miles. Locations within this radius have **no travel fee**.

**Default:** `15.0` miles

**Example:**
```yaml
serviceRadiusMiles: 20.0  # 20-mile service area
```

## How It Works

1. **Distance Calculation Origin (Priority Order):**
   - **Priority 1:** `lat`/`lng` fields (if provided) - **RECOMMENDED**
   - **Priority 2:** `distanceOrigin` (if provided, parsed as address or coordinates string)
   - **Priority 3:** `officeAddress` (fallback)
   - **Priority 4:** Hardcoded defaults (if location section missing)
   
   The origin is used for both Distance Matrix API (driving distance) and Haversine formula (straight-line).

2. **Service Area Check:**
   - Uses `serviceRadiusMiles` to determine if a location is within the service area
   - Locations within the radius: **$0 travel fee**
   - Locations outside the radius: Travel fee calculated based on distance

3. **Fallback Behavior:**
   - If `location` section is missing, system uses hardcoded defaults:
     - Address: `"4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"`
     - Coordinates: `38.6255, -90.2456`
     - Service radius: `15.0` miles

## Example Configuration

**Complete example for STL Party Helpers (Recommended):**

```yaml
id: stlpartyhelpers
displayName: "STL Party Helpers"
timezone: "America/Chicago"
currency: "usd"

location:
  officeAddress: "4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"
  distanceOrigin: "4220 Duncan Ave., Ste. 201, St. Louis, MO 63110"
  lat: 38.6255   # Recommended: Use explicit coordinates
  lng: -90.2456  # Recommended: Use explicit coordinates
  serviceRadiusMiles: 15.0

# ... rest of config
```

**Example with warehouse:**

```yaml
id: mybusiness
displayName: "My Business"
timezone: "America/Chicago"
currency: "usd"

location:
  officeAddress: "123 Main St, City, State 12345"  # Office for display
  distanceOrigin: "456 Warehouse Rd, City, State 12345"  # Calculate from warehouse
  serviceRadiusMiles: 20.0  # 20-mile service area

# ... rest of config
```

## Migration

If you're upgrading from the old hardcoded location:

1. Add the `location` section to your business YAML file
2. Set `officeAddress` to your actual office address
3. Set `distanceOrigin` to the same address (or different if needed)
4. Set `serviceRadiusMiles` to your desired service radius (default: 15.0)

The system will automatically use these values instead of the hardcoded defaults.

## Technical Details

- Location info is loaded from business config when available
- Falls back to constants if business config is not available
- Distance calculations use the configured origin point
- Service radius is used in travel fee calculations
- All location-related functions support both address strings and coordinates

## Related Documentation

- [Distance Matrix API Setup](./DISTANCE_MATRIX_API_SETUP.md) - For actual driving distance calculations
- [Enable Maps APIs](./ENABLE_MAPS_APIS.md) - For Google Maps API setup
