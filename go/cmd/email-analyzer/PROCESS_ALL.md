# Processing All 86,000 Emails

## Strategy for Full Processing

The local analyzer can process all your emails efficiently. Here's how:

## Quick Commands

### Process All Emails (Recommended)
```bash
cd go/cmd/email-analyzer
go run main.go -all -workers 5 -v
```

### Process with Custom Settings
```bash
# Process all with larger batches and shorter delays
go run main.go -all -batch 100 -delay 100 -v

# Process 10,000 at a time (resume between runs)
go run main.go -max 10000 -v
```

## Processing Options

### Option 1: Process All at Once
```bash
go run main.go -all -workers 5 -v
```
- Processes ALL available emails
- Auto-saves state every 2 minutes
- Can be interrupted and resumed
- Best for: Let it run overnight

### Option 2: Process in Chunks
```bash
# First run: 10,000 emails
go run main.go -max 10000 -workers 5 -v

# Second run: Resume and process next 10,000
go run main.go -max 10000 -workers 5 -resume -spreadsheet YOUR_ID -v

# Continue until done...
```
- Better for: Monitoring progress
- Can stop/start as needed
- Check results between runs

### Option 3: Process with Rate Limiting
```bash
# Slower but safer (respects quota)
go run main.go -all -workers 3 -batch 50 -delay 500 -v
```
- Larger delays between batches
- Safer for quota limits
- Takes longer but more reliable

## Progress Tracking

The analyzer shows:
- **Real-time progress**: Emails processed per second
- **Batch updates**: Every 10 batches
- **Auto-saves**: State saved every 2 minutes
- **Final summary**: Complete stats when done

### Example Output
```
📦 Batch 1 | Processed: 45 | Skipped: 5 | Elapsed: 30s
📦 Batch 10 | Processed: 450 | Skipped: 50 | Elapsed: 5m
  📊 Progress: 1000 processed | 100 skipped | 3.2 emails/sec
  💾 Auto-saved state
...
✅ PROCESSING COMPLETE!
  📧 Processed: 86000 emails
  ⏱️  Time elapsed: 7h 30m
  🚀 Processing rate: 3.2 emails/second
```

## Resume Functionality

If interrupted, resume with:
```bash
go run main.go -all -workers 5 -resume -spreadsheet YOUR_SPREADSHEET_ID -v
```

The analyzer will:
- Load state from spreadsheet
- Continue from last position
- Skip already processed emails

## Rate Limits & Quota

### Gmail API Limits
- **Daily quota**: 1,000,000 quota units
- **Per message read**: ~5 quota units
- **Safe processing**: ~200,000 emails/day

### Recommendations

**For 86,000 emails:**
- **Fast**: `-all -batch 100 -delay 100` (takes ~6-8 hours)
- **Safe**: `-all -batch 50 -delay 500` (takes ~12-15 hours, respects quota)
- **Chunked**: Process 10K at a time, resume between runs

## Monitoring

### Check Progress
1. **Watch console output** - Real-time progress
2. **Check spreadsheet** - State sheet shows progress
3. **Check logs** - Verbose mode shows details

### Estimate Time Remaining
```
If processing at 3 emails/sec:
- 10,000 emails = ~55 minutes
- 50,000 emails = ~4.5 hours  
- 86,000 emails = ~8 hours
```

## Best Practices

1. **Start Small**: Test with 100-1000 emails first
2. **Monitor First Run**: Watch for errors
3. **Use Resume**: Process in chunks if needed
4. **Check Quota**: Monitor GCP Console for quota usage
5. **Save State**: Auto-saves every 2 minutes (and after batch writes)
6. **Workers**: Prefer 3–5 workers for long runs

## Troubleshooting

### "Quota Exceeded"
- Wait for daily reset
- Reduce batch size: `-batch 25`
- Increase delay: `-delay 1000`

### "Too Many Requests"
- Increase delay: `-delay 1000`
- Reduce batch size: `-batch 25`

### Process Interrupted
- Use `-resume` flag
- Provide `-spreadsheet` ID
- Will continue from last position

## Example: Full Processing Session

```bash
# Day 1: Test run
go run main.go -max 1000 -v
# Review results, check spreadsheet

# Day 2: Process 20,000
go run main.go -max 20000 -v

# Day 3: Process remaining (resume)
go run main.go -all -resume -spreadsheet YOUR_ID -v
```

## Ready to Process All?

```bash
# Just run this and let it work!
cd go/cmd/email-analyzer
go run main.go -all -v
```

It will process all 86,000 emails, show progress, auto-save state, and give you a complete summary when done! 🚀
