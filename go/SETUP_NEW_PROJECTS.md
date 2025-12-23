# Setting Up New GCP Projects

## Quick Start

### Step 1: Link Billing Account

First, link a billing account to the projects:

```bash
# List billing accounts
gcloud beta billing accounts list

# Link to dev project
gcloud billing projects link bizops360-dev --billing-account=BILLING_ACCOUNT_ID

# Link to prod project (optional for now)
gcloud billing projects link bizops360-prod --billing-account=BILLING_ACCOUNT_ID
```

### Step 2: Complete Setup and Deploy

Run the complete setup script which will:
- Enable all required APIs
- Create Artifact Registry repository
- Create secrets
- Optionally deploy

```bash
cd go
bash scripts/complete-setup-dev.sh
```

## Manual Steps

If you prefer to do it manually:

### 1. Enable APIs

```bash
gcloud config set project bizops360-dev

gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable secretmanager.googleapis.com
```

### 2. Create Artifact Registry

```bash
gcloud artifacts repositories create bizops360-dev \
    --repository-format=docker \
    --location=us-central1 \
    --description="Docker repository for bizops360-dev"
```

### 3. Configure Docker Auth

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

### 4. Create Secrets

```bash
# API Key
echo -n "your-api-key-here" | gcloud secrets create svc-api-key-dev \
    --data-file=- \
    --replication-policy="automatic"

# Stripe Test (optional)
echo -n "sk_test_..." | gcloud secrets create stripe-secret-key-test \
    --data-file=- \
    --replication-policy="automatic"

# Stripe Prod (optional)
echo -n "sk_live_..." | gcloud secrets create stripe-secret-key-prod \
    --data-file=- \
    --replication-policy="automatic"
```

### 5. Deploy

```bash
cd go
bash scripts/deploy-dev.sh
```

## Scripts Available

- `scripts/create-projects.sh` - Create projects (already done)
- `scripts/link-billing.sh` - Link billing accounts
- `scripts/complete-setup-dev.sh` - Complete setup for dev
- `scripts/create-secrets.sh` - Create secrets only
- `scripts/deploy-dev.sh` - Deploy to dev

