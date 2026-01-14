# Enable Google Sheets API

## Issue
The analysis failed because Google Sheets API is not enabled in your project.

## Solution

1. **Open Google Cloud Console:**
   https://console.developers.google.com/apis/api/sheets.googleapis.com/overview?project=254762852232

2. **Click "Enable"** to enable the Google Sheets API

3. **Wait 1-2 minutes** for the API to propagate

4. **Run the analysis again:**
   ```bash
   cd go/cmd/email-analyzer
   export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
   go run main.go -max 5000 -workers 5 -job "JOB-1-INITIAL-INDEXING" -job-name "Initial Indexing" -v
   ```

## Required APIs

Make sure these APIs are enabled:
- ✅ Gmail API (likely already enabled)
- ❌ Google Sheets API (needs to be enabled)
- ✅ Google Drive API (for folder access, if needed)

## Quick Enable Link

Click here to enable: https://console.developers.google.com/apis/api/sheets.googleapis.com/overview?project=254762852232
