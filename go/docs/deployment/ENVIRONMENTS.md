# Environment Configuration Guide

## Overview

The Go API supports separate **dev** and **prod** environments, matching the JavaScript API setup. Each environment has its own:

- GCP Project (`bizops360-dev` / `bizops360-prod`)
- Cloud Run Service (`bizops360-api-go-dev` / `bizops360-api-go-prod`)
- Secret Manager secrets
- Configuration settings
- Resource limits

## Environment Structure

```
go/
├── config/
│   └── environments/
│       ├── dev.yaml      # Dev environment config
│       └── prod.yaml     # Prod environment config
├── scripts/
│   ├── deploy-dev.sh     # Deploy to dev
│   ├── deploy-prod.sh    # Deploy to prod
│   ├── deploy-all.sh     # Deploy to both
│   ├── build-dev.sh      # Build dev image
│   ├── build-prod.sh     # Build prod image
│   ├── test-dev.sh       # Test dev deployment
│   └── test-prod.sh      # Test prod deployment
├── Dockerfile.dev        # Dev Dockerfile (with debug tools)
└── Dockerfile.prod       # Prod Dockerfile (optimized)
```

## Environment Differences

### Development (`dev`)

- **Project**: `bizops360-dev`
- **Service**: `bizops360-api-go-dev`
- **Log Level**: `debug`
- **Resources**: 
  - Memory: 512Mi
  - CPU: 1
  - Min instances: 0 (scale to zero)
  - Max instances: 5
- **Features**:
  - Detailed logging enabled
  - Debug endpoints enabled
  - Uses Stripe test keys by default
- **Dockerfile**: `Dockerfile.dev` (includes debugging tools)

### Production (`prod`)

- **Project**: `bizops360-prod`
- **Service**: `bizops360-api-go-prod`
- **Log Level**: `info`
- **Resources**:
  - Memory: 1Gi
  - CPU: 2
  - Min instances: 1 (always warm)
  - Max instances: 20
- **Features**:
  - Optimized logging
  - Debug endpoints disabled
  - Uses Stripe live keys
- **Dockerfile**: `Dockerfile.prod` (distroless, optimized)

## Deployment

### Prerequisites

1. **GCP Projects**: Ensure both projects exist
   ```bash
   gcloud projects list
   ```

2. **Secrets**: Create required secrets in each project
   ```bash
   # Dev secrets
   gcloud secrets create svc-api-key-dev --project=bizops360-dev
   gcloud secrets create stripe-secret-key-test --project=bizops360-dev
   
   # Prod secrets
   gcloud secrets create svc-api-key-prod --project=bizops360-prod
   gcloud secrets create stripe-secret-key-prod --project=bizops360-prod
   ```

3. **Artifact Registry**: Enable in both projects
   ```bash
   gcloud services enable artifactregistry.googleapis.com --project=bizops360-dev
   gcloud services enable artifactregistry.googleapis.com --project=bizops360-prod
   ```

### Deploy to Development

```bash
cd go
./scripts/deploy-dev.sh
```

This will:
1. Build Docker image for dev
2. Push to Artifact Registry (`gcr.io/bizops360-dev/bizops360-api-go-dev`)
3. Deploy to Cloud Run in `bizops360-dev` project
4. Configure secrets and environment variables

### Deploy to Production

```bash
cd go
./scripts/deploy-prod.sh
```

**⚠️ Warning**: This script requires confirmation before deploying to production.

### Deploy to Both

```bash
cd go
./scripts/deploy-all.sh
```

Deploys to dev first, waits 10 seconds, then deploys to prod.

## Building Images

### Build Dev Image

```bash
cd go
./scripts/build-dev.sh
```

Builds using `Dockerfile.dev` with debugging symbols.

### Build Prod Image

```bash
cd go
./scripts/build-prod.sh
```

Builds using `Dockerfile.prod` optimized for production.

## Testing

### Test Dev Environment

```bash
cd go
./scripts/test-dev.sh
```

Tests:
- Health endpoint (no auth)
- Estimate endpoint (with API key)
- Deposit calculation (with API key)

### Test Prod Environment

```bash
cd go
./scripts/test-prod.sh
```

**⚠️ Warning**: Requires confirmation before testing production.

