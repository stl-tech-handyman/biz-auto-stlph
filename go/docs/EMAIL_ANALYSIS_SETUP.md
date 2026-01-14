# Email Analysis Service - Setup & Usage Guide

## Overview

The Email Analysis Service processes Gmail emails and extracts business data into Google Sheets. It runs on Cloud Run and can process thousands of emails efficiently.

## Two Ways to Run Email Analysis

1. **Cloud Run API (service endpoints)** — `POST /api/email-analysis/analyze`
2. **Local Analyzer CLI (recommended for big one-time backfills)** — `go/cmd/email-analyzer`

The local analyzer includes:
- fast processing using **goroutine workers** (`-workers N`)
- robust **resume** using Gmail Message IDs
- spreadsheet **lock** to avoid multi-process conflicts
- a live dashboard at `/email-analysis-dashboard.html`

## Prerequisites

1. **Gmail API Credentials** - Same credentials used for sending emails
2. **Google Sheets API Access** - Enabled in your GCP project
3. **Environment Variables** - `GMAIL_CREDENTIALS_JSON` must be set

## Environment Setup

### Required Environment Variable

```bash
GMAIL_CREDENTIALS_JSON=<path to credentials.json or JSON string>
```

The credentials must have these scopes:
- `https://www.googleapis.com/auth/gmail.readonly` (for reading emails)
- `https://www.googleapis.com/auth/spreadsheets` (for writing to Sheets)

## API Endpoints

### 1. Analyze Emails

**POST** `/api/email-analysis/analyze`

**Headers:**
```
X-API-Key: <your-api-key>
Content-Type: application/json
```

**Request Body:**
```json
{
  "max_emails": 100,
  "query": "from:zapier.com OR subject:\"New Lead\"",
  "resume": false,
  "spreadsheet_id": ""
}
```

**Parameters:**
- `max_emails` (optional): Maximum emails to process (default: 100, max: 1000 per run)
- `query` (optional): Gmail search query (default: searches for form submissions)
- `resume` (optional): Resume from last checkpoint (default: false)
- `spreadsheet_id` (optional): Use existing spreadsheet (default: creates new one)

**Response:**
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

### 2. Get Status

**GET** `/api/email-analysis/status?spreadsheet_id=<id>`

**Headers:**
```
X-API-Key: <your-api-key>
```

**Response:**
```json
{
  "last_processed_index": 50,
  "total_processed": 45,
  "last_run": "2025-01-15T10:30:00Z",
  "spreadsheet_id": "abc123..."
}
```

## How to Run

### 1. Local Testing

```bash
# Set environment variable
export GMAIL_CREDENTIALS_JSON="/path/to/credentials.json"

# Start server
cd go
go run cmd/api/main.go

# In another terminal, test the endpoint
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "max_emails": 50,
    "query": "from:zapier.com",
    "resume": false
  }'
```

### 2. Cloud Run Deployment

The service is automatically included when you deploy your API. Make sure:

1. **Set environment variable in Cloud Run:**
   ```bash
   gcloud run services update <service-name> \
     --set-env-vars GMAIL_CREDENTIALS_JSON=<json-string>
   ```

2. **Or use Secret Manager:**
   ```bash
   # Store credentials in Secret Manager
   echo -n '{"type":"service_account",...}' | \
     gcloud secrets create gmail-credentials --data-file=-
   
   # Grant access
   gcloud secrets add-iam-policy-binding gmail-credentials \
     --member=serviceAccount:<service-account>@<project>.iam.gserviceaccount.com \
     --role=roles/secretmanager.secretAccessor
   ```

## Processing Strategy

### First Run (Small Test)
```bash
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"max_emails": 50}'
```

### Process More Emails
```bash
# Process 500 emails
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"max_emails": 500, "resume": true}'
```

### Resume Processing
```bash
# Get current status first
curl -X GET "http://localhost:8080/api/email-analysis/status?spreadsheet_id=<id>" \
  -H "X-API-Key: your-api-key"

# Resume from where you left off
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "max_emails": 1000,
    "resume": true,
    "spreadsheet_id": "<your-spreadsheet-id>"
  }'
```

