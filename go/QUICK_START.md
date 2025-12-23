# Quick Start - Deploy to bizops360-dev

## Prerequisites

- gcloud CLI installed and authenticated
- Docker installed and running
- Billing account available

## Step 1: Link Billing Account

```bash
# List billing accounts
gcloud beta billing accounts list

# Link billing (replace BILLING_ACCOUNT_ID with actual ID)
gcloud billing projects link bizops360-dev --billing-account=BILLING_ACCOUNT_ID
```

## Step 2: Complete Setup and Deploy

Run the automated setup script:

```bash
cd go
bash scripts/complete-setup-dev.sh
```

This script will:
1. ✅ Enable all required APIs
2. ✅ Create Artifact Registry repository
3. ✅ Configure Docker authentication
4. ✅ Create secrets (API key, Stripe keys)
5. ✅ Deploy the service

## Alternative: Manual Setup

If you prefer manual steps, see [SETUP_NEW_PROJECTS.md](SETUP_NEW_PROJECTS.md)

## After Deployment

Once deployed, you'll get a service URL like:
`https://bizops360-api-go-dev-XXXXX-uc.a.run.app`

Test it:
```bash
curl https://YOUR-SERVICE-URL/api/health
```

## Getting API Key

After deployment, get your API key:
```bash
gcloud secrets versions access latest --secret="svc-api-key-dev" --project="bizops360-dev"
```

## Troubleshooting

### Billing Not Linked
- Error: "Billing account not found"
- Solution: Link billing account (Step 1)

### APIs Not Enabled
- Error: "API not enabled"
- Solution: Run `bash scripts/complete-setup-dev.sh` to enable all APIs

### Docker Push Fails
- Error: "denied: Permission denied"
- Solution: Run `gcloud auth configure-docker us-central1-docker.pkg.dev`

### Secrets Not Found
- Error: "Secret not found"
- Solution: Create secrets using `bash scripts/create-secrets.sh`