## Environment Variables

### Development

```bash
ENV=dev
LOG_LEVEL=debug
PORT=8080
CONFIG_DIR=/app/config
TEMPLATES_DIR=/app/templates
SERVICE_API_KEY=<from Secret Manager>
STRIPE_SECRET_KEY_TEST=<from Secret Manager>
```

### Production

```bash
ENV=prod
LOG_LEVEL=info
PORT=8080
CONFIG_DIR=/app/config
TEMPLATES_DIR=/app/templates
SERVICE_API_KEY=<from Secret Manager>
STRIPE_SECRET_KEY_PROD=<from Secret Manager>
STRIPE_SECRET_KEY_TEST=<from Secret Manager>
```

## Secret Manager

### Required Secrets

**Development (`bizops360-dev`)**:
- `svc-api-key-dev` - API key for authentication
- `stripe-secret-key-test` - Stripe test API key

**Production (`bizops360-prod`)**:
- `svc-api-key-prod` - API key for authentication
- `stripe-secret-key-prod` - Stripe live API key
- `stripe-secret-key-test` - Stripe test API key (optional)

### Creating Secrets

```bash
# Dev API key
echo -n "your-dev-api-key" | gcloud secrets create svc-api-key-dev \
  --data-file=- --project=bizops360-dev

# Prod API key
echo -n "your-prod-api-key" | gcloud secrets create svc-api-key-prod \
  --data-file=- --project=bizops360-prod

# Stripe keys
echo -n "sk_test_..." | gcloud secrets create stripe-secret-key-test \
  --data-file=- --project=bizops360-dev

echo -n "sk_live_..." | gcloud secrets create stripe-secret-key-prod \
  --data-file=- --project=bizops360-prod
```

## Service URLs

After deployment, get service URLs:

```bash
# Dev URL
gcloud run services describe bizops360-api-go-dev \
  --project=bizops360-dev \
  --region=us-central1 \
  --format="value(status.url)"

# Prod URL
gcloud run services describe bizops360-api-go-prod \
  --project=bizops360-prod \
  --region=us-central1 \
  --format="value(status.url)"
```

## Monitoring

### View Logs

```bash
# Dev logs
gcloud logging read "resource.type=cloud_run_revision AND \
  resource.labels.service_name=bizops360-api-go-dev" \
  --project=bizops360-dev --limit=50

# Prod logs
gcloud logging read "resource.type=cloud_run_revision AND \
  resource.labels.service_name=bizops360-api-go-prod" \
  --project=bizops360-prod --limit=50
```

### Check Service Status

```bash
# Dev status
gcloud run services describe bizops360-api-go-dev \
  --project=bizops360-dev \
  --region=us-central1

# Prod status
gcloud run services describe bizops360-api-go-prod \
  --project=bizops360-prod \
  --region=us-central1
```

## Rollback

To rollback to a previous revision:

```bash
# List revisions
gcloud run revisions list \
  --service=bizops360-api-go-prod \
  --project=bizops360-prod \
  --region=us-central1

# Rollback to specific revision
gcloud run services update-traffic bizops360-api-go-prod \
  --to-revisions=REVISION_NAME=100 \
  --project=bizops360-prod \
  --region=us-central1
```

## Best Practices

1. **Always deploy to dev first** - Test thoroughly before prod
2. **Use separate API keys** - Never share keys between environments
3. **Monitor logs** - Check logs after deployment
4. **Test endpoints** - Run test scripts after deployment
5. **Tag images** - Use version tags for production deployments
6. **Keep secrets secure** - Never commit secrets to git

## Troubleshooting

### Build Fails

```bash
# Check Docker is running
docker ps

# Check GCP authentication
gcloud auth configure-docker
```

### Deployment Fails

```bash
# Check project access
gcloud projects get-iam-policy bizops360-dev
gcloud projects get-iam-policy bizops360-prod

# Check Cloud Run API is enabled
gcloud services list --enabled --project=bizops360-dev
```

### Service Not Responding

```bash
# Check service status
gcloud run services describe bizops360-api-go-dev --project=bizops360-dev

# Check logs
gcloud logging read "resource.type=cloud_run_revision" --project=bizops360-dev --limit=20
```

