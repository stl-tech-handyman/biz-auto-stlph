@echo off
REM Test Email Analysis Service for Windows
REM Usage: test-email-analysis.bat [api-key] [max-emails]

set API_KEY=%1
if "%API_KEY%"=="" set API_KEY=your-api-key-here

set MAX_EMAILS=%2
if "%MAX_EMAILS%"=="" set MAX_EMAILS=50

set BASE_URL=%3
if "%BASE_URL%"=="" set BASE_URL=http://localhost:8080

echo Testing Email Analysis Service
echo ==============================
echo API Key: %API_KEY:~0,10%...
echo Max Emails: %MAX_EMAILS%
echo Base URL: %BASE_URL%
echo.

echo 1. Starting email analysis...
curl -X POST "%BASE_URL%/api/email-analysis/analyze" ^
  -H "X-API-Key: %API_KEY%" ^
  -H "Content-Type: application/json" ^
  -d "{\"max_emails\": %MAX_EMAILS%, \"query\": \"from:zapier.com OR subject:\\\"New Lead\\\"\", \"resume\": false}"

echo.
echo.
echo Done!
echo.
echo To check status, use:
echo curl -X GET "%BASE_URL%/api/email-analysis/status?spreadsheet_id=YOUR_SPREADSHEET_ID" -H "X-API-Key: %API_KEY%"