## Logging

### Debug Logs

All operations are logged to `.cursor/debug.log` with detailed information:
- Batch processing progress
- Email extraction results
- Errors and warnings
- State updates

### View Logs

```bash
# Tail the log file
tail -f .cursor/debug.log

# Search for email analysis logs
grep "email-analysis" .cursor/debug.log

# View recent processing
grep "Processing batch" .cursor/debug.log | tail -20
```

### Cloud Run Logs

```bash
# View logs in Cloud Run
gcloud run services logs read <service-name> --limit 50

# Follow logs
gcloud run services logs tail <service-name>
```

## Google Sheets Output

The service creates a spreadsheet with these sheets:

### 1. Raw Data
All processed emails with extracted fields:
- Email ID, Thread ID, Date, From, Subject
- Is Test, Is Confirmation, Client Email
- Event Date, Total Cost, Rate, Hours, Helpers
- Occasion, Status, Guests, Deposit
- Email Type, Conversation ID, Message Number

### 2. State
Processing state (for resume functionality):
- LastProcessedIndex
- TotalProcessed
- LastRun
- SpreadsheetID

### 3. Processing Log
Processing history and errors

## Rate Limits & Best Practices

### Gmail API Limits
- **Quota**: 1,000,000 quota units per day
- **Read operation**: ~5 quota units per message
- **Safe batch size**: 50-100 emails per request

### Recommendations

1. **Start Small**: Test with 50-100 emails first
2. **Batch Processing**: Process in chunks of 500-1000 emails
3. **Resume Feature**: Use `resume: true` to continue from last position
4. **Monitor Quota**: Check GCP Console for quota usage

### Processing 86,000 Emails

```bash
# Strategy: Process in batches of 1000
for i in {1..86}; do
  curl -X POST http://localhost:8080/api/email-analysis/analyze \
    -H "X-API-Key: your-api-key" \
    -H "Content-Type: application/json" \
    -d "{\"max_emails\": 1000, \"resume\": true, \"spreadsheet_id\": \"<id>\"}"
  
  sleep 60  # Wait 1 minute between batches
done
```

## Troubleshooting

### Error: "GMAIL_CREDENTIALS_JSON not set"
- Set the environment variable before starting the server
- For Cloud Run, add it via `gcloud run services update`

### Error: "failed to create sheets service"
- Ensure Google Sheets API is enabled in GCP Console
- Check that credentials have `spreadsheets` scope

### Error: "quota exceeded"
- Wait for quota reset (daily)

## Local Analyzer (CLI) — Quick Start

See full docs in `go/cmd/email-analyzer/README.md`.

```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="/path/to/credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"

# Full run, 5 concurrent goroutine workers
go run main.go -all -workers 5 -job "JOB-1-INITIAL-INDEXING" -job-name "Initial Indexing" -v
```

Notes:
- Prefer **3–5** workers for long runs.
- Use `-resume -spreadsheet <id>` to restart without losing progress.
- The analyzer prints a dashboard URL: `/email-analysis-dashboard.html?spreadsheet_id=...`
- Reduce batch size
- Add delays between requests

### No emails found
- Check your Gmail query
- Verify you have emails matching the query
- Try a simpler query like `from:zapier.com`

## Next Steps

After processing emails:

1. **Review Raw Data** sheet in Google Sheets
2. **Identify Organizations** - Group by domain (future feature)
3. **Generate Reports** - Create summary sheets (future feature)
4. **Set up Automation** - Schedule regular processing (future feature)

## Example Workflow

```bash
# 1. Initial test (50 emails)
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"max_emails": 50}'

# 2. Check results in spreadsheet URL from response

# 3. Process more (500 emails)
curl -X POST http://localhost:8080/api/email-analysis/analyze \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "max_emails": 500,
    "resume": true,
    "spreadsheet_id": "<from-step-1>"
  }'

# 4. Continue until all processed
# Use resume: true to continue from last position
```

## Support

Check logs in `.cursor/debug.log` for detailed error messages and processing information.
