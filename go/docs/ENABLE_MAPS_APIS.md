# Enable Google Maps APIs - Quick Guide

## APIs You Need to Enable

For your event location address box with autocomplete, enable these APIs in the **bizops360-maps** project:

### Required APIs:

1. **Places API (New)** ✅ (Modern replacement for legacy Places API)
   - Service: `places-backend.googleapis.com`
   - Link: https://console.cloud.google.com/apis/library/places-backend.googleapis.com?project=bizops360-maps

2. **Maps JavaScript API** ✅ (For address autocomplete widget)
   - Service: `maps-javascript-api.googleapis.com`
   - Link: https://console.cloud.google.com/apis/library/maps-javascript-api.googleapis.com?project=bizops360-maps

3. **Geocoding API** ✅ (Convert addresses to coordinates)
   - Service: `geocoding-api.googleapis.com`
   - Link: https://console.cloud.google.com/apis/library/geocoding-api.googleapis.com?project=bizops360-maps

4. **Distance Matrix API** ✅ (Get actual driving distance between locations)
   - Service: `distance-matrix-api.googleapis.com`
   - Link: https://console.cloud.google.com/apis/library/distance-matrix-api.googleapis.com?project=bizops360-maps
   - **Note:** This is used for calculating travel fees based on actual driving distance instead of straight-line distance

### Optional (if needed):

5. **Routes API** (For detailed routing/directions - alternative to Distance Matrix)
   - Service: `routes.googleapis.com`
   - Link: https://console.cloud.google.com/apis/library/routes.googleapis.com?project=bizops360-maps

## Quick Enable All APIs

**Direct link to enable all at once:**
https://console.cloud.google.com/apis/library?project=bizops360-maps&q=maps

## Steps:

1. Click each link above, or go to: https://console.cloud.google.com/apis/library?project=bizops360-maps
2. Search for each API name
3. Click **"ENABLE"** for each one
4. Wait for each to finish enabling

## After Enabling:

Make sure your API key restrictions include:
- **Places API (New)** - `places-backend.googleapis.com`
- **Maps JavaScript API** - `maps-javascript-api.googleapis.com`
- **Geocoding API** - `geocoding-api.googleapis.com`
- **Distance Matrix API** - `distance-matrix-api.googleapis.com`

## Check What's Enabled

```bash
gcloud services list --enabled --project=bizops360-maps | grep -i "maps\|places\|geocod\|distance"
```

## Distance Matrix API

For **actual driving distance** calculations (instead of straight-line distance), see:
- [DISTANCE_MATRIX_API_SETUP.md](./DISTANCE_MATRIX_API_SETUP.md) - Complete setup guide

The Distance Matrix API is now integrated and will automatically be used when available, with automatic fallback to Haversine formula if the API is unavailable.
