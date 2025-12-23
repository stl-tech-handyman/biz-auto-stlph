# Google Maps API Setup Guide

This guide explains how to set up a dedicated Google Cloud project for Google Maps API usage in your event location address box.

## Why a Separate Project?

Creating a separate project (`bizops360-maps`) for Google Maps API provides:

- **Better Cost Tracking**: Isolate Maps API costs from other services
- **Security**: Separate API keys and access controls
- **Billing Clarity**: Easier to track and budget for Maps API usage
- **Independent Management**: Update Maps API settings without affecting main project

## Quick Setup

### 1. Create the Project and API Key

Run the setup script:

**On Linux/Mac:**
```bash
cd go
bash scripts/setup-maps-project.sh
```

**On Windows (PowerShell):**
```powershell
cd go
.\scripts\setup-maps-project.ps1
```

This script will:
- Create a new Google Cloud project: `bizops360-maps`
- Enable required Google Maps APIs:
  - Maps JavaScript API (for address autocomplete)
  - Geocoding API (for address to coordinates conversion)
  - Places API (for place autocomplete)
- Create an API key (or prompt you to create one manually)
- Store the API key in Secret Manager

**Note**: If the script cannot create the API key automatically, it will provide a link to create it manually in the Google Cloud Console.

### 2. Restrict Your API Key

**IMPORTANT**: After creating the API key, you must restrict it for security.

1. Go to [Google Cloud Console - Credentials](https://console.cloud.google.com/apis/credentials?project=bizops360-maps)
2. Click on your API key
3. Under **Application restrictions**, select **HTTP referrers (web sites)**
4. Add your website domains:
   - `https://yourdomain.com/*`
   - `https://*.yourdomain.com/*`
   - `http://localhost:*` (for local development)
5. Under **API restrictions**, select **Restrict key**
6. Select only:
   - Maps JavaScript API
   - Geocoding API
   - Places API
7. Click **Save**

### 3. Get Your API Key

To retrieve the API key:

**On Linux/Mac:**
```bash
cd go
bash scripts/get-maps-api-key.sh
```

**On Windows (PowerShell):**
```powershell
cd go
gcloud secrets versions access latest --secret="maps-api-key" --project="bizops360-maps"
```

Or manually:

```bash
gcloud secrets versions access latest --secret="maps-api-key" --project="bizops360-maps"
```

### 4. Add to Your Website Form

Add the API key to your website's Google Maps integration:

```html
<!-- Example: Google Maps JavaScript API -->
<script src="https://maps.googleapis.com/maps/api/js?key=YOUR_API_KEY&libraries=places"></script>
```

Or in your JavaScript configuration:

```javascript
const GOOGLE_MAPS_API_KEY = 'YOUR_API_KEY';
```

## Project Details

- **Project ID**: `bizops360-maps`
- **Project Name**: BizOps360 Maps API
- **Region**: us-central1
- **Secret Name**: `maps-api-key`

## APIs Enabled

The following Google Maps APIs are enabled in this project:

1. **Maps JavaScript API** (`maps-javascript-api.googleapis.com`)
   - Used for interactive maps and address autocomplete
   - Required for the address input box

2. **Geocoding API** (`geocoding-api.googleapis.com`)
   - Converts addresses to coordinates (lat/lng)
   - Used for location-based features

3. **Places API** (`places-api.googleapis.com`)
   - Provides place autocomplete and place details
   - Enhances address input with suggestions

## Monitoring Usage

### View API Usage

```bash
# View API usage dashboard
open https://console.cloud.google.com/apis/dashboard?project=bizops360-maps
```

### Check Billing

```bash
# View billing dashboard
open https://console.cloud.google.com/billing?project=bizops360-maps
```

### View Quotas

```bash
# Check API quotas
gcloud services list --enabled --project=bizops360-maps
```

## Cost Management

Google Maps API has a free tier with monthly credits:

- **Maps JavaScript API**: $200 free credit/month
- **Geocoding API**: $200 free credit/month  
- **Places API**: $200 free credit/month

After free credits, pay-as-you-go pricing applies. Monitor usage in the Google Cloud Console.

### Set Up Billing Alerts

1. Go to [Billing](https://console.cloud.google.com/billing?project=bizops360-maps)
2. Click **Budgets & alerts**
3. Create a budget to get notified when spending exceeds thresholds

## Troubleshooting

### API Key Not Working

1. **Check API restrictions**: Make sure your domain is added to HTTP referrers
2. **Verify APIs enabled**: Ensure Maps JavaScript API, Geocoding API, and Places API are enabled
3. **Check billing**: Verify billing account is linked and active
4. **Review quotas**: Check if you've exceeded API quotas

### Common Errors

- **RefererNotAllowedMapError**: Your domain is not in the allowed referrers list
- **ApiNotActivatedMapError**: The Maps JavaScript API is not enabled
- **OverQueryLimitError**: You've exceeded the API quota

### Get Help

- [Google Maps Platform Documentation](https://developers.google.com/maps/documentation)
- [Google Cloud Support](https://cloud.google.com/support)

## Security Best Practices

1. **Always restrict API keys** by HTTP referrer and API restrictions
2. **Never commit API keys** to version control
3. **Use Secret Manager** to store keys securely
4. **Rotate keys periodically** if compromised
5. **Monitor usage** for unusual patterns

## Next Steps

After setup:

1. ✅ Restrict your API key (see step 2 above)
2. ✅ Add the API key to your website form
3. ✅ Test the address autocomplete functionality
4. ✅ Set up billing alerts
5. ✅ Monitor usage regularly

