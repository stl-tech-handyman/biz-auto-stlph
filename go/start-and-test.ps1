# Complete startup and test script
# Kills processes on ports 8080/8081, starts server, and tests all endpoints

$PORT = "8080"
$ALT_PORT = "8081"

# Set environment variables
$env:STRIPE_SECRET_KEY_PROD="sk_live_YOUR_PROD_KEY_HERE"
$env:SERVICE_API_KEY="test-api-key-12345"
$env:ENV="dev"
$env:LOG_LEVEL="debug"
$env:CONFIG_DIR="./config"
$env:TEMPLATES_DIR="./templates"

Write-Host "🔄 Restarting and Testing Go API Server..." -ForegroundColor Yellow
Write-Host ""

# Kill processes on both ports
Write-Host "Killing processes on ports $PORT and $ALT_PORT..." -ForegroundColor Cyan
$processes8080 = Get-NetTCPConnection -LocalPort $PORT -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
$processes8081 = Get-NetTCPConnection -LocalPort $ALT_PORT -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique

$allProcesses = @()
if ($processes8080) { $allProcesses += $processes8080 }
if ($processes8081) { $allProcesses += $processes8081 }

if ($allProcesses.Count -gt 0) {
    foreach ($pid in $allProcesses) {
        $procInfo = Get-Process -Id $pid -ErrorAction SilentlyContinue
        if ($procInfo) {
            Write-Host "   Killing process: $pid ($($procInfo.ProcessName))" -ForegroundColor Red
            Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
        }
    }
    Write-Host "   Waiting 3 seconds for ports to be released..." -ForegroundColor Cyan
    Start-Sleep -Seconds 3
}

# Check if port 8080 is free
$processAfter = Get-NetTCPConnection -LocalPort $PORT -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
if ($processAfter) {
    Write-Host "⚠️  Port $PORT is still in use, switching to port $ALT_PORT" -ForegroundColor Yellow
    $PORT = $ALT_PORT
} else {
    Write-Host "✅ Port $PORT is free" -ForegroundColor Green
}

$env:PORT=$PORT

Write-Host ""
Write-Host "🚀 Starting Go API Server on port $PORT..." -ForegroundColor Green
Write-Host ""
Write-Host "Configuration:" -ForegroundColor Yellow
Write-Host "  Stripe: LIVE (Production) Key" -ForegroundColor Red
Write-Host "  API Key: $env:SERVICE_API_KEY" -ForegroundColor Cyan
Write-Host "  Port: $env:PORT" -ForegroundColor Cyan
Write-Host "  Environment: $env:ENV" -ForegroundColor Cyan
Write-Host ""
Write-Host "Server URL: http://localhost:$env:PORT" -ForegroundColor Green
Write-Host ""
Write-Host "⚠️  WARNING: Using LIVE Stripe key - real charges will occur!" -ForegroundColor Red
Write-Host ""
Write-Host "Starting server in background..." -ForegroundColor Cyan
Write-Host ""

# Start server in background
$job = Start-Job -ScriptBlock {
    param($port, $stripeKey, $apiKey, $env, $logLevel, $configDir, $templatesDir)
    $env:STRIPE_SECRET_KEY_PROD=$stripeKey
    $env:SERVICE_API_KEY=$apiKey
    $env:ENV=$env
    $env:PORT=$port
    $env:LOG_LEVEL=$logLevel
    $env:CONFIG_DIR=$configDir
    $env:TEMPLATES_DIR=$templatesDir
    Set-Location $using:PWD
    go run ./cmd/api
} -ArgumentList $PORT, $env:STRIPE_SECRET_KEY_PROD, $env:SERVICE_API_KEY, $env:ENV, $env:LOG_LEVEL, $env:CONFIG_DIR, $env:TEMPLATES_DIR

Write-Host "Waiting for server to start..." -ForegroundColor Cyan
Start-Sleep -Seconds 8

# Test server
$BASE_URL = "http://localhost:$PORT"
Write-Host ""
Write-Host "🧪 Testing server..." -ForegroundColor Cyan

# Test health
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/api/health" -Method GET -TimeoutSec 5
    Write-Host "✅ Health check: OK" -ForegroundColor Green
    Write-Host "   Service: $($health.service)" -ForegroundColor Gray
    Write-Host "   Environment: $($health.environment)" -ForegroundColor Gray
    if ($health.debug) {
        Write-Host "   SERVICE_API_KEY set: $($health.debug.serviceApiKeySet)" -ForegroundColor $(if ($health.debug.serviceApiKeySet) { "Green" } else { "Red" })
    } else {
        Write-Host "   ⚠️  Debug field not found - server may not have restarted with new code" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ Health check failed: $_" -ForegroundColor Red
}

# Test protected endpoint
Write-Host ""
Write-Host "Testing protected endpoint..." -ForegroundColor Cyan
try {
    $stripeTest = Invoke-RestMethod -Uri "$BASE_URL/api/stripe/test" -Method POST -Headers @{"X-Api-Key"=$env:SERVICE_API_KEY; "Content-Type"="application/json"} -Body '{}' -TimeoutSec 5
    Write-Host "✅ Stripe test: OK" -ForegroundColor Green
} catch {
    $errorMsg = $_.ErrorDetails.Message
    if ($errorMsg -like "*Service Configuration Error*") {
        Write-Host "❌ Stripe test: SERVICE_API_KEY not configured" -ForegroundColor Red
        Write-Host "   Server is running but cannot see SERVICE_API_KEY environment variable" -ForegroundColor Yellow
    } else {
        Write-Host "⚠️  Stripe test: $errorMsg" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "=== Server Status ===" -ForegroundColor Cyan
Write-Host "URL: $BASE_URL" -ForegroundColor Green
Write-Host "Job ID: $($job.Id)" -ForegroundColor Gray
Write-Host ""
Write-Host "To stop server: Stop-Job -Id $($job.Id); Remove-Job -Id $($job.Id)" -ForegroundColor Yellow
Write-Host ""
Write-Host "Press any key to continue (server will keep running)..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

