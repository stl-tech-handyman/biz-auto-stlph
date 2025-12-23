# Google Maps API - Quick Start

Get your Google Maps API key in 3 steps:

## Step 1: Create the Project

Run the setup script (you'll need your billing account ID):

**Linux/Mac:**
```bash
cd go
bash scripts/setup-maps-project.sh
```

**Windows:**
```powershell
cd go
.\scripts\setup-maps-project.ps1
```

## Step 2: Restrict Your API Key

1. Go to: https://console.cloud.google.com/apis/credentials?project=bizops360-maps
2. Click on your API key
3. Set **Application restrictions** → HTTP referrers → Add your domains
4. Set **API restrictions** → Restrict key → Select:
   - Maps JavaScript API
   - Geocoding API
   - Places API
5. Save

## Step 3: Get Your API Key

**Linux/Mac:**
```bash
cd go
bash scripts/get-maps-api-key.sh
```

**Windows:**
```powershell
gcloud secrets versions access latest --secret="maps-api-key" --project="bizops360-maps"
```

Copy the key and add it to your website form!

---

**Project**: `bizops360-maps`  
**Full Guide**: See [GOOGLE_MAPS_SETUP.md](./GOOGLE_MAPS_SETUP.md)

