# Local development startup script for Go API
# WARNING: This uses LIVE Stripe keys - be careful!

$PORT = "8080"
$ALT_PORT = "8081"

# Check if port 8080 is available
$process = Get-NetTCPConnection -LocalPort $PORT -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
if ($process) {
    Write-Host "⚠️  Port $PORT is in use, switching to port $ALT_PORT" -ForegroundColor Yellow
    $PORT = $ALT_PORT
}

# Load from .env file if it exists, otherwise use defaults
# Note: .env file is loaded automatically by the Go code via godotenv
# These are fallback values if .env doesn't exist
if (-not $env:STRIPE_SECRET_KEY_TEST) {
    $env:STRIPE_SECRET_KEY_TEST="sk_test_YOUR_TEST_KEY_HERE"
}
if (-not $env:STRIPE_SECRET_KEY_PROD) {
    $env:STRIPE_SECRET_KEY_PROD="sk_live_YOUR_PROD_KEY_HERE"
}
if (-not $env:SERVICE_API_KEY) {
    $env:SERVICE_API_KEY="test-api-key-12345"
}

# Gmail credentials (optional - required for email features)
# Try to load from .env file first, then try to fetch from Secret Manager
if (-not $env:GMAIL_CREDENTIALS_JSON) {
    # Try to fetch credentials automatically
    $credScript = Join-Path $PSScriptRoot "scripts\get-gmail-credentials.ps1"
    if (Test-Path $credScript) {
        Write-Host "ℹ️  Fetching Gmail credentials from Secret Manager..." -ForegroundColor Cyan
        & $credScript | Out-Null
    }
    
    # If still not set, check for .env file
    if (-not $env:GMAIL_CREDENTIALS_JSON) {
        $envFile = Join-Path $PSScriptRoot ".env"
        if (Test-Path $envFile) {
            $envContent = Get-Content $envFile -Raw
            if ($envContent -match "GMAIL_CREDENTIALS_JSON=(.+)") {
                $credPath = $matches[1].Trim()
                if (Test-Path $credPath) {
                    $env:GMAIL_CREDENTIALS_JSON = Get-Content $credPath -Raw
                    Write-Host "✅ Loaded Gmail credentials from .env file" -ForegroundColor Green
                }
            }
        }
    }
    
    if (-not $env:GMAIL_CREDENTIALS_JSON) {
        Write-Host "⚠️  GMAIL_CREDENTIALS_JSON not set - email features will not work" -ForegroundColor Yellow
        Write-Host "   Run: powershell -ExecutionPolicy Bypass -File scripts/get-gmail-credentials.ps1" -ForegroundColor Gray
    }
}
if (-not $env:GMAIL_FROM) {
    $env:GMAIL_FROM="team@stlpartyhelpers.com"
}

$env:ENV="dev"
$env:PORT=$PORT
$env:LOG_LEVEL="debug"
$env:CONFIG_DIR="./config"
$env:TEMPLATES_DIR="./templates"

Write-Host "🚀 Starting Go API Server..." -ForegroundColor Green
Write-Host ""
Write-Host "Configuration:" -ForegroundColor Yellow
Write-Host "  Stripe: LIVE (Production) Key" -ForegroundColor Red
Write-Host "  API Key: $env:SERVICE_API_KEY" -ForegroundColor Cyan
Write-Host "  Port: $env:PORT" -ForegroundColor Cyan
Write-Host "  Environment: $env:ENV" -ForegroundColor Cyan
if ($env:GMAIL_CREDENTIALS_JSON) {
    Write-Host "  Email: Configured" -ForegroundColor Green
} else {
    Write-Host "  Email: Not configured" -ForegroundColor Yellow
}
Write-Host ""
if ($env:STRIPE_SECRET_KEY_PROD -and $env:STRIPE_SECRET_KEY_PROD -notmatch "YOUR_PROD_KEY") {
    Write-Host "⚠️  WARNING: Using LIVE Stripe key - real charges will occur!" -ForegroundColor Red
} else {
    Write-Host "ℹ️  Using test keys (safe for development)" -ForegroundColor Green
}
Write-Host ""
Write-Host "Test the final invoice endpoint:" -ForegroundColor Yellow
Write-Host "  curl -X POST http://localhost:$env:PORT/api/stripe/final-invoice \""
Write-Host "    -H `"X-Api-Key: $env:SERVICE_API_KEY`" \""
Write-Host "    -H `"Content-Type: application/json`" \""
Write-Host "    -d '{\"email\":\"your-email@example.com\",\"name\":\"Test Customer\",\"totalAmount\":1000.0,\"depositPaid\":400.0}'"
Write-Host ""
Write-Host "Press Ctrl+C to stop the server"
Write-Host ""

go run ./cmd/api
