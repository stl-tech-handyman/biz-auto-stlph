# Starting Email Analysis

## Quick Start

Once you have `GMAIL_CREDENTIALS_JSON` set, run:

```bash
cd go/cmd/email-analyzer
go run main.go -all -workers 5 -job "JOB-1-INITIAL-INDEXING" -job-name "Initial Indexing and Classification Systematization" -v
```

## What Happens

1. **Creates Spreadsheet** - Automatically creates a new Google Sheet
2. **Acquires Lock** - Prevents concurrent processing conflicts
3. **Processes Emails** - Uses `-workers 5` goroutine workers inside a single run
4. **Tracks Progress** - Saves state every 25 emails and every 2 minutes
5. **Job Tracking** - All emails stamped with Job ID
6. **Resume Safe** - Uses Message IDs for perfect resume capability

## Output

The spreadsheet ID will be printed when created. Use it in the dashboard:
- Open: `/email-analysis-dashboard.html?spreadsheet_id=YOUR_ID`
- Or add it manually using "Add New" button

## Monitoring

- **Dashboard**: View live stats at `/email-analysis-dashboard.html`
- **Google Sheet**: Direct link shown in dashboard
- **Console**: Real-time progress with processing rates

## Resume

If stopped, resume with:
```bash
go run main.go -all -workers 5 -resume -spreadsheet YOUR_SPREADSHEET_ID -v
```
