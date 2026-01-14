# 🔒 Concurrent Processing with Lock System

## ✅ Features

The analyzer now supports **safe concurrent processing** with:
- ✅ **Lock mechanism** - Prevents multiple agents from conflicting
- ✅ **Auto-expiration** - Locks expire after 1 minute (auto-cleanup)
- ✅ **Idempotent mode** - Can recreate sheets/clear data
- ✅ **Agent IDs** - Each process has unique identifier
- ✅ **Lock refresh** - Extends lock while processing

## 🚀 Running Multiple Agents

### Agent 1
```bash
go run main.go -all -agent "agent-1" -v
```

### Agent 2 (Different Query/Spreadsheet)
```bash
go run main.go -all -agent "agent-2" -query "from:zapier.com" -v
```

### Agent 3 (Different Spreadsheet)
```bash
go run main.go -all -agent "agent-3" -spreadsheet DIFFERENT_SPREADSHEET_ID -v
```

## 🔒 Lock System

### How It Works

1. **Acquire Lock**
   - Agent requests lock before processing
   - Lock expires in 1 minute
   - Only one active lock per spreadsheet

2. **Refresh Lock**
   - Lock refreshed every 30 seconds
   - Extends expiration to 1 minute from now
   - Keeps lock alive while processing

3. **Release Lock**
   - Released when process completes
   - Released on error/exit (defer)
   - Auto-expires after 1 minute if process dies

### Lock Sheet Structure

| Agent ID | Created At | Expires At | Status |
|----------|------------|------------|--------|
| agent-1 | 2025-01-15T10:00:00Z | 2025-01-15T10:01:00Z | ACTIVE |
| agent-2 | 2025-01-15T10:05:00Z | 2025-01-15T10:06:00Z | EXPIRED |

## 🧹 Auto-Cleanup

### Apps Script Cleanup Service

**Setup:**
1. Open Apps Script editor
2. Create new project
3. Copy `go/scripts/lock-cleanup.gs` code
4. Set spreadsheet ID in script properties
5. Set up trigger: Every minute

**Run:**
```javascript
// In Apps Script
setupLockCleanup();  // One-time setup
cleanupExpiredLocks(); // Runs automatically on schedule
```

**Manual cleanup:**
```javascript
manualCleanup(); // Run manually if needed
getLockStatus();  // Check current locks
```

## 🔄 Idempotent Mode

### Recreate Sheets

```bash
# Clear all data and recreate sheets
go run main.go -all -idempotent -v
```

**What it does:**
- Clears all existing locks
- Clears all data in sheets (except State/Locks)
- Recreates headers
- Starts fresh processing

## 📊 Concurrent Processing Strategies

### Strategy 1: Different Spreadsheets
```bash
# Agent 1: Process to spreadsheet 1
go run main.go -all -agent "agent-1" -spreadsheet SPREADSHEET_1 -v

# Agent 2: Process to spreadsheet 2
go run main.go -all -agent "agent-2" -spreadsheet SPREADSHEET_2 -v
```
✅ **No conflicts** - Different spreadsheets

### Strategy 2: Different Queries (Same Spreadsheet)
```bash
# Agent 1: Process zapier emails
go run main.go -max 10000 -agent "agent-1" -query "from:zapier.com" -spreadsheet ID -v

# Agent 2: Process form submissions (waits for lock)
go run main.go -max 10000 -agent "agent-2" -query "subject:\"Form Submission\"" -spreadsheet ID -v
```
⚠️ **Sequential** - Second agent waits for first to finish

### Strategy 3: Chunked Processing
```bash
# Agent 1: Process first 20K
go run main.go -max 20000 -agent "agent-1" -spreadsheet ID -v

# Agent 2: Process next 20K (after agent 1 finishes)
go run main.go -max 20000 -agent "agent-2" -resume -spreadsheet ID -v
```
✅ **Sequential chunks** - Each agent processes different range

## 🔍 Monitoring Locks

### Check Lock Status (Apps Script)
```javascript
getLockStatus();
```

### Check Lock Status (Go)
```bash
# Verbose mode shows lock activity
go run main.go -max 100 -v
# Shows: 🔒 Lock acquired, 🔒 Lock refreshed, 🔓 Lock released
```

### Check Locks Sheet
Open spreadsheet → Locks sheet:
- See all active/expired locks
- Check expiration times
- Monitor agent activity

## ⚙️ Command Options

```
-agent string
    Agent ID for concurrent processing (auto-generated if empty)

-idempotent
    Recreate sheets if they exist (idempotent mode)

-resume
    Resume from last position (works with locks)
```

## 🛡️ Safety Features

### ✅ Lock Expiration
- Locks expire after 1 minute
- Auto-refreshed every 30 seconds
- Prevents deadlocks from crashed processes

### ✅ Cleanup Service
- Apps Script runs every minute
- Removes expired locks automatically
- Prevents lock buildup

### ✅ Agent Identification
- Each process has unique agent ID
- Tracks which agent holds lock
- Helps debug concurrent issues

### ✅ Idempotent Mode
- Can safely recreate sheets
- Clears old data
- Starts fresh processing

## 📝 Example Workflow

### Setup Cleanup Service
```javascript
// In Apps Script
setupLockCleanup();
// Set trigger: cleanupExpiredLocks, Every minute
```

### Run Multiple Agents
```bash
# Terminal 1
go run main.go -all -agent "agent-1" -v

# Terminal 2 (waits for lock)
go run main.go -all -agent "agent-2" -v
# Will show: ❌ Lock already held by another agent
```

### Resume After Lock Release
```bash
# After agent-1 finishes, agent-2 can run
go run main.go -all -agent "agent-2" -resume -spreadsheet ID -v
```

## 🎯 Best Practices

1. **Use Different Spreadsheets** for true parallelism
2. **Use Different Queries** for sequential processing
3. **Monitor Locks** - Check Locks sheet regularly
4. **Setup Cleanup** - Run Apps Script cleanup service
5. **Use Agent IDs** - Identify which process is running

## 🚨 Troubleshooting

### "Lock already held"
- Another agent is processing
- Wait for it to finish
- Or cleanup expired locks manually

### "Lock not found"
- Lock expired or was cleaned up
- Process may have crashed
- Just restart with new agent ID

### Stuck Locks
- Run cleanup service manually
- Or use idempotent mode to clear all

## ✅ Ready for Concurrent Processing!

You can now safely run multiple agents:
- ✅ Different spreadsheets = true parallelism
- ✅ Same spreadsheet = sequential (with locks)
- ✅ Auto-cleanup = no stuck locks
- ✅ Idempotent mode = fresh start anytime

**Just use `-agent` flag to identify each process!** 🚀
