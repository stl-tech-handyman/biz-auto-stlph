# 🚀 Start Email Analysis Now

## Quick Start

The dashboard is now set to auto-detect the spreadsheet. Just start the analysis:

```bash
cd go/cmd/email-analyzer

# Set credentials (use your actual path or JSON string)
export GMAIL_CREDENTIALS_JSON="/path/to/credentials.json"
# OR
export GMAIL_CREDENTIALS_JSON='{"type":"service_account",...}'

# Start analysis (will create spreadsheet automatically)
go run main.go -max 5000 -workers 5 -job "JOB-1-INITIAL-INDEXING" -job-name "Initial Indexing" -v
```

## What Happens

1. ✅ **Creates Spreadsheet** - Automatically creates Google Sheet
2. ✅ **Prints Dashboard URL** - Copy the URL shown in console
3. ✅ **Opens Dashboard** - Paste URL in browser, stats auto-load
4. ✅ **No Input Needed** - Dashboard detects spreadsheet automatically

## Dashboard

Once analysis starts, you'll see:
```
✅ Created spreadsheet: https://docs.google.com/spreadsheets/d/...
📊 Spreadsheet ID: abc123...
🌐 Dashboard URL: http://localhost:8080/email-analysis-dashboard.html?spreadsheet_id=abc123...
```

Just open that dashboard URL - it will auto-detect and show stats!
