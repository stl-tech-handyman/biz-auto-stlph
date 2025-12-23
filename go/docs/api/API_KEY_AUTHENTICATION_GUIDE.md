# API Key Authentication Guide

## Overview

All public-facing API endpoints now require API key authentication via the `X-Api-Key` header.

## Protected Endpoints

The following endpoints require API key authentication:

- **All Stripe endpoints** (`/api/stripe/*`)
  - `POST /api/stripe/deposit`
  - `GET /api/stripe/deposit/calculate`
  - `POST /api/stripe/deposit/with-email`
  - `GET /api/stripe/deposit/amount`

## Public Endpoints (No Auth Required)

These endpoints remain public and do not require authentication:

- `GET /` - Service information
- `GET /api/health` - Health check
- `GET /api/health/ready` - Readiness check
- `GET /api/health/live` - Liveness check
- `GET /api/geocode` - Geocoding service info

## Getting Your API Key

### For Development Environment

```bash
gcloud secrets versions access latest --secret=svc-api-key-dev --project=bizops360-dev
```

### For Production Environment

```bash
gcloud secrets versions access latest --secret=svc-api-key-prod --project=bizops360-prod
```

## Using the API Key

Include the API key in the `X-Api-Key` header for all protected endpoints:

```bash
curl -X POST "https://stlph-api-nhrhozfuaq-uc.a.run.app/api/stripe/deposit" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY_HERE" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "estimatedTotal": 10000
  }'
```

## Error Responses

### 401 Unauthorized - Missing API Key

```json
{
  "error": "Unauthorized",
  "message": "API key is required",
  "hint": "Include X-Api-Key header in your request"
}
```

### 401 Unauthorized - Invalid API Key

```json
{
  "error": "Unauthorized",
  "message": "Invalid API key"
}
```

### 500 Internal Server Error - API Key Not Configured

```json
{
  "error": "Service Configuration Error",
  "message": "API authentication is not properly configured"
}
```

## Security Notes

- API keys are stored securely in Google Cloud Secret Manager
- Keys are automatically injected as environment variables during Cloud Run deployment
- Never commit API keys to version control
- Rotate keys regularly for security best practices
- Invalid authentication attempts are logged for security monitoring

## Implementation Details

The API key authentication is implemented as Express middleware in `services/stlph-api/middleware/auth.js` and is applied to all Stripe routes automatically.

