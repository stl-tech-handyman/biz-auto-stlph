# Test all API endpoints
# This script tests all available endpoints to ensure server is working correctly

$API_KEY = "test-api-key-12345"
$BASE_URL = "http://localhost:8080"

# Try 8081 if 8080 doesn't respond
try {
    $response = Invoke-WebRequest -Uri "$BASE_URL/api/health" -Method GET -TimeoutSec 2 -ErrorAction Stop
    if ($response.StatusCode -ne 200) {
        $BASE_URL = "http://localhost:8081"
        Write-Host "Port 8080 returned non-200, trying 8081..." -ForegroundColor Yellow
    }
} catch {
    $BASE_URL = "http://localhost:8081"
    Write-Host "Port 8080 not responding, trying 8081..." -ForegroundColor Yellow
}

Write-Host "🧪 Testing all API endpoints..." -ForegroundColor Cyan
Write-Host "Base URL: $BASE_URL" -ForegroundColor Cyan
Write-Host "API Key: $API_KEY" -ForegroundColor Cyan
Write-Host ""

$passed = 0
$failed = 0

# Test function
function Test-Endpoint {
    param(
        [string]$Method,
        [string]$Path,
        [string]$Body = "",
        [bool]$RequiresAuth = $false,
        [string]$Description = ""
    )
    
    $url = "$BASE_URL$Path"
    $headers = @{}
    
    if ($RequiresAuth) {
        $headers["X-Api-Key"] = $API_KEY
    }
    
    $headers["Content-Type"] = "application/json"
    
    Write-Host "Testing: $Method $Path" -ForegroundColor Yellow -NoNewline
    
    try {
        if ($Method -eq "GET") {
            $response = Invoke-WebRequest -Uri $url -Method GET -Headers $headers -ErrorAction Stop
        } else {
            $response = Invoke-WebRequest -Uri $url -Method $Method -Headers $headers -Body $Body -ErrorAction Stop
        }
        
        $status = $response.StatusCode
        if ($status -ge 200 -and $status -lt 300) {
            Write-Host " ✅ OK ($status)" -ForegroundColor Green
            $script:passed++
            return $true
        } elseif ($status -eq 401 -or $status -eq 403) {
            Write-Host " ⚠️  Auth required ($status)" -ForegroundColor Yellow
            $script:passed++
            return $true
        } elseif ($status -eq 500) {
            $content = $response.Content | ConvertFrom-Json
            if ($content.error -eq "Service Configuration Error") {
                Write-Host " ❌ SERVICE_API_KEY not configured ($status)" -ForegroundColor Red
            } else {
                Write-Host " ❌ Server error ($status): $($content.message)" -ForegroundColor Red
            }
            $script:failed++
            return $false
        } else {
            Write-Host " ⚠️  Status: $status" -ForegroundColor Yellow
            $script:passed++
            return $true
        }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        $errorMsg = $_.Exception.Message
        Write-Host " ❌ Failed: $errorMsg ($statusCode)" -ForegroundColor Red
        $script:failed++
        return $false
    }
}

# Health endpoints (no auth)
Write-Host "=== Health Endpoints ===" -ForegroundColor Cyan
Test-Endpoint -Method "GET" -Path "/" -RequiresAuth $false -Description "Root endpoint"
Test-Endpoint -Method "GET" -Path "/api/health" -RequiresAuth $false -Description "Health check"
Test-Endpoint -Method "GET" -Path "/api/health/ready" -RequiresAuth $false -Description "Readiness check"
Test-Endpoint -Method "GET" -Path "/api/health/live" -RequiresAuth $false -Description "Liveness check"
Write-Host ""

# Stripe endpoints (require auth)
Write-Host "=== Stripe Endpoints ===" -ForegroundColor Cyan
Test-Endpoint -Method "POST" -Path "/api/stripe/test" -Body '{}' -RequiresAuth $true -Description "Stripe test"
Test-Endpoint -Method "GET" -Path "/api/stripe/deposit/calculate?estimate=1000" -RequiresAuth $true -Description "Calculate deposit"
Test-Endpoint -Method "POST" -Path "/api/stripe/deposit/amount" -Body '{"estimatedTotal":1000.0}' -RequiresAuth $true -Description "Get deposit amount"
Write-Host ""

# Estimate endpoints (require auth)
Write-Host "=== Estimate Endpoints ===" -ForegroundColor Cyan
Test-Endpoint -Method "POST" -Path "/api/estimate" -Body '{"eventDate":"2025-12-25","durationHours":4,"numHelpers":2}' -RequiresAuth $true -Description "Calculate estimate"
Test-Endpoint -Method "GET" -Path "/api/estimate/special-dates?years=5" -RequiresAuth $true -Description "Get special dates"
Write-Host ""

# Email endpoints (require auth)
Write-Host "=== Email Endpoints ===" -ForegroundColor Cyan
Test-Endpoint -Method "POST" -Path "/api/email/test" -Body '{"to":"test@example.com","subject":"Test","html":"<p>Test</p>"}' -RequiresAuth $true -Description "Test email"
Write-Host ""

# Public endpoints (no auth)
Write-Host "=== Public Endpoints ===" -ForegroundColor Cyan
Test-Endpoint -Method "GET" -Path "/swagger" -RequiresAuth $false -Description "Swagger UI"
Test-Endpoint -Method "GET" -Path "/api/openapi.json" -RequiresAuth $false -Description "OpenAPI spec"
Write-Host ""

# Summary
Write-Host "=== Test Summary ===" -ForegroundColor Cyan
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor $(if ($failed -eq 0) { "Green" } else { "Red" })
Write-Host ""

if ($failed -eq 0) {
    Write-Host "✅ All tests passed!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "❌ Some tests failed. Check the output above." -ForegroundColor Red
    exit 1
}

