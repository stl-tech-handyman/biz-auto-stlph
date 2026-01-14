# Email Analyzer - Local Run in Cursor

This is a **local, AI-assisted email analyzer** that runs directly in Cursor. Unlike the Cloud Run version, this can adapt and reason as it processes emails.

## Features

✅ **Email Migration Handling**
- Automatically treats `stlpartyhelpers@gmail.com` as `team@stlpartyhelpers.com`
- Detects forwarded conversations
- Tracks email mapping

✅ **Conversation Detection**
- Identifies when client wrote to old email but conversation continues from new email
- Groups related emails by thread + normalized client email

✅ **AI-Assisted Processing**
- Can be enhanced with reasoning logic
- Verbose logging for debugging
- Adapts patterns as it learns

## Quick Start

### 1. Set Environment Variable

```bash
# Windows
set GMAIL_CREDENTIALS_JSON=C:\path\to\credentials.json

# Linux/Mac
export GMAIL_CREDENTIALS_JSON=/path/to/credentials.json
```

### 2. Run

**Windows:**
```cmd
cd go\scripts
run-email-analyzer.bat
```

**Linux/Mac:**
```bash
cd go/scripts
chmod +x run-email-analyzer.sh
./run-email-analyzer.sh
```

**Or directly:**
```bash
cd go/cmd/email-analyzer
go run main.go -max 50 -v
```

## Command Line Options

```
-max int
    Maximum emails to process (default: 100)

-query string
    Gmail search query (default: searches for form submissions)

-spreadsheet string
    Existing spreadsheet ID (creates new if empty)

-resume
    Resume from last position

-v
    Verbose logging
```

## Examples

### Process 50 emails (test)
```bash
go run main.go -max 50 -v
```

### Process 500 emails with custom query
```bash
go run main.go -max 500 -query "from:zapier.com" -v
```

### Resume processing
```bash
go run main.go -max 1000 -resume -spreadsheet YOUR_SPREADSHEET_ID -v
```

## Output

Creates a Google Sheet with:

1. **Raw Data** - All processed emails with:
   - Original and normalized emails
   - Migration detection flags
   - Forwarded from tracking
   - All extracted data

2. **Email Mapping** - Shows how emails were normalized:
   - `stlpartyhelpers@gmail.com` → `team@stlpartyhelpers.com`

3. **State** - Processing state for resume

4. **Processing Log** - Processing history

## Email Migration Logic

The analyzer automatically:

1. **Normalizes emails**: `stlpartyhelpers@gmail.com` → `team@stlpartyhelpers.com`
2. **Detects forwarding**: Flags emails that were forwarded
3. **Tracks conversations**: Groups by thread + normalized client email
4. **Maps relationships**: Shows original → normalized email mapping

## Verbose Mode

Use `-v` flag to see:
- Each batch being processed
- Migration detections
- Email normalization
- Processing progress

## Next Steps

This local version can be enhanced with:
- Pattern learning from samples
- Adaptive extraction rules
- Better conversation grouping
- AI reasoning for edge cases

Run it, review the results, and we can improve the logic based on what we find!
