# Email Analyzer (Local) — End-to-End Flow

This document describes the **actual runtime flow** of the local analyzer (`go/cmd/email-analyzer`) and how it achieves:
- fast processing via **goroutine workers**
- safety via **spreadsheet lock**
- resumability via **message-id based checkpoints**
- observability via the **Email Analysis Dashboard**

## High-level architecture

There is **one process** (one CLI run) and inside it:

1. **Producer (1 goroutine)**
   - pages through Gmail search results (`Users.Messages.List`)
   - pushes message IDs to a buffered channel

2. **Workers (`-workers N` goroutines)**
   - read message IDs from the channel
   - fetch full messages (`Users.Messages.Get`)
   - extract structured fields
   - append rows to an in-memory batch

3. **Batch writer (synchronized)**
   - only one worker at a time writes batches to Google Sheets
   - state is checkpointed after writes

## Why “agent” still exists

`-agent` identifies the **process** (not the goroutines). It is used for:
- lock identity in the `Locks` sheet
- job stats attribution in the `Job Stats` sheet
- dashboard reporting

## Locking model (process-level safety)

Before processing, the analyzer writes an **ACTIVE lock row** into the spreadsheet.

- **Goal**: prevent two separate processes from writing the same spreadsheet concurrently
- Locks expire after ~1 minute and are refreshed periodically while the process is alive
- If a process crashes, the lock expires and can be cleaned up

Important: the lock is about **multiple processes**, not goroutines. Goroutines are safe because they share memory and synchronize with mutexes.

## Resume model (idempotent message ID processing)

The analyzer uses Gmail’s **Message ID** as the unique “already processed” key.

- A processed message ID is recorded in:
  - `Raw Data` sheet (source-of-truth audit trail)
  - `State` sheet (`ProcessedIDs` list + counts)
- On `-resume`, the analyzer loads processed IDs and skips duplicates.

This makes the run **restart-safe** even if interrupted.

## Step-by-step flow

1. **Start**
   - parse flags (`-workers`, `-resume`, `-spreadsheet`, `-job`, etc.)
   - initialize Gmail + Sheets clients

2. **Spreadsheet**
   - if `-spreadsheet` empty: create a new spreadsheet
   - initialize headers/sheets (optionally clear data in `-idempotent`)

3. **Acquire lock**
   - write lock row (`ACTIVE`) with expiry
   - refresh periodically

4. **Load state**
   - if `-resume`: load processed message IDs and `LastIndex`
   - otherwise start fresh

5. **Producer goroutine**
   - Gmail search query returns message IDs in pages
   - pushes IDs to `messageChan`

6. **Worker goroutines**
   - read from `messageChan`
   - skip if message ID already processed
   - fetch full message
   - extract fields (client email, pricing hints, test/confirmation, etc.)
   - append row to batch

7. **Write batches + checkpoint**
   - when batch reaches threshold (25 rows):
     - append rows to `Raw Data`
     - save `State` with updated counts + processed IDs
     - update `Job Stats`

8. **Finish**
   - write remaining batch (if any)
   - save final state
   - release lock

## Dashboard flow

The dashboard reads from the spreadsheet via the API endpoint:

- `GET /api/email-analysis/stats?spreadsheet_id=...`

It renders:
- processed counts
- job history
- active agents (locks)
- direct Google Sheet link

## Recommended settings

- `-workers 3`: safest default for long runs
- `-workers 5`: faster, usually fine
- `-workers > 5`: only if you’re sure you won’t hit Gmail API throttling

## Where to look next

- Local analyzer docs: `go/cmd/email-analyzer/README.md`
- Launch instructions: `go/cmd/email-analyzer/LAUNCH.md`
- Concurrency & lock behavior: `go/cmd/email-analyzer/CONCURRENT_PROCESSING.md`
- Resume behavior: `go/cmd/email-analyzer/RESUME_GUIDE.md`
