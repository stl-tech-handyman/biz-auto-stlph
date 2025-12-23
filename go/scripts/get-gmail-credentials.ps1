# Get Gmail credentials from GCP Secret Manager for local development

$ErrorActionPreference = "Stop"

$EMAIL_PROJECT = "bizops360-email-dev"
$SECRET_NAME = "gmail-credentials-json"

Write-Host "[INFO] Fetching Gmail credentials from Secret Manager..." -ForegroundColor Blue
Write-Host "[INFO] Project: $EMAIL_PROJECT" -ForegroundColor Blue
Write-Host "[INFO] Secret: $SECRET_NAME" -ForegroundColor Blue

# Check if secret exists
$secretCheck = gcloud secrets describe $SECRET_NAME --project=$EMAIL_PROJECT 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Secret '$SECRET_NAME' not found in project '$EMAIL_PROJECT'" -ForegroundColor Red
    Write-Host "[INFO] Make sure the secret exists and you have access to it" -ForegroundColor Yellow
    exit 1
}

# Get secret value
Write-Host "[INFO] Retrieving secret value..." -ForegroundColor Blue
$credentials = gcloud secrets versions access latest --secret=$SECRET_NAME --project=$EMAIL_PROJECT 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Failed to retrieve secret value" -ForegroundColor Red
    Write-Host $credentials
    exit 1
}

# Save to temporary file
$tempFile = "$env:TEMP\gmail-credentials-$(Get-Date -Format 'yyyyMMdd-HHmmss').json"
$credentials | Out-File -FilePath $tempFile -Encoding UTF8 -NoNewline

Write-Host "[SUCCESS] Credentials saved to: $tempFile" -ForegroundColor Green
Write-Host ""
Write-Host "To use these credentials, set the environment variable:" -ForegroundColor Yellow
Write-Host "  `$env:GMAIL_CREDENTIALS_JSON = Get-Content '$tempFile' -Raw" -ForegroundColor Cyan
Write-Host ""
Write-Host "Or add to your .env file:" -ForegroundColor Yellow
Write-Host "  GMAIL_CREDENTIALS_JSON=$tempFile" -ForegroundColor Cyan
Write-Host ""
Write-Host "Also set GMAIL_FROM (the email address to send from):" -ForegroundColor Yellow
Write-Host "  `$env:GMAIL_FROM = 'team@stlpartyhelpers.com'" -ForegroundColor Cyan
Write-Host ""

# Optionally set environment variable for current session
$env:GMAIL_CREDENTIALS_JSON = Get-Content $tempFile -Raw
$env:GMAIL_FROM = "team@stlpartyhelpers.com"

Write-Host "[INFO] Environment variables set for current PowerShell session" -ForegroundColor Green
Write-Host "[INFO] GMAIL_CREDENTIALS_JSON: Set (from file)" -ForegroundColor Gray
Write-Host "[INFO] GMAIL_FROM: $env:GMAIL_FROM" -ForegroundColor Gray
Write-Host ""
Write-Host "Note: These environment variables are only set for this PowerShell session." -ForegroundColor Yellow
Write-Host "To make them permanent, add them to your .env file or set them in your startup script." -ForegroundColor Yellow

