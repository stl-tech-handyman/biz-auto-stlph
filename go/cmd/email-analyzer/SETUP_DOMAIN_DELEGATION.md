# Setup Domain-Wide Delegation for Email Analysis

## Current Error
```
unauthorized_client: Client is unauthorized to retrieve access tokens using this method
```

This means the service account needs domain-wide delegation configured.

## Steps to Fix

### 1. Get Service Account Client ID

1. Open Google Cloud Console: https://console.cloud.google.com/iam-admin/serviceaccounts
2. Find your service account (the one in `gmail-credentials.json`)
3. Click on it → **Details** tab
4. Copy the **Client ID** (looks like: `123456789012345678901`)

### 2. Enable Domain-Wide Delegation in Google Workspace

1. Go to Google Workspace Admin Console: https://admin.google.com
2. Navigate to: **Security** → **API Controls** → **Domain-wide Delegation**
3. Click **Add new**
4. Enter:
   - **Client ID**: (paste from step 1)
   - **OAuth Scopes** (one per line):
     ```
     https://www.googleapis.com/auth/gmail.readonly
     https://www.googleapis.com/auth/spreadsheets
     ```
5. Click **Authorize**

### 3. Wait 1-2 Minutes

Wait for changes to propagate.

### 4. Test Again

```bash
cd go/cmd/email-analyzer
export GMAIL_CREDENTIALS_JSON="../../.secrets/gmail-credentials.json"
export GMAIL_FROM="team@stlpartyhelpers.com"
go run main.go -max 5 -workers 3 -job "JOB-1-TEST" -job-name "Test" -v
```

## Required Scopes

Make sure these are added to domain-wide delegation:
- `https://www.googleapis.com/auth/gmail.readonly` (for reading emails)
- `https://www.googleapis.com/auth/spreadsheets` (for creating/writing to sheets)

## Alternative: Use OAuth2 Instead

If domain-wide delegation is not possible, you could use OAuth2 with a refresh token, but domain-wide delegation is recommended for service accounts.
