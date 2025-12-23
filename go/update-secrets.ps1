# Update Cloud Run service with new secrets (without rebuilding)

$ErrorActionPreference = "Stop"

$PROJECT_ID = "bizops360-dev"
$SERVICE_NAME = "bizops360-api-go-dev"
$REGION = "us-central1"
$GMAIL_FROM = "team@stlpartyhelpers.com"

Write-Host "[INFO] Updating Cloud Run service with new secrets..." -ForegroundColor Blue
Write-Host "[INFO] Project: $PROJECT_ID" -ForegroundColor Blue
Write-Host "[INFO] Service: $SERVICE_NAME" -ForegroundColor Blue

# Set project
gcloud config set project $PROJECT_ID

# Prepare secrets
$SECRET_ARGS = "SERVICE_API_KEY=svc-api-key-dev:latest,STRIPE_SECRET_KEY_TEST=stripe-secret-key-test:latest,STRIPE_SECRET_KEY_PROD=stripe-secret-key-prod:latest"

# Check if Gmail credentials exist in main project
$gmailCheck = gcloud secrets describe gmail-credentials-json --project=$PROJECT_ID 2>&1
if ($LASTEXITCODE -eq 0) {
    $SECRET_ARGS += ",GMAIL_CREDENTIALS_JSON=gmail-credentials-json:latest"
    Write-Host "[INFO] Gmail credentials found in $PROJECT_ID" -ForegroundColor Green
} else {
    Write-Host "[WARNING] Gmail credentials not found in $PROJECT_ID (email features may not work)" -ForegroundColor Yellow
}

# Update Cloud Run service
Write-Host "[INFO] Updating service..." -ForegroundColor Blue
gcloud run services update $SERVICE_NAME `
    --project=$PROJECT_ID `
    --region=$REGION `
    --update-secrets=$SECRET_ARGS `
    --update-env-vars="ENV=dev,LOG_LEVEL=debug,CONFIG_DIR=/app/config,TEMPLATES_DIR=/app/templates,GMAIL_FROM=$GMAIL_FROM"

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Service updated successfully!" -ForegroundColor Green
    
    $SERVICE_URL = gcloud run services describe $SERVICE_NAME `
        --project=$PROJECT_ID `
        --region=$REGION `
        --format="value(status.url)"
    
    Write-Host "[INFO] Service URL: $SERVICE_URL" -ForegroundColor Blue
    Write-Host "[INFO] Health check: $SERVICE_URL/api/health" -ForegroundColor Blue
} else {
    Write-Host "[ERROR] Failed to update service" -ForegroundColor Red
    exit 1
}

