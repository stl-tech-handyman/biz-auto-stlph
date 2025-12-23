# Quick Stripe Integration Test

$ErrorActionPreference = "Continue"

Write-Host "=== Testing Stripe Integration ===" -ForegroundColor Cyan
Write-Host ""

# Get service URL
$SERVICE_URL = "https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app"

# Get API key
Write-Host "[INFO] Retrieving API key from Secret Manager..." -ForegroundColor Blue
$API_KEY = gcloud secrets versions access latest --secret="svc-api-key-dev" --project="bizops360-dev" 2>&1

if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrEmpty($API_KEY)) {
    Write-Host "[ERROR] Could not retrieve API key" -ForegroundColor Red
    exit 1
}

Write-Host "[SUCCESS] API key retrieved" -ForegroundColor Green
Write-Host ""

# Test 1: Health check (no auth)
Write-Host "[TEST 1] Health Check (no auth required)..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$SERVICE_URL/api/health" -Method GET -ErrorAction Stop
    Write-Host "[SUCCESS] Health check passed" -ForegroundColor Green
    Write-Host "Status: $($healthResponse.status)" -ForegroundColor Gray
} catch {
    Write-Host "[ERROR] Health check failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 2: Stripe test endpoint
Write-Host "[TEST 2] Stripe Integration Test..." -ForegroundColor Yellow
try {
    $headers = @{
        "X-Api-Key" = $API_KEY
        "Content-Type" = "application/json"
    }
    $body = @{} | ConvertTo-Json
    $stripeResponse = Invoke-RestMethod -Uri "$SERVICE_URL/api/stripe/test" -Method POST -Headers $headers -Body $body -ErrorAction Stop
    Write-Host "[SUCCESS] Stripe test passed!" -ForegroundColor Green
    Write-Host "Response: $($stripeResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Gray
} catch {
    Write-Host "[ERROR] Stripe test failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody" -ForegroundColor Red
    }
}
Write-Host ""

# Test 3: Deposit calculation
Write-Host "[TEST 3] Deposit Calculation..." -ForegroundColor Yellow
try {
    $headers = @{
        "X-Api-Key" = $API_KEY
    }
    $depositResponse = Invoke-RestMethod -Uri "$SERVICE_URL/api/stripe/deposit/calculate?estimate=1000" -Method GET -Headers $headers -ErrorAction Stop
    Write-Host "[SUCCESS] Deposit calculation passed!" -ForegroundColor Green
    Write-Host "Deposit amount: $($depositResponse.deposit)" -ForegroundColor Gray
    Write-Host "Response: $($depositResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Gray
} catch {
    Write-Host "[ERROR] Deposit calculation failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 4: Estimate calculation
Write-Host "[TEST 4] Estimate Calculation..." -ForegroundColor Yellow
try {
    $headers = @{
        "X-Api-Key" = $API_KEY
        "Content-Type" = "application/json"
    }
    $body = @{
        eventDate = "2025-12-25"
        durationHours = 4
        numHelpers = 2
    } | ConvertTo-Json
    $estimateResponse = Invoke-RestMethod -Uri "$SERVICE_URL/api/estimate" -Method POST -Headers $headers -Body $body -ErrorAction Stop
    Write-Host "[SUCCESS] Estimate calculation passed!" -ForegroundColor Green
    Write-Host "Estimated total: $($estimateResponse.estimatedTotal)" -ForegroundColor Gray
} catch {
    Write-Host "[ERROR] Estimate calculation failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Testing Complete ===" -ForegroundColor Cyan

