# Postman Collection Documentation

## Overview

This document describes the Postman collection for the BizOps360 Go API. The collection includes all available endpoints with example requests and curl commands.

## Collection Location

The Postman collection file is located at:
- **File**: `go/postman/BizOps360-Go-API.postman_collection.json`
- **Import**: Import this file into Postman to get all API endpoints ready to test

## Environment Variables

The collection uses Postman environment variables. Create an environment with:

- `base_url_dev`: `https://bizops360-api-go-dev-nhrhozfuaq-uc.a.run.app`
- `base_url_prod`: `https://bizops360-api-go-prod-XXXXX-uc.a.run.app` (update when prod is deployed)
- `api_key`: Your API key from Secret Manager (`svc-api-key-dev` or `svc-api-key-prod`)

## How to Use

1. **Import Collection**: Import `BizOps360-Go-API.postman_collection.json` into Postman
2. **Create Environment**: Create a new environment with the variables above
3. **Select Environment**: Select your environment (dev or prod)
4. **Set API Key**: Update the `api_key` variable with your actual API key
5. **Test Endpoints**: All endpoints are ready to test with example data

## Adding New Endpoints

**IMPORTANT**: When adding a new endpoint to the Go API, you MUST:

1. **Add to Postman Collection**: Add the new endpoint to `BizOps360-Go-API.postman_collection.json`
2. **Include curl Command**: Every request in Postman should have a curl command in the description
3. **Update This Document**: Add the endpoint to the endpoints list below
4. **Test**: Verify the endpoint works in Postman before committing

### Format for curl Commands

```bash
# Description of what this endpoint does
curl -X METHOD "{{base_url}}/endpoint/path" \
  -H "X-Api-Key: {{api_key}}" \
  -H "Content-Type: application/json" \
  -d '{
    "field": "value"
  }'
```

## Endpoints

### Health & Info

#### GET /
Root endpoint - Service information
- **Auth**: None required
- **curl**: See collection

#### GET /api/health
Health check endpoint
- **Auth**: None required
- **curl**: See collection

#### GET /api/health/ready
Readiness check endpoint
- **Auth**: None required
- **curl**: See collection

#### GET /api/health/live
Liveness check endpoint
- **Auth**: None required
- **curl**: See collection

### Stripe Endpoints

#### POST /api/stripe/deposit
Create a deposit payment
- **Auth**: API Key required
- **curl**: See collection

#### GET /api/stripe/deposit/calculate
Calculate deposit amount from estimate
- **Auth**: API Key required
- **Query Params**: `estimate` (dollars), `deposit` (dollars)
- **curl**: See collection

#### POST /api/stripe/deposit/with-email
Create deposit and send email
- **Auth**: API Key required
- **curl**: See collection

#### POST /api/stripe/test
Test Stripe integration
- **Auth**: API Key required
- **curl**: See collection

### Estimate Endpoints

#### POST /api/estimate
Calculate event estimate
- **Auth**: API Key required
- **Body**: `{ "eventDate": "2025-12-25", "durationHours": 4, "numHelpers": 2 }`
- **curl**: See collection

#### GET /api/estimate/special-dates
Get special dates (holidays, etc.)
- **Auth**: API Key required
- **Query Params**: `years` (optional, default: 5)
- **curl**: See collection

### Email Endpoints

#### POST /api/email/test
Send test email
- **Auth**: API Key required
- **Body**: `{ "to": "test@example.com", "subject": "Test", "html": "<p>Test</p>" }`
- **curl**: See collection

#### POST /api/email/booking-deposit
Send booking deposit email
- **Auth**: API Key required
- **Body**: `{ "name": "John Doe", "email": "john@example.com" }`
- **curl**: See collection

### V1 API (Pipeline-based)

#### POST /v1/form-events
Submit form event
- **Auth**: None required
- **curl**: See collection

#### POST /v1/triggers
Trigger pipeline execution
- **Auth**: None required
- **curl**: See collection

## Getting API Keys

### Development
```bash
gcloud secrets versions access latest --secret="svc-api-key-dev" --project="bizops360-dev"
```

### Production
```bash
gcloud secrets versions access latest --secret="svc-api-key-prod" --project="bizops360-prod"
```

## Testing Workflow

1. **Start with Health Check**: Test `/api/health` to verify service is running
2. **Test Public Endpoints**: Test root `/` and health endpoints
3. **Set API Key**: Update environment variable with your API key
4. **Test Authenticated Endpoints**: Test endpoints that require API key
5. **Verify Responses**: Check that responses match expected format

## Troubleshooting

### 401 Unauthorized
- Check that `api_key` environment variable is set correctly
- Verify API key matches the environment (dev vs prod)
- Check that `X-Api-Key` header is being sent

### 400 Bad Request
- Verify request body matches expected format
- Check required fields are present
- Validate data types (dates, numbers, etc.)

### 500 Internal Server Error
- Check Cloud Run logs: `gcloud run services logs read bizops360-api-go-dev --project bizops360-dev`
- Verify secrets are configured correctly
- Check environment variables

## Maintenance

- **Update Collection**: When endpoints change, update the Postman collection
- **Update curl Commands**: Keep curl commands in sync with Postman requests
- **Test Regularly**: Run through collection after deployments
- **Document Changes**: Update this document when adding/modifying endpoints

