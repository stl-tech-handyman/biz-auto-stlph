# Deploy Go API to Development Environment (PowerShell)

$ErrorActionPreference = "Stop"

# Configuration
$PROJECT_ID = "bizops360-dev"
$REGION = "us-central1"
$SERVICE_NAME = "bizops360-api-go-dev"
$IMAGE_NAME = "gcr.io/$PROJECT_ID/$SERVICE_NAME"

Write-Host "[INFO] Deploying Go API to DEV environment..." -ForegroundColor Blue
Write-Host "[INFO] Project: $PROJECT_ID" -ForegroundColor Blue
Write-Host "[INFO] Service: $SERVICE_NAME" -ForegroundColor Blue
Write-Host "[INFO] Region: $REGION" -ForegroundColor Blue

# Set project
gcloud config set project $PROJECT_ID

# Build Docker image
Write-Host "[INFO] Building Docker image..." -ForegroundColor Blue
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent (Split-Path -Parent $scriptDir)
Set-Location $projectRoot

docker build -f go/Dockerfile.dev -t "$IMAGE_NAME`:latest" .

# Push to Artifact Registry
Write-Host "[INFO] Pushing image to Artifact Registry..." -ForegroundColor Blue
docker push "$IMAGE_NAME`:latest"

# Prepare secrets
$SECRET_ARGS = "SERVICE_API_KEY=svc-api-key-dev:latest"

$testSecret = gcloud secrets describe stripe-secret-key-test --project=$PROJECT_ID 2>&1
if ($LASTEXITCODE -eq 0) {
    $SECRET_ARGS += ",STRIPE_SECRET_KEY_TEST=stripe-secret-key-test:latest"
} else {
    Write-Host "[ERROR] Secret stripe-secret-key-test not found in $PROJECT_ID" -ForegroundColor Red
    exit 1
}

$prodSecret = gcloud secrets describe stripe-secret-key-prod --project=$PROJECT_ID 2>&1
if ($LASTEXITCODE -eq 0) {
    $SECRET_ARGS += ",STRIPE_SECRET_KEY_PROD=stripe-secret-key-prod:latest"
}

# Deploy to Cloud Run
Write-Host "[INFO] Deploying to Cloud Run..." -ForegroundColor Blue
gcloud run deploy $SERVICE_NAME `
    --project=$PROJECT_ID `
    --region=$REGION `
    --image="$IMAGE_NAME`:latest" `
    --platform=managed `
    --allow-unauthenticated `
    --port=8080 `
    --memory=512Mi `
    --cpu=1 `
    --min-instances=0 `
    --max-instances=5 `
    --concurrency=80 `
    --timeout=300 `
    --set-secrets=$SECRET_ARGS `
    --set-env-vars="ENV=dev,LOG_LEVEL=debug,CONFIG_DIR=/app/config,TEMPLATES_DIR=/app/templates" `
    --labels="env=dev,service=api-go,type=cloud-run"

# Get service URL
$SERVICE_URL = gcloud run services describe $SERVICE_NAME `
    --project=$PROJECT_ID `
    --region=$REGION `
    --format="value(status.url)"

Write-Host "[SUCCESS] Deployment complete!" -ForegroundColor Green
Write-Host "[INFO] Service URL: $SERVICE_URL" -ForegroundColor Blue
Write-Host "[INFO] Health check: $SERVICE_URL/api/health" -ForegroundColor Blue
Write-Host "[INFO] Stripe endpoint: $SERVICE_URL/api/stripe/deposit/calculate" -ForegroundColor Blue

