# Setup Google Maps API project: bizops360-maps
# This project is dedicated to Google Maps API usage for event location address boxes

$ErrorActionPreference = "Stop"

# Configuration
$PROJECT_ID = "bizops360-maps"
$PROJECT_NAME = "BizOps360 Maps API"
$REGION = "us-central1"

Write-Host "=== Setting up Google Maps API Project ===" -ForegroundColor Cyan
Write-Host "Project ID: $PROJECT_ID"
Write-Host "Project Name: $PROJECT_NAME"
Write-Host ""

# Check if gcloud is installed
if (-not (Get-Command gcloud -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: gcloud CLI is not installed" -ForegroundColor Red
    exit 1
}

# Get billing account
Write-Host "To find your billing account, run: gcloud beta billing accounts list" -ForegroundColor Yellow
$BILLING_ACCOUNT = Read-Host "Enter billing account ID (e.g., 01ABCD-23EFGH-456789)"

if ([string]::IsNullOrWhiteSpace($BILLING_ACCOUNT)) {
    Write-Host "ERROR: Billing account is required" -ForegroundColor Red
    exit 1
}

# Check if project already exists
Write-Host "Checking if project exists..." -ForegroundColor Cyan
$projectExists = $false
try {
    gcloud projects describe $PROJECT_ID 2>&1 | Out-Null
    if ($LASTEXITCODE -eq 0) {
        $projectExists = $true
        Write-Host "Project $PROJECT_ID already exists, using existing project" -ForegroundColor Yellow
    }
} catch {
    $projectExists = $false
}

if (-not $projectExists) {
    # Create project
    Write-Host "Creating project $PROJECT_ID..." -ForegroundColor Cyan
    gcloud projects create $PROJECT_ID --name="$PROJECT_NAME" --set-as-default
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Failed to create project $PROJECT_ID" -ForegroundColor Red
        exit 1
    }
    Write-Host "Project $PROJECT_ID created" -ForegroundColor Green
}

# Link billing account
Write-Host "Linking billing account..." -ForegroundColor Cyan
gcloud billing projects link $PROJECT_ID --billing-account=$BILLING_ACCOUNT 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Host "WARNING: Failed to link billing account (may already be linked)" -ForegroundColor Yellow
}

# Set as current project
gcloud config set project $PROJECT_ID

# Enable required APIs
Write-Host "Enabling Google Maps APIs..." -ForegroundColor Cyan
$MAPS_APIS = @(
    "maps-javascript-api.googleapis.com",
    "geocoding-api.googleapis.com",
    "places-api.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com"
)

foreach ($api in $MAPS_APIS) {
    Write-Host "  Enabling $api..." -ForegroundColor Cyan
    gcloud services enable $api --project=$PROJECT_ID 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  WARNING: Failed to enable $api (may already be enabled)" -ForegroundColor Yellow
    }
}

Write-Host "All APIs enabled" -ForegroundColor Green

# Create API key
Write-Host ""
Write-Host "Creating Google Maps API key..." -ForegroundColor Cyan
Write-Host "NOTE: API key creation via CLI may not work. You may need to create it manually." -ForegroundColor Yellow
Write-Host ""

# Try to create API key
$API_KEY = $null
try {
    $apiKeyOutput = gcloud alpha services api-keys create --display-name="Maps API Key for Event Location" --project=$PROJECT_ID 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        # Try to extract or retrieve the key
        $apiKeyId = gcloud alpha services api-keys list --project=$PROJECT_ID --format="value(name)" --filter="displayName:'Maps API Key for Event Location'" 2>&1 | Select-Object -First 1
        
        if ($apiKeyId) {
            $API_KEY = gcloud alpha services api-keys get-key-string $apiKeyId --project=$PROJECT_ID --format="value(keyString)" 2>&1
            if ($LASTEXITCODE -eq 0 -and $API_KEY) {
                Write-Host "API key created: $($API_KEY.Substring(0, [Math]::Min(20, $API_KEY.Length)))..." -ForegroundColor Green
            }
        }
    }
} catch {
    Write-Host "Could not create API key via CLI" -ForegroundColor Yellow
}

# If API key creation failed, prompt for manual entry
if ([string]::IsNullOrWhiteSpace($API_KEY)) {
    Write-Host ""
    Write-Host "Please create the API key manually:" -ForegroundColor Yellow
    Write-Host "  1. Go to: https://console.cloud.google.com/apis/credentials?project=$PROJECT_ID" -ForegroundColor Cyan
    Write-Host "  2. Click 'Create Credentials' -> 'API Key'" -ForegroundColor Cyan
    Write-Host "  3. Copy the API key" -ForegroundColor Cyan
    Write-Host ""
    $API_KEY = Read-Host "Enter your Google Maps API key (or press Enter to skip)"
}

# Save API key to Secret Manager if we have one
if (-not [string]::IsNullOrWhiteSpace($API_KEY)) {
    Write-Host "Saving API key to Secret Manager..." -ForegroundColor Cyan
    
    $SECRET_NAME = "maps-api-key"
    
    # Check if secret exists
    $secretExists = $false
    try {
        gcloud secrets describe $SECRET_NAME --project=$PROJECT_ID 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            $secretExists = $true
        }
    } catch {
        $secretExists = $false
    }
    
    if (-not $secretExists) {
        $API_KEY | gcloud secrets create $SECRET_NAME --data-file=- --replication-policy="automatic" --project=$PROJECT_ID
        Write-Host "Secret $SECRET_NAME created in Secret Manager" -ForegroundColor Green
    } else {
        $API_KEY | gcloud secrets versions add $SECRET_NAME --data-file=- --project=$PROJECT_ID
        Write-Host "New version added to secret $SECRET_NAME" -ForegroundColor Green
    }
}

# Display summary
Write-Host ""
Write-Host "=== Setup Complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "Project Details:" -ForegroundColor Cyan
Write-Host "  Project ID: $PROJECT_ID"
Write-Host "  Project Name: $PROJECT_NAME"
Write-Host "  Region: $REGION"
Write-Host ""

if (-not [string]::IsNullOrWhiteSpace($API_KEY)) {
    Write-Host "API Key Information:" -ForegroundColor Cyan
    Write-Host "  API Key: $API_KEY"
    Write-Host "  Secret Name: $SECRET_NAME"
    Write-Host ""
    Write-Host "IMPORTANT: Restrict your API key in Google Cloud Console:" -ForegroundColor Yellow
    Write-Host "  https://console.cloud.google.com/apis/credentials?project=$PROJECT_ID"
    Write-Host ""
    Write-Host "Recommended restrictions:" -ForegroundColor Cyan
    Write-Host "  1. Application restrictions: HTTP referrers (web sites)"
    Write-Host "     Add your website domains (e.g., *.yourdomain.com/*)"
    Write-Host "  2. API restrictions: Restrict to:"
    Write-Host "     - Maps JavaScript API"
    Write-Host "     - Geocoding API"
    Write-Host "     - Places API"
    Write-Host ""
}

Write-Host "Next Steps:" -ForegroundColor Cyan
Write-Host "  1. Restrict your API key (see above)"
Write-Host "  2. Add the API key to your website form"
Write-Host "  3. Test the address autocomplete functionality"
Write-Host ""

Write-Host "To retrieve the API key later:" -ForegroundColor Cyan
Write-Host "  gcloud secrets versions access latest --secret=`"$SECRET_NAME`" --project=`"$PROJECT_ID`""
Write-Host ""

Write-Host "To view API usage and billing:" -ForegroundColor Cyan
Write-Host "  https://console.cloud.google.com/apis/dashboard?project=$PROJECT_ID"
Write-Host ""

