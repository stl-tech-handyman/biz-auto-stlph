# Email Analyzer: Verification and Testing Guide

## Overview

The email analyzer now includes comprehensive verification and logging to ensure data is being written to the Google Sheet correctly. This document explains how to verify everything is working.

## Test Mode

Run with `-test` flag to process only 10 emails and verify writes:

```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"

go run main.go -test -idempotent -spreadsheet "YOUR_SPREADSHEET_ID" \
  -job "JOB-TEST" -job-name "Test verification" \
  -query "" -batch 10 -workers 1
```

**What test mode does:**
1. ✅ Writes a test row to verify sheet access
2. ✅ Reads it back to confirm write succeeded
3. ✅ Processes 10 emails
4. ✅ Verifies each batch write
5. ✅ Shows total row count after each write

## Verification Features

### 1. Write Verification
- After each batch write, checks that the API returned a success response
- Logs the number of rows updated
- Shows timestamp for each write operation

### 2. Row Count Tracking
- After each batch, reads the sheet to get total row count
- Shows: `📊 Sheet now has X total data rows (excluding header)`
- This confirms data is accumulating in the sheet

### 3. Error Detection
- If write fails, shows detailed error message
- Restores batch on error (doesn't lose data)
- Logs batch details for debugging

## What to Look For

### ✅ Good Signs:
- `✅ Confirmed: X rows written to spreadsheet`
- `📊 Sheet now has X total data rows`
- `✅ Successfully wrote batch of X emails to spreadsheet`
- Row count increases after each batch

### ⚠️ Warning Signs:
- `⚠️ Warning: Expected X rows, but Y rows were updated` - May indicate partial write
- `❌ ERROR writing batch` - Write failed, check error message
- Row count not increasing - Data may not be writing

## Full Run with Verification

To process all emails with full verification:

```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"

go run main.go -all -workers 1 -idempotent \
  -spreadsheet "YOUR_SPREADSHEET_ID" \
  -job "JOB-FULL-VERIFY" \
  -job-name "Full analysis with verification" \
  -query "" -batch 25 -delay 100 -v \
  > ../../email-analyzer-verify.log 2>&1 &
```

## Monitoring Progress

### 1. Check Log File
```bash
tail -f go/email-analyzer-verify.log
```

Look for:
- `💾 Writing batch of X emails to spreadsheet...`
- `✅ Successfully wrote batch of X emails`
- `📊 Sheet now has X total data rows`

### 2. Check Google Sheet Directly
- Open the spreadsheet
- Go to "Raw Data" tab
- Check row count (should increase over time)
- Verify latest rows have data

### 3. Check Dashboard
- Open: `http://localhost:8080/email-analysis-dashboard.html?spreadsheet_id=YOUR_ID`
- Should show increasing "Total Processed" count
- If stuck at 0, check API logs

## Troubleshooting

### Sheet Not Updating

1. **Check if process is running:**
   ```bash
   ps aux | grep "[g]o run main.go"
   ```

2. **Check for errors in log:**
   ```bash
   grep -i "error\|failed\|warning" go/email-analyzer-verify.log | tail -20
   ```

3. **Verify sheet permissions:**
   - Service account must have edit access
   - Check GMAIL_CREDENTIALS_JSON is set correctly
   - Check GMAIL_FROM is set to team@stlpartyhelpers.com

4. **Test write capability:**
   ```bash
   go run main.go -test -spreadsheet "YOUR_ID" -query ""
   ```

### Row Count Not Increasing

- Check if batches are being written (look for "💾 Writing batch" messages)
- Verify no errors in log
- Check Google Sheet directly - sometimes API is slow to reflect changes
- Wait 30 seconds and check again

### Verification Warnings

- `⚠️ Warning: Written message ID X not found` - Usually a timing issue
- The write succeeded (you'll see "Successfully wrote")
- Check row count - if it's increasing, data is being written
- This warning is informational, not critical

## Expected Output

### Successful Run:
```
📋 Initializing sheets...
✅ Sheets initialized successfully
🚀 Starting email processing...
📦 Fetching batch 1...
  📦 Batch 1: Fetched 25 message IDs | Has next page: true
  [Worker 0] ✅ Processed: abc123...
  💾 Writing batch of 25 emails to spreadsheet...
  ✅ Confirmed: 25 rows written to spreadsheet
  📊 Sheet now has 25 total data rows (excluding header)
  ✅ Successfully wrote batch of 25 emails to spreadsheet
```

## Quick Verification Commands

```bash
# Check if process is running
ps aux | grep "[g]o run main.go"

# Check latest writes
tail -50 go/email-analyzer-verify.log | grep -E "(Writing|Successfully|Sheet now has)"

# Count total batches written
grep -c "Successfully wrote batch" go/email-analyzer-verify.log

# Check for errors
grep -i "error\|failed" go/email-analyzer-verify.log | tail -10
```

## Next Steps

Once verification shows data is writing correctly:
1. Monitor the log file for progress
2. Check Google Sheet periodically to confirm rows are accumulating
3. Use dashboard to see live stats
4. Let it run until all emails are processed
