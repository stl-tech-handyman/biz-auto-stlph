# 🚀 START HERE: Email Analyzer with Verification

## Current Status

✅ **Process is running** - Processing ALL emails with verification
✅ **14+ batches written** - Data is being written to spreadsheet
✅ **Verification enabled** - Each write is confirmed

**Spreadsheet:** https://docs.google.com/spreadsheets/d/1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU

## Quick Verification (Right Now)

### 1. Check if Process is Running
```bash
ps aux | grep "[g]o run main.go"
```

### 2. Check Latest Writes
```bash
tail -30 go/email-analyzer-verify.log | grep -E "(Writing|Successfully|Sheet now has)"
```

### 3. Check Google Sheet
- Open the spreadsheet link above
- Go to "Raw Data" tab
- Check row count (should be increasing)
- Scroll to bottom - should see new rows

### 4. Count Total Batches Written
```bash
grep -c "Successfully wrote batch" go/email-analyzer-verify.log
```

## If You Need to Restart

### Option 1: Test Mode First (Recommended)
```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"

go run main.go -test -idempotent \
  -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" \
  -job "JOB-TEST" -job-name "Test verification" \
  -query "" -batch 10 -workers 1
```

### Option 2: Full Run
```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"

nohup go run main.go -all -workers 1 -idempotent \
  -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" \
  -job "JOB-FULL-ALL" -job-name "Process ALL emails with verification" \
  -query "" -batch 25 -delay 100 -v \
  > ../../email-analyzer-full.log 2>&1 &

echo "Started (PID: $!)"
```

## What You'll See

### ✅ Good Signs (Everything Working):
```
💾 Writing batch of 25 emails to spreadsheet...
✅ Confirmed: 25 rows written to spreadsheet
📊 Sheet now has 350 total data rows (excluding header)
✅ Successfully wrote batch of 25 emails to spreadsheet
```

### ⚠️ Warnings (Usually OK):
```
⚠️ Warning: Written message ID X not found in sheet verification
```
- This is a timing issue, not a failure
- If you see "Successfully wrote" and row count increases, it's working

### ❌ Errors (Needs Attention):
```
❌ ERROR writing batch: [error details]
```
- Check error message
- Verify credentials are set
- Check sheet permissions

## Monitoring Commands

```bash
# Watch log in real-time
tail -f go/email-analyzer-verify.log

# Check progress every minute
watch -n 60 'tail -20 go/email-analyzer-verify.log | grep -E "(Progress|Wrote|Sheet now has)"'

# Count batches written
grep -c "Successfully wrote batch" go/email-analyzer-verify.log

# Check for errors
grep -i "error\|failed" go/email-analyzer-verify.log | tail -10
```

## Verification Checklist

Every 5-10 minutes, verify:

- [ ] Process is running (`ps aux | grep "[g]o run main.go"`)
- [ ] Log shows "Successfully wrote batch" messages
- [ ] Log shows "Sheet now has X total data rows" (increasing)
- [ ] Google Sheet "Raw Data" tab has new rows
- [ ] No ERROR messages in log
- [ ] Row count in sheet matches log count (within ~25 rows)

## If Something's Wrong

1. **Check process status** - Is it running?
2. **Check log for errors** - Any ERROR messages?
3. **Test write capability** - Run `-test` mode
4. **Check credentials** - Are env vars set?
5. **Check sheet directly** - Open Google Sheet and verify rows

## Expected Progress

- **First batch**: ~30 seconds
- **First 100 emails**: ~2-3 minutes  
- **First 1,000 emails**: ~20-30 minutes
- **All 84K+ emails**: ~8-10 hours

## Resume if Interrupted

```bash
go run main.go -all -workers 1 -resume \
  -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" \
  -job "JOB-RESUME" -job-name "Resume processing" \
  -query "" -batch 25 -delay 100 -v
```

## Files to Check

- **Log file**: `go/email-analyzer-verify.log` (or `go/email-analyzer-full.log`)
- **Google Sheet**: https://docs.google.com/spreadsheets/d/1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU
- **Dashboard**: http://localhost:8080/email-analysis-dashboard.html?spreadsheet_id=1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU

## Key Points

1. **Verification is built-in** - Each write is confirmed
2. **Row count is tracked** - Shows data accumulating
3. **Errors are logged** - Check log if something seems wrong
4. **Test mode available** - Use `-test` to verify before full run
5. **Resume supported** - Can continue if interrupted
