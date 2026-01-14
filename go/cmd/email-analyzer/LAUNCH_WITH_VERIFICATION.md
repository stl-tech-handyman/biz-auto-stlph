# Launch Email Analyzer with Full Verification

## Quick Start - Process ALL Emails with Verification

### Step 1: Set Environment Variables
```bash
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"
```

### Step 2: Run Test Mode First (Verify Writes Work)
```bash
cd go/cmd/email-analyzer

go run main.go -test -idempotent \
  -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" \
  -job "JOB-TEST-VERIFY" \
  -job-name "Test: verify writes work" \
  -query "" -batch 10 -workers 1
```

**Expected output:**
- ✅ TEST PASSED: Successfully wrote test row
- ✅ TEST PASSED: Successfully verified read-back
- Processes 10 emails
- Shows verification for each write

### Step 3: Launch Full Processing
```bash
cd go/cmd/email-analyzer

nohup go run main.go -all -workers 1 -idempotent \
  -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" \
  -job "JOB-FULL-ALL-EMAILS" \
  -job-name "Process ALL emails - most recent to archive" \
  -query "" -batch 25 -delay 100 -v \
  > ../../email-analyzer-full.log 2>&1 &

echo "Started (PID: $!)"
```

## What to Monitor

### 1. Check Log File (Every Few Minutes)
```bash
tail -50 go/email-analyzer-full.log
```

**Look for:**
- `💾 Writing batch of X emails to spreadsheet...`
- `✅ Confirmed: X rows written to spreadsheet`
- `📊 Sheet now has X total data rows`
- `✅ Successfully wrote batch of X emails`

**Bad signs:**
- `❌ ERROR writing batch` - Write failed
- `⚠️ Warning: Expected X rows, but Y rows were updated` - Partial write

### 2. Check Google Sheet Directly
1. Open: https://docs.google.com/spreadsheets/d/1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU
2. Go to "Raw Data" tab
3. Check row count (should increase over time)
4. Scroll to bottom - should see new rows being added

### 3. Check Process Status
```bash
# Is it running?
ps aux | grep "[g]o run main.go"

# How many batches written?
grep -c "Successfully wrote batch" go/email-analyzer-full.log

# Latest progress
tail -20 go/email-analyzer-full.log | grep -E "(Progress|Wrote|Sheet now has)"
```

## Verification Checklist

After 5 minutes, verify:

- [ ] Process is still running (`ps aux | grep "[g]o run main.go"`)
- [ ] Log shows "Successfully wrote batch" messages
- [ ] Log shows "Sheet now has X total data rows" increasing
- [ ] Google Sheet "Raw Data" tab has new rows
- [ ] Row count in sheet matches or is close to log count
- [ ] No ERROR messages in log

## If Sheet Not Updating

### Check 1: Process Status
```bash
ps aux | grep "[g]o run main.go"
```
If not running, restart with command above.

### Check 2: Errors in Log
```bash
grep -i "error\|failed" go/email-analyzer-full.log | tail -20
```

### Check 3: Test Write Capability
```bash
go run main.go -test -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" -query ""
```

### Check 4: Verify Credentials
```bash
echo $GMAIL_CREDENTIALS_JSON
echo $GMAIL_FROM
```

## Expected Timeline

- **Test mode (10 emails)**: ~10-20 seconds
- **First 100 emails**: ~2-3 minutes
- **First 1,000 emails**: ~20-30 minutes
- **All 84K+ emails**: ~8-10 hours (at ~2-3 emails/sec)

## Key Verification Messages

### ✅ Good (Data is Writing):
```
✅ Confirmed: 25 rows written to spreadsheet
📊 Sheet now has 125 total data rows (excluding header)
✅ Successfully wrote batch of 25 emails to spreadsheet
```

### ⚠️ Warning (May be OK):
```
⚠️ Warning: Written message ID X not found in sheet verification
```
- This is usually a timing issue
- If you see "Successfully wrote" and row count increases, it's working

### ❌ Error (Needs Fix):
```
❌ ERROR writing batch: [error message]
```
- Check error message
- Verify credentials
- Check sheet permissions

## Resume if Interrupted

If process stops, resume with:
```bash
go run main.go -all -workers 1 -resume \
  -spreadsheet "1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU" \
  -job "JOB-RESUME" -job-name "Resume processing" \
  -query "" -batch 25 -delay 100 -v
```

## Full Command Reference

```bash
# Test mode (10 emails, verify writes)
go run main.go -test -idempotent -spreadsheet "ID" -query "" -workers 1

# Full run (all emails, most recent first)
go run main.go -all -workers 1 -idempotent \
  -spreadsheet "ID" -job "JOB-NAME" -job-name "Description" \
  -query "" -batch 25 -delay 100 -v

# Resume (continue from where left off)
go run main.go -all -workers 1 -resume \
  -spreadsheet "ID" -job "JOB-RESUME" -job-name "Resume" \
  -query "" -batch 25 -delay 100 -v
```

## Current Status

**Spreadsheet ID:** `1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU`

**Dashboard:** http://localhost:8080/email-analysis-dashboard.html?spreadsheet_id=1AqrivHOOF1HxjW2t_1XnGUhN4lwcdSv18VacyH36HeU

**Log File:** `go/email-analyzer-full.log`
