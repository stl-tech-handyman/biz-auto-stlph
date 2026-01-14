# ✅ Email Analysis Service - Ready to Test!

## 🎉 Setup Complete

The email analysis service is now integrated into your Cloud Run API and ready to test!

## 📋 What Was Created

1. **Google Sheets Client** (`go/internal/infra/sheets/sheets.go`)
   - Handles all Google Sheets API operations
   - Creates spreadsheets, writes data, manages state

2. **Email Analysis Service** (`go/internal/services/email_analysis/service.go`)
   - Processes Gmail emails
   - Extracts business data (client emails, pricing, events, etc.)
   - Writes to Google Sheets
   - Manages processing state

3. **HTTP Handler** (`go/internal/http/handlers/email_analysis_handler.go`)
   - REST API endpoints for email analysis
   - Status checking

4. **Router Integration** (`go/internal/http/router.go`)
   - Added routes: `/api/email-analysis/analyze` and `/api/email-analysis/status`
   - Protected with API key authentication

5. **Gmail Client Extensions** (`go/internal/infra/email/gmail.go`)
   - Added `GetService()` and `GetFromEmail()` methods

## 🚀 How to Run & Test

### Step 1: Start the Server

```bash
cd go
go run cmd/api/main.go
```

Or if already running, restart it.

### Step 2: Test with Small Batch (Recommended First)

**Windows:**
```cmd
cd go\scripts
test-email-analysis.bat your-api-key 50
```

**Linux/Mac:**
```bash
cd go/scripts
chmod +x test-email-analysis.sh
./test-email-analysis.sh your-api-key 50
```

**Or manually:**
```bash
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "max_emails": 50,
    "query": "from:zapier.com",
    "resume": false
  }'
```

### Step 3: Check the Response

You'll get a response like:
```json
{
  "processed": 45,
  "skipped": 5,
  "total": 45,
  "spreadsheet_id": "abc123...",
  "spreadsheet_url": "https://docs.google.com/spreadsheets/d/abc123...",
  "next_index": 50,
  "has_more": true
}
```

**Open the `spreadsheet_url` in your browser** to see the results!

### Step 4: Check Status

```bash
curl -X GET "http://localhost:8080/api/email-analysis/status?spreadsheet_id=YOUR_SPREADSHEET_ID" \
  -H "X-API-Key: your-api-key"
```

### Step 5: Process More Emails

```bash
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "max_emails": 500,
    "resume": true,
    "spreadsheet_id": "YOUR_SPREADSHEET_ID"
  }'
```

## 📊 What Gets Created

### Google Sheets Structure

1. **Raw Data Sheet**
   - All processed emails with extracted fields
   - Columns: Email ID, Date, From, Subject, Client Email, Total Cost, Rate, etc.

2. **State Sheet**
   - Processing state (for resume functionality)
   - Tracks: LastProcessedIndex, TotalProcessed, LastRun

3. **Processing Log Sheet**
   - Processing history and errors

## 📝 Logging

### Debug Logs

All operations are logged to `.cursor/debug.log`:

```bash
# View recent logs
tail -f .cursor/debug.log

# Filter for email analysis
grep "email-analysis" .cursor/debug.log

# View processing batches
grep "Processing batch" .cursor/debug.log
```

### What's Logged

- Batch processing progress
- Email extraction results
- Errors and warnings
- State updates
- Spreadsheet operations

## 🔍 Monitoring

### Check Processing Progress

1. **Via API:**
   ```bash
   curl -X GET "http://localhost:8080/api/email-analysis/status?spreadsheet_id=ID" \
     -H "X-API-Key: your-api-key"
   ```

2. **Via Google Sheets:**
   - Open the spreadsheet
   - Check "State" sheet for progress
   - Check "Raw Data" sheet for processed emails

3. **Via Logs:**
   ```bash
   tail -f .cursor/debug.log | grep "email-analysis"
   ```

## ⚙️ Configuration

### Environment Variables

**Required:**
```bash
GMAIL_CREDENTIALS_JSON=/path/to/credentials.json
# OR
GMAIL_CREDENTIALS_JSON='{"type":"service_account",...}'
```

**For Cloud Run:**
```bash
gcloud run services update <service-name> \
  --set-env-vars GMAIL_CREDENTIALS_JSON=<json-string>
```

### API Key

Use your existing API key (same one used for other endpoints).

## 🎯 Processing Strategy

### For Testing (Start Here)
```bash
# Process 50 emails
curl -X POST ... -d '{"max_emails": 50}'
```

### For Production
```bash
# Process 1000 emails per run
curl -X POST ... -d '{"max_emails": 1000, "resume": true, "spreadsheet_id": "ID"}'

# Repeat until all processed
```

### For 86,000 Emails

Process in batches of 1000, using `resume: true`:

```bash
# Run multiple times, each time it resumes from last position
for i in {1..86}; do
  curl -X POST ... -d '{"max_emails": 1000, "resume": true, "spreadsheet_id": "ID"}'
  sleep 60
done
```

## 🐛 Troubleshooting

### "GMAIL_CREDENTIALS_JSON not set"
- Set the environment variable before starting server
- Check it's exported: `echo $GMAIL_CREDENTIALS_JSON`

### "failed to create sheets service"
- Ensure Google Sheets API is enabled in GCP Console
- Check credentials have `spreadsheets` scope

### "No emails found"
- Check your Gmail query
- Try simpler query: `from:zapier.com`
- Verify you have matching emails

### "quota exceeded"
- Wait for daily quota reset
- Reduce batch size
- Add delays between requests

## ✅ Next Steps

1. **Test with 50 emails** - Verify it works
2. **Check the spreadsheet** - Review extracted data
3. **Process more** - Scale up to 500-1000 emails
4. **Monitor logs** - Watch for errors
5. **Resume processing** - Use `resume: true` to continue

## 📚 Documentation

Full documentation: `go/docs/EMAIL_ANALYSIS_SETUP.md`

## 🎉 Ready to Go!

Everything is set up and ready. Start with a small test (50 emails) and scale up from there!

**First command to run:**
```bash
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"max_emails": 50}'
```

Then open the spreadsheet URL from the response!

