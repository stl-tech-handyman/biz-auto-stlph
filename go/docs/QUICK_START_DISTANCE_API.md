# Quick Start: Distance Matrix API Setup

## TL;DR - Get It Working in 5 Minutes

### 1. Enable the API
👉 **Click here:** https://console.cloud.google.com/apis/library/distance-matrix-api.googleapis.com?project=bizops360-maps&enableapi=true

Click **"ENABLE"** and wait ~10 seconds.

### 2. Update API Key Restrictions
👉 **Click here:** https://console.cloud.google.com/apis/credentials?project=bizops360-maps

1. Find your API key (the one in `GOOGLE_MAPS_API_KEY` environment variable)
2. Click on it
3. Under "API restrictions", make sure **Distance Matrix API** is checked
4. Click **"SAVE"**

### 3. Verify It Works
Send a test quote request and check server logs for:
- ✅ `"Using Distance Matrix API for distance calculation"` = Success!
- ⚠️ `"falling back to Haversine formula"` = API not enabled or key issue

## That's It!

The system will now use **actual driving distances** instead of straight-line distances for travel fee calculations.

## Need More Details?

- **Full setup guide:** [DISTANCE_MATRIX_API_SETUP.md](./DISTANCE_MATRIX_API_SETUP.md)
- **Business location config:** [BUSINESS_LOCATION_CONFIG.md](./BUSINESS_LOCATION_CONFIG.md)
- **All Maps APIs:** [ENABLE_MAPS_APIS.md](./ENABLE_MAPS_APIS.md)

## Troubleshooting

**API not working?**
1. Check API is enabled: https://console.cloud.google.com/apis/library?project=bizops360-maps
2. Check API key restrictions include Distance Matrix API
3. Check `GOOGLE_MAPS_API_KEY` environment variable is set

**Still using Haversine?**
- That's OK! The system falls back gracefully. Check API enablement and key restrictions.
