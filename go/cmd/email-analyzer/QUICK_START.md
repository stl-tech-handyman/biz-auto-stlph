# 🚀 Email Analyzer - Quick Start (Local in Cursor)

## ✅ Ready to Run!

This local analyzer runs in Cursor with AI-assisted reasoning and handles:
- ✅ Email migration: `stlpartyhelpers@gmail.com` → `team@stlpartyhelpers.com`
- ✅ Forwarding detection: Identifies when conversations were forwarded
- ✅ Conversation grouping: Groups by thread + normalized client email

## Quick Start

### 1. Set Environment Variable

```bash
# Windows
set GMAIL_CREDENTIALS_JSON=C:\path\to\credentials.json

# Linux/Mac  
export GMAIL_CREDENTIALS_JSON=/path/to/credentials.json
```

### 2. Run It!

**Windows:**
```cmd
cd go\scripts
run-email-analyzer.bat
```

**Or directly:**
```cmd
cd go\cmd\email-analyzer
email-analyzer.exe -max 50 -workers 3 -v
```

**Linux/Mac:**
```bash
cd go/cmd/email-analyzer
go run main.go -max 50 -workers 3 -v
```

## What It Does

1. **Reads emails** from Gmail using your query
2. **Normalizes emails** - converts old email to new email
3. **Detects forwarding** - flags migrated conversations
4. **Extracts data** - client emails, pricing, events, etc.
5. **Writes to Google Sheets** - creates spreadsheet with results

## Command Options

```
-max int          Maximum emails to process (0 = process all available)
-query string     Gmail search query
-spreadsheet string  Existing spreadsheet ID
-resume           Resume from last position
-workers int      Number of concurrent workers (recommended: 3-5)
-v                Verbose logging (see what's happening)
```

## Examples

### Test Run (50 emails)
```bash
go run main.go -max 50 -workers 3 -v
```

### Process 500 emails
```bash
go run main.go -max 500 -workers 5 -v
```

### Custom query
```bash
go run main.go -max 100 -workers 5 -query "from:zapier.com" -v
```

### Resume processing
```bash
go run main.go -max 1000 -workers 5 -resume -spreadsheet YOUR_ID -v
```

## Recommended Worker Counts

- **3 workers**: safe default for long runs
- **5 workers**: faster, still usually safe
- **> 5**: may hit Gmail API rate limits depending on your quota / traffic

## Output

Creates a Google Sheet with:

- **Raw Data** - All emails with:
  - Original and normalized emails
  - Migration detection
  - Forwarded from tracking
  - All extracted fields

- **Email Mapping** - Shows email normalization

- **State** - Processing state for resume

## Email Migration Features

✅ **Automatic Normalization**
- `stlpartyhelpers@gmail.com` → `team@stlpartyhelpers.com`
- Applied to From, To, and Client Email fields

✅ **Forwarding Detection**
- Flags when emails were forwarded from old to new
- Tracks "Forwarded From" field

✅ **Conversation Continuity**
- Groups conversations even when email changed
- Uses thread ID + normalized client email

## Verbose Mode

Use `-v` to see:
- Each batch being processed
- Migration detections: `🔄 Migration detected: old -> new`
- Email normalization
- Processing progress

## Next Steps

1. **Run a test** (50 emails) to see it work
2. **Check the spreadsheet** - review extracted data
3. **Process more** - scale up gradually
4. **Review results** - we can improve extraction logic based on findings

## Troubleshooting

**"GMAIL_CREDENTIALS_JSON not set"**
- Set the environment variable before running

**"No emails found"**
- Check your query
- Try: `-query "from:zapier.com"`

**Build errors**
- Run: `go mod tidy` in the email-analyzer directory

## Ready!

Just run it and watch it process your emails! 🎉
