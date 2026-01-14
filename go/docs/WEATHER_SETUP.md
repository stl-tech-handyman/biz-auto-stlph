# Weather Integration Setup Guide

## Overview
Weather forecasts are automatically included in quote emails for events that are **less than 10 days away**. The system provides weather conditions and smart recommendations, especially for outdoor events.

## API Setup: OpenWeatherMap (Recommended)

**Why OpenWeatherMap?**
- ✅ **Free tier**: 1,000 API calls/day (plenty for most use cases)
- ✅ **No credit card required** for free tier
- ✅ **Easy setup**: Just need an API key
- ✅ **Reliable**: Industry-standard weather API
- ✅ **Good documentation**: Well-documented API

### Step 1: Get Your API Key

1. Go to [OpenWeatherMap Sign Up](https://home.openweathermap.org/users/sign_up)
2. Create a free account (no credit card needed)
3. Once logged in, go to [API Keys](https://home.openweathermap.org/api_keys)
4. Copy your API key (it may take a few minutes to activate)

### Step 2: Set Environment Variable

Add the API key to your environment:

**Windows (PowerShell):**
```powershell
$env:OPENWEATHERMAP_API_KEY="your-api-key-here"
```

**Windows (Command Prompt):**
```cmd
set OPENWEATHERMAP_API_KEY=your-api-key-here
```

**Linux/Mac:**
```bash
export OPENWEATHERMAP_API_KEY="your-api-key-here"
```

**For Production (Secret Manager):**
Store the key in your secret manager and reference it:
```bash
OPENWEATHERMAP_API_KEY=$(gcloud secrets versions access latest --secret="openweathermap-api-key")
```

### Step 3: Verify Setup

The server will log on startup:
- ✅ `Weather service initialized` - Success!
- ⚠️ `Weather service not available` - Check your API key

## Alternative: WeatherAPI.com

If you prefer WeatherAPI.com (1 million calls/month free tier):

1. Sign up at [WeatherAPI.com](https://www.weatherapi.com/signup.aspx)
2. Get your API key
3. Update `go/internal/infra/weather/weather.go` to use WeatherAPI endpoints
4. Set `WEATHERAPI_KEY` environment variable

## How It Works

1. **Event Date Check**: Only fetches weather for events < 10 days away
2. **Geocoding**: Uses existing Google Maps geocoding to get lat/lng from address
3. **Weather Fetch**: Calls OpenWeatherMap API for forecast
4. **Smart Recommendations**: 
   - Temperature warnings (too hot/cold)
   - Precipitation alerts (rain/snow)
   - Wind warnings
   - Outdoor event specific recommendations

## Weather Display in Email

For events < 10 days away, the email will show:

```
🌤️ Weather Forecast: Clear, clear sky (72°F)
Weather considerations: Perfect weather conditions expected
```

For outdoor events with rain:
```
🌤️ Weather Forecast: Rain, light rain (65°F)
⚠️ Rain expected — consider backup indoor space or tent coverage
```

## Cost Considerations

- **OpenWeatherMap Free Tier**: 1,000 calls/day
- **Typical Usage**: ~50-100 quotes/day = well within free tier
- **Upgrade**: If you exceed, paid plans start at $40/month for 100,000 calls/day

## Troubleshooting

**Weather not showing in emails?**
1. Check API key is set: `echo $OPENWEATHERMAP_API_KEY`
2. Check server logs for "Weather service initialized"
3. Verify event is < 10 days away
4. Verify address can be geocoded (check geocoding logs)

**API errors?**
1. Verify API key is correct
2. Check API key is activated (may take 10-60 minutes after signup)
3. Check rate limits (free tier: 60 calls/minute)
4. Verify address format is correct

## Testing

Test weather integration:
1. Create a quote for an event 5-9 days away
2. Include a valid address
3. Check email preview - weather should appear after "Event Details" section

## Future Enhancements

Potential improvements:
- Cache weather data (reduce API calls)
- More sophisticated outdoor event detection
- Hourly forecasts for same-day events
- Historical weather data for planning
