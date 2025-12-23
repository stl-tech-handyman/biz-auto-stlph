# Deployment Instructions

## Quick Deploy Commands (Bash)

### Step 1: Authenticate with gcloud

```bash
# Check current authentication
gcloud auth list

# If needed, login (will open browser)
gcloud auth login

# Set application-default credentials (needed for Docker)
gcloud auth application-default login
```

### Step 2: Configure Docker for Google Container Registry

```bash
# Configure Docker to use gcloud credentials
gcloud auth configure-docker
```

### Step 3: Set the Project (if not already set)

```bash
# For dev deployment
gcloud config set project bizops360-dev

# For prod deployment
gcloud config set project bizops360-prod
```

### Step 4: Deploy

```bash
cd go

# Deploy to DEV (recommended first)
./scripts/deploy-dev.sh

# OR deploy to PROD (requires confirmation)
./scripts/deploy-prod.sh

# OR deploy to both
./scripts/deploy-all.sh
```

## One-Line Commands

### Deploy to Dev:
```bash
cd go && gcloud config set project bizops360-dev && ./scripts/deploy-dev.sh
```

### Deploy to Prod:
```bash
cd go && gcloud config set project bizops360-prod && ./scripts/deploy-prod.sh
```

## Troubleshooting

### If authentication fails:
```bash
# Re-authenticate
gcloud auth login
gcloud auth application-default login
```

### If Docker push fails:
```bash
# Re-configure Docker
gcloud auth configure-docker
```

### If secrets are missing:
```bash
# Check existing secrets
gcloud secrets list --project=bizops360-dev

# Create missing secrets (example)
echo -n "your-api-key" | gcloud secrets create svc-api-key-dev \
  --data-file=- --project=bizops360-dev
```


