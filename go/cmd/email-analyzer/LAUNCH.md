# 🚀 How to Launch Email Analyzer with Workers

## Quick Start

```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"
go run main.go -all -workers 5 -job "JOB-1-INITIAL-INDEXING" -job-name "Initial Indexing" -v
```

## Command Options

### Basic Usage
```bash
# Process all emails with 5 workers
go run main.go -all -workers 5 -v

# Process specific number of emails
go run main.go -max 1000 -workers 3 -v

# Resume from previous run
go run main.go -all -workers 5 -resume -spreadsheet SPREADSHEET_ID -v
```

### Flags

- `-workers N` - Number of concurrent workers (recommended: 3-5, max: 10)
- `-all` - Process all available emails
- `-max N` - Process maximum N emails
- `-resume` - Resume from last position (uses spreadsheet state)
- `-spreadsheet ID` - Use existing spreadsheet (creates new if empty)
- `-job "JOB-ID"` - Job identifier for tracking
- `-job-name "Name"` - Human-readable job name
- `-v` - Verbose logging
- `-query "..."` - Custom Gmail search query

## How Workers Work

1. **Producer Goroutine**: Fetches batches from Gmail API and sends message IDs to a channel
2. **Worker Goroutines** (N workers): Pull message IDs from channel and process them concurrently
3. **Coordination**: Mutexes protect shared state (counters, batch, processed IDs)
4. **Batch Writing**: Synchronized - only one worker writes batches at a time

## Performance

- **1 worker**: ~1-2 emails/second
- **3 workers**: ~3-5 emails/second  
- **5 workers**: ~5-8 emails/second
- **10 workers**: May hit Gmail API rate limits

## Example: Full Analysis

```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"

# Start full analysis with 5 workers
go run main.go \
  -all \
  -workers 5 \
  -job "JOB-1-INITIAL-INDEXING" \
  -job-name "Initial Indexing and Classification" \
  -v
```

The analyzer will:
- ✅ Create a new spreadsheet automatically
- ✅ Show the spreadsheet URL and dashboard link
- ✅ Process emails concurrently using 5 workers
- ✅ Auto-save state every 2 minutes
- ✅ Show real-time progress
- ✅ Resume safely if interrupted

## Monitoring

Open the dashboard URL shown in the output:
```
http://localhost:8080/email-analysis-dashboard.html?spreadsheet_id=SPREADSHEET_ID
```

The dashboard shows:
- Total processed/skipped
- Active workers
- Job history
- Direct link to Google Sheet
