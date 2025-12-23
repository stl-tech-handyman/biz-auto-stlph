@echo off
echo ========================================
echo  Final Invoice Test - Quick Test
echo ========================================
echo.

set API_KEY=test-api-key-12345
set BASE_URL=http://localhost:8080

echo Enter your email address:
set /p TEST_EMAIL=

if "%TEST_EMAIL%"=="" (
    echo Error: Email is required
    pause
    exit /b 1
)

echo.
echo [1/2] Creating final invoice...
curl -s -X POST "%BASE_URL%/api/stripe/final-invoice" ^
  -H "X-Api-Key: %API_KEY%" ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"%TEST_EMAIL%\",\"name\":\"Test Customer\",\"totalAmount\":1000.0,\"depositPaid\":400.0}" > temp_invoice.json

type temp_invoice.json
echo.

echo [2/2] Extracting invoice data and sending email...
powershell -Command "$json = Get-Content temp_invoice.json | ConvertFrom-Json; $url = $json.invoice.url; $balance = $json.details.remainingBalance; curl -s -X POST 'http://localhost:8080/api/email/final-invoice' -H 'X-Api-Key: test-api-key-12345' -H 'Content-Type: application/json' -d (\"{\\\"name\\\":\\\"Test Customer\\\",\\\"email\\\":\\\"%TEST_EMAIL%\\\",\\\"totalAmount\\\":1000.0,\\\"depositPaid\\\":400.0,\\\"remainingBalance\\\":\" + $balance + \",\\\"invoiceUrl\\\":\\\"\" + $url + \"\\\"}\")"

echo.
echo Done! Check your email: %TEST_EMAIL%
del temp_invoice.json 2>nul
pause

