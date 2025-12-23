# Test script for final invoice functionality
# Make sure the server is running first!

$API_KEY = $env:SERVICE_API_KEY
if (-not $API_KEY) {
    $API_KEY = "test-api-key-12345"
}

$BASE_URL = "http://localhost:8080"
$TEST_EMAIL = Read-Host "Enter your email address for testing"

Write-Host "🧪 Testing Final Invoice Functionality" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""

# Step 1: Create final invoice
Write-Host "1️⃣ Creating final invoice..." -ForegroundColor Yellow
$invoiceBody = @{
    email = $TEST_EMAIL
    name = "Test Customer"
    totalAmount = 1000.0
    depositPaid = 400.0
    currency = "usd"
    description = "Final payment for test event"
} | ConvertTo-Json

try {
    $invoiceResponse = Invoke-RestMethod -Uri "$BASE_URL/api/stripe/final-invoice" `
        -Method POST `
        -Headers @{
            "X-Api-Key" = $API_KEY
            "Content-Type" = "application/json"
        } `
        -Body $invoiceBody

    Write-Host "✅ Invoice created successfully!" -ForegroundColor Green
    Write-Host ""
    $invoiceResponse | ConvertTo-Json -Depth 10
    
    $invoiceURL = $invoiceResponse.invoice.url
    $remainingBalance = $invoiceResponse.details.remainingBalance
    
    if (-not $invoiceURL) {
        Write-Host "❌ Failed to get invoice URL from response" -ForegroundColor Red
        exit 1
    }
    
    Write-Host ""
    Write-Host "Invoice URL: $invoiceURL" -ForegroundColor Cyan
    Write-Host ""
    
    # Step 2: Send email
    Write-Host "2️⃣ Sending final invoice email..." -ForegroundColor Yellow
    $emailBody = @{
        name = "Test Customer"
        email = $TEST_EMAIL
        totalAmount = 1000.0
        depositPaid = 400.0
        remainingBalance = $remainingBalance
        invoiceUrl = $invoiceURL
    } | ConvertTo-Json
    
    try {
        $emailResponse = Invoke-RestMethod -Uri "$BASE_URL/api/email/final-invoice" `
            -Method POST `
            -Headers @{
                "X-Api-Key" = $API_KEY
                "Content-Type" = "application/json"
            } `
            -Body $emailBody
        
        Write-Host "✅ Email sent successfully!" -ForegroundColor Green
        Write-Host ""
        $emailResponse | ConvertTo-Json -Depth 10
        
        Write-Host ""
        Write-Host "✅ Test complete! Check your email: $TEST_EMAIL" -ForegroundColor Green
        Write-Host "   Invoice URL: $invoiceURL" -ForegroundColor Cyan
        
    } catch {
        Write-Host "❌ Failed to send email: $_" -ForegroundColor Red
        Write-Host "   Make sure GMAIL_CREDENTIALS_JSON or EMAIL_SERVICE_URL is configured" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host "❌ Failed to create invoice: $_" -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit 1
}

