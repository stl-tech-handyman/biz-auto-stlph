# Go API Dev/Prod Setup Summary

## ✅ Completed

Successfully created separate **dev** and **prod** environments for the Go API, matching the JavaScript API structure.

## Environment Structure

### Development (`bizops360-dev`)
- **Service Name**: `bizops360-api-go-dev`
- **GCP Project**: `bizops360-dev`
- **Docker Image**: `gcr.io/bizops360-dev/bizops360-api-go-dev`
- **Resources**: 512Mi RAM, 1 CPU, 0-5 instances
- **Log Level**: `debug`
- **Features**: Debug tools, detailed logging, Stripe test keys

### Production (`bizops360-prod`)
- **Service Name**: `bizops360-api-go-prod`
- **GCP Project**: `bizops360-prod`
- **Docker Image**: `gcr.io/bizops360-prod/bizops360-api-go-prod`
- **Resources**: 1Gi RAM, 2 CPU, 1-20 instances
- **Log Level**: `info`
- **Features**: Optimized build, minimal logging, Stripe live keys

## Files Created

### Configuration Files
- `go/config/environments/dev.yaml` - Dev environment config
- `go/config/environments/prod.yaml` - Prod environment config
- `go/internal/config/environment.go` - Environment detection logic

### Dockerfiles
- `go/Dockerfile.dev` - Development build (with debug tools)
- `go/Dockerfile.prod` - Production build (optimized, distroless)
- `go/Dockerfile` - Default (production)

### Deployment Scripts

**Bash (Linux/Mac/Git Bash):**
- `go/scripts/deploy-dev.sh` - Deploy to dev
- `go/scripts/deploy-prod.sh` - Deploy to prod (with confirmation)
- `go/scripts/deploy-all.sh` - Deploy to both
- `go/scripts/build-dev.sh` - Build dev image
- `go/scripts/build-prod.sh` - Build prod image
- `go/scripts/test-dev.sh` - Test dev deployment
- `go/scripts/test-prod.sh` - Test prod deployment

**PowerShell (Windows):**
- `go/scripts/deploy-dev.ps1` - Deploy to dev
- `go/scripts/deploy-prod.ps1` - Deploy to prod (with confirmation)

**Makefile:**
- `go/Makefile` - Cross-platform commands

### Documentation
- `go/ENVIRONMENTS.md` - Detailed environment guide
- `go/README_ENVIRONMENTS.md` - Quick start guide

## Usage

### Quick Deploy

**Development:**
```bash
cd go
./scripts/deploy-dev.sh
# or
make deploy-dev
# or (Windows)
.\scripts\deploy-dev.ps1
```

**Production:**
```bash
cd go
./scripts/deploy-prod.sh
# or
make deploy-prod
# or (Windows)
.\scripts\deploy-prod.ps1
```

### Testing

```bash
# Test dev
./scripts/test-dev.sh
make test-dev

# Test prod
./scripts/test-prod.sh
make test-prod
```

## Environment Detection

The Go API automatically detects the environment from the `ENV` environment variable:

- `ENV=dev` → Development mode
- `ENV=prod` → Production mode
- Default → Development (safe default)

The code uses:
- `config.GetEnvironment()` - Returns current environment
- `config.IsProduction()` - Returns true if prod
- `config.IsDevelopment()` - Returns true if dev
- `config.GetServiceName()` - Returns appropriate service name

## Secret Management

### Required Secrets

**Development (`bizops360-dev`):**
- `svc-api-key-dev` - API authentication key
- `stripe-secret-key-test` - Stripe test API key

**Production (`bizops360-prod`):**
- `svc-api-key-prod` - API authentication key
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

## Deployment Flow

1. **Build** - Creates Docker image with appropriate Dockerfile
2. **Push** - Pushes to Artifact Registry in respective project
3. **Deploy** - Deploys to Cloud Run with environment-specific settings
4. **Configure** - Sets secrets and environment variables
5. **Verify** - Tests endpoints to confirm deployment

## Differences from JS API

### Same:
- ✅ Separate GCP projects (`bizops360-dev` / `bizops360-prod`)
- ✅ Separate Cloud Run services
- ✅ Separate Secret Manager secrets
- ✅ Environment variable `ENV` for detection
- ✅ Different resource allocations

### Enhanced:
- ✅ Separate Dockerfiles (dev with debug tools, prod optimized)
- ✅ Environment detection helpers (`IsProduction()`, `IsDevelopment()`)
- ✅ Cross-platform deployment scripts (Bash + PowerShell)
- ✅ Makefile for easy commands
- ✅ Comprehensive test scripts

## Testing

All environment detection tests pass:
```bash
$ go test ./internal/config/... -v
=== RUN   TestGetEnvironment
--- PASS: TestGetEnvironment (0.00s)
=== RUN   TestIsProduction
--- PASS: TestIsProduction (0.00s)
=== RUN   TestIsDevelopment
--- PASS: TestIsDevelopment (0.00s)
=== RUN   TestGetServiceName
--- PASS: TestGetServiceName (0.00s)
PASS
```

## Next Steps

1. **Create Secrets**: Set up Secret Manager secrets in both projects
2. **Enable APIs**: Ensure Artifact Registry and Cloud Run APIs are enabled
3. **Deploy Dev**: Test deployment to dev environment first
4. **Verify**: Run test scripts to confirm endpoints work
5. **Deploy Prod**: Deploy to production after dev verification

## See Also

- `go/ENVIRONMENTS.md` - Complete environment guide
- `go/README_ENVIRONMENTS.md` - Quick reference
- `go/scripts/` - All deployment scripts
- `go/Makefile` - Cross-platform commands

