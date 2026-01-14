# Google Distance Matrix API Setup Guide

## Overview

The Distance Matrix API is now integrated to calculate **actual driving distances** instead of straight-line distances (Haversine formula). This provides more accurate travel fee calculations.

## What You Need to Do

### Step 1: Get/Verify Your Google Cloud Project

Make sure you have access to the **bizops360-maps** Google Cloud project:
- Project ID: `bizops360-maps`
- If you don't have access, contact your administrator

### Step 2: Enable Distance Matrix API

1. Go to Google Cloud Console: https://console.cloud.google.com/apis/library?project=bizops360-maps
2. Search for **"Distance Matrix API"**
3. Click on **Distance Matrix API**
4. Click **"ENABLE"** button
5. Wait for the API to be enabled (usually takes a few seconds)

**Direct link:** https://console.cloud.google.com/apis/library/distance-matrix-api.googleapis.com?project=bizops360-maps

**Alternative:** Use the quick enable link:
https://console.cloud.google.com/apis/library/distance-matrix-api.googleapis.com?project=bizops360-maps&enableapi=true

### Step 3: Get Your API Key

If you don't have an API key yet:

1. Go to: https://console.cloud.google.com/apis/credentials?project=bizops360-maps
2. Click **"+ CREATE CREDENTIALS"** → **"API key"**
3. Copy the API key (you'll restrict it in the next step)

If you already have an API key, skip to Step 4.

### Step 4: Update API Key Restrictions

Make sure your API key has **Distance Matrix API** enabled:

1. Go to: https://console.cloud.google.com/apis/credentials?project=bizops360-maps
2. Find your API key (the one used in `GOOGLE_MAPS_API_KEY` or `MAPS_API_KEY`)
3. Click on it to edit
4. Under **"API restrictions"**, select **"Restrict key"**
5. Under **"Select APIs"**, search for and select:
   - ✅ **Distance Matrix API** (`distance-matrix-api.googleapis.com`)
   - ✅ **Geocoding API** (`geocoding-api.googleapis.com`) - if not already added
   - ✅ **Places API (New)** (`places-backend.googleapis.com`) - if not already added
   - ✅ **Maps JavaScript API** (`maps-javascript-api.googleapis.com`) - if not already added
6. Click **"SAVE"**

**Important:** The service name for Distance Matrix API is: `distance-matrix-api.googleapis.com`

### Step 5: Set API Key Environment Variable

Set your API key as an environment variable. The service uses the same API key as Geocoding API:

**Option 1: Environment Variable (Local Development)**
```bash
export GOOGLE_MAPS_API_KEY="your-api-key-here"
# OR
export MAPS_API_KEY="your-api-key-here"
```

**Option 2: Google Cloud Secret Manager (Production)**
1. Store the API key in Secret Manager:
   ```bash
   gcloud secrets create google-maps-api-key --data-file=- --project=bizops360-maps
   ```
2. Reference it in your deployment config

**Option 3: .env file (Local Development)**
Create a `.env` file:
```
GOOGLE_MAPS_API_KEY=your-api-key-here
```

**Priority:** `GOOGLE_MAPS_API_KEY` (primary) > `MAPS_API_KEY` (fallback)

### Step 6: Verify Setup

Test that everything is working:

1. **Check API is enabled:**
   ```bash
   gcloud services list --enabled --project=bizops360-maps | grep distance
   ```
   Should show: `distance-matrix-api.googleapis.com`

2. **Test the API key:**
   ```bash
   curl "https://maps.googleapis.com/maps/api/distancematrix/json?origins=4220+Duncan+Ave+St+Louis+MO&destinations=Gateway+Arch+St+Louis+MO&key=YOUR_API_KEY"
   ```
   Should return JSON with distance information (not an error)

3. **Check server logs** when sending a quote - you should see:
   - `"Using Distance Matrix API for distance calculation"` (success)

## How It Works

### Current Implementation

1. **Primary Method:** Distance Matrix API (driving distance)
   - Gets actual driving distance from office to event location
   - More accurate for travel fee calculations
   - Includes duration information (for future use)

2. **Fallback Method:** Haversine formula (straight-line distance)
   - Used if Distance Matrix API is unavailable or fails
   - Ensures the system continues to work even if API has issues
   - Less accurate but still functional

### Code Flow

```
Event Location Address
    ↓
Try Distance Matrix API
    ↓ (if fails)
Fallback to Geocoding + Haversine
    ↓
Calculate Travel Fee
```

## Testing

After enabling the API, test it by:

1. Sending a quote request with an event location
2. Check server logs for:
   - `"Using Distance Matrix API for distance calculation"` (success)
   - `"Distance Matrix API failed, falling back to Haversine formula"` (fallback)

## Pricing

Distance Matrix API pricing (as of 2024):
- **Free tier:** $200 credit per month (covers ~40,000 requests)
- **After free tier:** $5.00 per 1,000 requests

**Note:** Each quote calculation = 1 API request. With the free tier, you can process ~40,000 quotes per month for free.

## Troubleshooting

### API Not Working

1. **Check API is enabled:**
   ```bash
   gcloud services list --enabled --project=bizops360-maps | grep distance
   ```

2. **Check API key restrictions:**
   - Make sure `distance-matrix-api.googleapis.com` is allowed
   - Make sure the API key is correct in your environment

3. **Check logs:**
   - Look for `"Distance Matrix service not available"` warnings
   - Check for API error messages in logs

### Fallback to Haversine

If you see `"falling back to Haversine formula"` in logs:
- The system is still working, just using less accurate distance
- Check API key and API enablement
- Check network connectivity to Google APIs

## Related APIs

You also need these APIs enabled (should already be enabled):
- ✅ **Geocoding API** - For address to coordinates conversion
- ✅ **Places API (New)** - For address autocomplete
- ✅ **Maps JavaScript API** - For address autocomplete widget

See [ENABLE_MAPS_APIS.md](./ENABLE_MAPS_APIS.md) for full list.
