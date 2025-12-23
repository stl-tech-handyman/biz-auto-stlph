# Email Service Setup Guide

## Overview

Email functionality is integrated directly into the Go API using Gmail API. A separate GCP project (`bizops360-email-dev`) is used to store Gmail credentials.

## Prerequisites

- GCP project `bizops360-email-dev` created
- Billing account linked
- Gmail API enabled
- Service account created

## Step 1: Link Billing Account

Link billing account to the email service project:

```bash
# Via web console:
https://console.cloud.google.com/billing/01C379-C9A8C1-3ED059/linkedprojects

# Or via CLI (if beta component installed):
gcloud beta billing projects link bizops360-email-dev --billing-account=01C379-C9A8C1-3ED059
```

## Step 2: Enable Required APIs

```bash
gcloud config set project bizops360-email-dev

gcloud services enable gmail.googleapis.com
gcloud services enable secretmanager.googleapis.com
gcloud services enable iam.googleapis.com
```

## Step 3: Create Service Account and Credentials

### Option A: OAuth2 Credentials (Recommended for Gmail API)

1. Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials?project=bizops360-email-dev)
2. Click "Create Credentials" → "OAuth client ID"
3. Application type: "Web application"
4. Authorized redirect URIs: `http://localhost:8080` (for testing)
5. Download JSON credentials
6. Save to Secret Manager:

```bash
gcloud secrets create gmail-oauth-credentials-json \
    --data-file=path/to/credentials.json \
    --replication-policy="automatic" \
    --project=bizops360-email-dev
```

### Option B: Service Account with Domain-Wide Delegation

1. Create service account (already done):
```bash
gcloud iam service-accounts create email-service-sa \
    --display-name="Email Service Account" \
    --project=bizops360-email-dev
```

2. Enable domain-wide delegation in Google Workspace Admin Console
3. Create key:
```bash
gcloud iam service-accounts keys create /tmp/email-key.json \
    --iam-account=email-service-sa@bizops360-email-dev.iam.gserviceaccount.com \
    --project=bizops360-email-dev
```

4. Save to Secret Manager:
```bash
gcloud secrets create gmail-credentials-json \
    --data-file=/tmp/email-key.json \
    --replication-policy="automatic" \
    --project=bizops360-email-dev
```

## Step 4: Grant Access to Main API Service Account

Grant the main API service account access to read the secret:

```bash
MAIN_PROJECT="bizops360-dev"
MAIN_PROJECT_NUMBER=$(gcloud projects describe $MAIN_PROJECT --format="value(projectNumber)")
MAIN_SA="${MAIN_PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

gcloud secrets add-iam-policy-binding gmail-credentials-json \
    --member="serviceAccount:${MAIN_SA}" \
    --role="roles/secretmanager.secretAccessor" \
    --project=bizops360-email-dev
```

## Step 5: Update Deployment

The deployment script automatically includes Gmail credentials if they exist:

```bash
cd go
bash scripts/deploy-dev.sh
```

The script will add:
```
GMAIL_CREDENTIALS_JSON=projects/bizops360-email-dev/secrets/gmail-credentials-json:latest
```

## Step 6: Test Email Sending

```bash
curl -X POST 'https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app/api/email/test' \
  -H 'X-Api-Key: dev-api-key-12345' \
  -H 'Content-Type: application/json' \
  -d '{
    "to": "test@example.com",
    "subject": "Test Email",
    "html": "<p>This is a test email</p>"
  }'
```

## Troubleshooting

### Error: "GMAIL_CREDENTIALS_JSON environment variable is not set"
- Check that secret exists in `bizops360-email-dev` project
- Verify main API service account has access to the secret
- Redeploy the service

### Error: "failed to parse Gmail credentials"
- Verify credentials JSON is valid
- Check that credentials are for OAuth2 or service account with domain-wide delegation

### Error: "insufficient permissions"
- For service accounts: Enable domain-wide delegation in Google Workspace
- For OAuth2: Complete OAuth flow to get refresh token

## Gmail API Setup Notes

**Important:** Gmail API requires one of the following:

1. **OAuth2 Credentials** - User must authorize the application
2. **Service Account with Domain-Wide Delegation** - Requires Google Workspace admin setup
3. **App Passwords** - For personal Gmail accounts (not recommended for production)

For production, use OAuth2 with refresh tokens or domain-wide delegation.

