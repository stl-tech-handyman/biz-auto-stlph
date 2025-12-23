# Create Google Maps API Key - Quick Guide

## Step 1: Link Billing Account (if not done)

1. Go to: https://console.cloud.google.com/billing?project=bizops360-maps
2. Select your billing account and link it to the project

## Step 2: Enable Required APIs

1. Go to: https://console.cloud.google.com/apis/library?project=bizops360-maps
2. Enable these APIs:
   - **Maps JavaScript API** - Search for "Maps JavaScript API" and click Enable
   - **Geocoding API** - Search for "Geocoding API" and click Enable  
   - **Places API** - Search for "Places API" and click Enable

## Step 3: Create API Key

1. Go directly to create API key: https://console.cloud.google.com/apis/credentials?project=bizops360-maps
2. Click **"+ CREATE CREDENTIALS"** → **"API key"**
3. Copy the API key that appears
4. Click **"RESTRICT KEY"** to secure it:
   - **Application restrictions**: Select "HTTP referrers (web sites)"
   - Add your website domains:
     - `https://yourdomain.com/*`
     - `https://*.yourdomain.com/*`
     - `http://localhost:*` (for local development)
   - **API restrictions**: Select "Restrict key"
   - Choose only:
     - Maps JavaScript API
     - Geocoding API
     - Places API
   - Click **"SAVE"**

## Step 4: Save to Secret Manager

After you have the API key, run:

```bash
# Replace YOUR_API_KEY with the actual key
echo -n "YOUR_API_KEY" | gcloud secrets create maps-api-key --data-file=- --replication-policy="automatic" --project=bizops360-maps
```

Or if the secret already exists:

```bash
echo -n "YOUR_API_KEY" | gcloud secrets versions add maps-api-key --data-file=- --project=bizops360-maps
```

## Quick Links

- **Create API Key**: https://console.cloud.google.com/apis/credentials?project=bizops360-maps
- **Enable APIs**: https://console.cloud.google.com/apis/library?project=bizops360-maps
- **Link Billing**: https://console.cloud.google.com/billing?project=bizops360-maps
- **View API Usage**: https://console.cloud.google.com/apis/dashboard?project=bizops360-maps

