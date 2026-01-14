# 🚀 Process All 86,000 Emails

## ✅ Ready for Full Processing!

The analyzer is now enhanced to handle **all your emails** efficiently with:
- ✅ Progress tracking
- ✅ Auto-save every 5 minutes
- ✅ Resume functionality
- ✅ Rate limiting
- ✅ Detailed statistics

## Quick Start - Process Everything

### Option 1: Process All at Once (Recommended)
```bash
cd go/cmd/email-analyzer
go run main.go -all -v
```

This will:
- Process ALL available emails
- Show real-time progress
- Auto-save state every 5 minutes
- Give you complete stats when done

### Option 2: Process in Chunks (Safer)
```bash
# First: 10,000 emails
go run main.go -max 10000 -v

# Then resume and continue
go run main.go -max 10000 -resume -spreadsheet YOUR_ID -v

# Repeat until all processed
```

## Command Options

```
-all              Process all available emails (no limit)
-max N            Process N emails (0 = all)
-batch N          Batch size (default: 50)
-delay N          Delay between batches in ms (default: 200)
-resume           Resume from last position
-spreadsheet ID   Use existing spreadsheet
-query "..."      Custom Gmail search query
-v                Verbose logging
```

## Examples

### Process All (Fast)
```bash
go run main.go -all -batch 100 -delay 100 -v
```
- Larger batches, shorter delays
- ~6-8 hours for 86K emails

### Process All (Safe)
```bash
go run main.go -all -batch 50 -delay 500 -v
```
- Smaller batches, longer delays
- ~12-15 hours, respects quota better

### Process 20,000 at a Time
```bash
go run main.go -max 20000 -v
```

### Resume Processing
```bash
go run main.go -all -resume -spreadsheet YOUR_SPREADSHEET_ID -v
```

## What You'll See

### Progress Updates
```
📦 Batch 1 | Processed: 45 | Skipped: 5 | Elapsed: 30s
📦 Batch 10 | Processed: 450 | Skipped: 50 | Elapsed: 5m
  📊 Progress: 1000 processed | 100 skipped | 3.2 emails/sec
  💾 Auto-saved state
```

### Final Summary
```
============================================================
✅ PROCESSING COMPLETE!
============================================================
  📧 Processed: 86000 emails
  ⏭️  Skipped: 5000 emails
  📊 Total processed (all time): 86000
  ⏱️  Time elapsed: 7h 30m
  🚀 Processing rate: 3.2 emails/second
  ⏳ Average time per email: 315ms
  📍 Next index: 91000
  📄 Spreadsheet: https://docs.google.com/spreadsheets/d/...
============================================================
```

## Processing Strategy

### For 86,000 Emails:

**Recommended Approach:**
```bash
# Start with a test run
go run main.go -max 1000 -v

# Review results, then process all
go run main.go -all -v
```

**Time Estimates:**
- **Fast mode**: ~6-8 hours (batch 100, delay 100ms)
- **Safe mode**: ~12-15 hours (batch 50, delay 500ms)
- **Chunked**: Process 10K-20K at a time, resume between runs

## Features

### ✅ Auto-Save
- State saved every 5 minutes
- Can interrupt and resume anytime
- Progress never lost

### ✅ Progress Tracking
- Real-time processing rate
- Batch progress updates
- Time estimates

### ✅ Email Migration
- Automatically handles `stlpartyhelpers@gmail.com` → `team@stlpartyhelpers.com`
- Detects forwarded conversations
- Tracks email mapping

### ✅ Resume Support
- Interrupt anytime (Ctrl+C)
- Resume with `-resume` flag
- Continues from last position

## Monitoring

### Watch Progress
1. **Console output** - Real-time updates
2. **Spreadsheet State sheet** - Shows current progress
3. **Verbose mode** - Detailed logging with `-v`

### Check Quota
- Monitor GCP Console for quota usage
- Safe processing: ~200K emails/day
- 86K emails = well within daily limit

## Troubleshooting

### "Quota Exceeded"
- Wait for daily reset (midnight PST)
- Reduce batch size: `-batch 25`
- Increase delay: `-delay 1000`

### Process Interrupted
- Just resume: `go run main.go -all -resume -spreadsheet YOUR_ID -v`
- Will continue from last saved position

### Slow Processing
- Increase batch size: `-batch 100`
- Decrease delay: `-delay 100`
- Check network connection

## Ready to Process All?

```bash
# Just run this!
cd go/cmd/email-analyzer
go run main.go -all -v
```

Let it run and it will process all 86,000 emails! 🎉

The analyzer will:
- ✅ Process all emails
- ✅ Show progress
- ✅ Auto-save state
- ✅ Handle email migration
- ✅ Give you complete stats

**Estimated time: 6-15 hours depending on settings**
