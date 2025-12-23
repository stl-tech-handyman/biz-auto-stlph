# Email Service Standard Operating Procedure (SOP)

## Purpose

This SOP documents the complete setup, configuration, and operation of the email service for BizOps360 Go API.

## Architecture

- **Email Service Project**: `bizops360-email-dev` (separate GCP project)
- **Main API Project**: `bizops360-dev` (uses email credentials from email project)
- **Implementation**: Gmail API integration directly in Go API
- **Credentials Storage**: Google Secret Manager in email service project

## Initial Setup

### Prerequisites

- GCP billing account: `01C379-C9A8C1-3ED059`
- Access to both projects: `bizops360-dev` and `bizops360-email-dev`
- Google Workspace admin access (for domain-wide delegation)

### Step 1: Create Email Service Project

```bash
# Create project
gcloud projects create bizops360-email-dev --name="BizOps360 Email Dev"

# Link billing account
# Via web: https://console.cloud.google.com/billing/01C379-C9A8C1-3ED059/linkedprojects
# Or CLI: gcloud beta billing projects link bizops360-email-dev --billing-account=01C379-C9A8C1-3ED059
```

### Step 2: Enable Required APIs

```bash
gcloud config set project bizops360-email-dev

gcloud services enable \
    gmail.googleapis.com \
    secretmanager.googleapis.com \
    iam.googleapis.com \
    cloudresourcemanager.googleapis.com
```

### Step 3: Create Service Account

```bash
gcloud iam service-accounts create email-service-sa \
    --display-name="Email Service Account" \
    --project=bizops360-email-dev
```

### Step 4: Generate Service Account Key

```bash
SERVICE_ACCOUNT="email-service-sa@bizops360-email-dev.iam.gserviceaccount.com"

gcloud iam service-accounts keys create /tmp/email-service-key.json \
    --iam-account="$SERVICE_ACCOUNT" \
    --project=bizops360-email-dev
```

### Step 5: Store Credentials in Secret Manager

```bash
gcloud secrets create gmail-credentials-json \
    --data-file=/tmp/email-service-key.json \
    --replication-policy="automatic" \
    --project=bizops360-email-dev
```

### Step 6: Grant Access to Main API Service Account

```bash
MAIN_PROJECT="bizops360-dev"
MAIN_PROJECT_NUMBER=$(gcloud projects describe $MAIN_PROJECT --format="value(projectNumber)")
MAIN_SA="${MAIN_PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

gcloud secrets add-iam-policy-binding gmail-credentials-json \
    --member="serviceAccount:${MAIN_SA}" \
    --role="roles/secretmanager.secretAccessor" \
    --project=bizops360-email-dev
```

### Step 7: Configure Gmail API

**Option A: Domain-Wide Delegation (Recommended for Workspace)**

1. Go to [Google Workspace Admin Console](https://admin.google.com)
2. Navigate to Security → API Controls → Domain-wide Delegation
3. Add new API client:
   - Client ID: Get from service account details
   - OAuth Scopes: `https://www.googleapis.com/auth/gmail.send`
4. Authorize

**Option B: OAuth2 Credentials (For personal Gmail)**

1. Go to [GCP Console Credentials](https://console.cloud.google.com/apis/credentials?project=bizops360-email-dev)
2. Create OAuth 2.0 Client ID
3. Download credentials JSON
4. Update secret with OAuth2 credentials

### Step 8: Deploy Main API with Email Support

```bash
cd go
bash scripts/deploy-dev.sh
```

The deployment script automatically includes Gmail credentials if they exist in the email service project.

## Testing

### Test Email Endpoint

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

### Expected Success Response

```json
{
  "ok": true,
  "message": "Email sent successfully",
  "result": {
    "messageId": "gmail-message-id",
    "success": true
  }
}
```

## Troubleshooting

### Issue: "GMAIL_CREDENTIALS_JSON environment variable is not set"

**Solution:**
1. Verify secret exists: `gcloud secrets list --project=bizops360-email-dev`
2. Check IAM permissions: `gcloud secrets get-iam-policy gmail-credentials-json --project=bizops360-email-dev`
3. Redeploy: `cd go && bash scripts/deploy-dev.sh`

### Issue: "failed to parse Gmail credentials"

**Solution:**
1. Verify credentials format: Should be valid JSON
2. Check secret content: `gcloud secrets versions access latest --secret=gmail-credentials-json --project=bizops360-email-dev`
3. Ensure correct credential type (service account vs OAuth2)

### Issue: "insufficient permissions" or "access denied"

**Solution:**
1. For service accounts: Enable domain-wide delegation in Google Workspace
2. Verify OAuth scopes include `https://www.googleapis.com/auth/gmail.send`
3. Check service account has necessary IAM roles

### Issue: "User not found" when using service account

**Solution:**
- Service accounts cannot impersonate users without domain-wide delegation
- Use OAuth2 credentials with user authorization, or
- Configure domain-wide delegation properly

## Maintenance

### Rotating Credentials

1. Generate new key:
```bash
gcloud iam service-accounts keys create /tmp/new-key.json \
    --iam-account=email-service-sa@bizops360-email-dev.iam.gserviceaccount.com \
    --project=bizops360-email-dev
```

2. Update secret:
```bash
gcloud secrets versions add gmail-credentials-json \
    --data-file=/tmp/new-key.json \
    --project=bizops360-email-dev
```

3. Redeploy API (will automatically use new version)

### Monitoring

- Check Cloud Run logs for email sending errors
- Monitor Gmail API quota usage
- Review Secret Manager access logs

### Backup

- Service account keys are stored in Secret Manager (automatic replication)
- Keep backup of service account key JSON files securely

## Security Considerations

1. **Secret Access**: Only main API service account has access to email credentials
2. **Least Privilege**: Service account has minimal required permissions
3. **Credential Rotation**: Rotate credentials periodically
4. **Audit Logging**: Monitor Secret Manager access logs
5. **Domain-Wide Delegation**: Limit to specific OAuth scopes

## Production Setup

For production (`bizops360-prod`):

1. Create production email project: `bizops360-email-prod`
2. Follow same setup steps
3. Use production Gmail account or workspace
4. Update production deployment script to use production credentials

## Related Documentation

- [Email Service Setup Guide](EMAIL_SERVICE_SETUP.md)
- [Gmail API Documentation](https://developers.google.com/gmail/api)
- [Domain-Wide Delegation Guide](https://developers.google.com/identity/protocols/oauth2/service-account#delegatingauthority)

