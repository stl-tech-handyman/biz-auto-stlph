# Dev and Prod Environments Setup

## Quick Start

### Deploy to Development

**Bash (Linux/Mac/Git Bash):**
```bash
cd go
./scripts/deploy-dev.sh
```

**PowerShell (Windows):**
```powershell
cd go
.\scripts\deploy-dev.ps1
```

**Make (Cross-platform):**
```bash
cd go
make deploy-dev
```

### Deploy to Production

**Bash:**
```bash
cd go
./scripts/deploy-prod.sh
```

**PowerShell:**
```powershell
cd go
.\scripts\deploy-prod.ps1
```

**Make:**
```bash
cd go
make deploy-prod
```

### Deploy to Both

```bash
cd go
./scripts/deploy-all.sh
# or
make deploy-all
```

## Environment Configuration

### Development (`bizops360-dev`)

- **Service**: `bizops360-api-go-dev`
- **Resources**: 512Mi RAM, 1 CPU, 0-5 instances
- **Log Level**: `debug`
- **Stripe**: Uses test keys
- **Dockerfile**: `Dockerfile.dev` (includes debugging tools)

### Production (`bizops360-prod`)

- **Service**: `bizops360-api-go-prod`
- **Resources**: 1Gi RAM, 2 CPU, 1-20 instances
- **Log Level**: `info`
- **Stripe**: Uses live keys
- **Dockerfile**: `Dockerfile.prod` (optimized, distroless)

## Testing Deployments

```bash
# Test dev
./scripts/test-dev.sh
# or
make test-dev

# Test prod
./scripts/test-prod.sh
# or
make test-prod
```

## Environment Variables

Both environments use the same environment variables, but with different values:

| Variable | Dev | Prod |
|----------|-----|------|
| `ENV` | `dev` | `prod` |
| `LOG_LEVEL` | `debug` | `info` |
| `SERVICE_API_KEY` | From `svc-api-key-dev` secret | From `svc-api-key-prod` secret |
| `STRIPE_SECRET_KEY_TEST` | From Secret Manager | From Secret Manager |
| `STRIPE_SECRET_KEY_PROD` | Optional | Required |

## Service URLs

After deployment, get URLs:

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

## See Also

- `go/ENVIRONMENTS.md` - Detailed environment guide
- `go/scripts/` - Deployment scripts
- `go/config/environments/` - Environment configuration files

