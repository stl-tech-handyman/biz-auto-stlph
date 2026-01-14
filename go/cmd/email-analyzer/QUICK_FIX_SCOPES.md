# Quick Fix: Add Sheets Scope to Domain-Wide Delegation

## Current Issue
Your client ID `111386381202667840149` already exists but only has:
- ✅ `https://www.googleapis.com/auth/gmail.readonly`

## What You Need
Add BOTH scopes (one per line in the OAuth scopes field):

```
https://www.googleapis.com/auth/gmail.readonly
https://www.googleapis.com/auth/spreadsheets
```

## Steps
1. In the modal dialog, update the "OAuth scopes" field
2. Add the Sheets scope (second line)
3. Keep "Overwrite existing client ID" checked ✅
4. Click "AUTHORIZE"
5. Wait 1-2 minutes
6. Test again!

## Test Command
```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"
go run main.go -max 5 -workers 3 -job "JOB-1-TEST" -v
```
