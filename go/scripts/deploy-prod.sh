#!/bin/bash
# Deploy Go API to Production Environment

set -e

# Configuration
PROJECT_ID="bizops360-prod"
REGION="us-central1"
SERVICE_NAME="bizops360-api-go-prod"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

# Check prerequisites
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed"
    exit 1
fi

# Confirmation prompt for production
print_warning "⚠️  You are about to deploy to PRODUCTION!"
read -p "Are you sure you want to continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    print_status "Deployment cancelled"
    exit 0
fi

print_status "Deploying Go API to PROD environment..."
print_status "Project: $PROJECT_ID"
print_status "Service: $SERVICE_NAME"
print_status "Region: $REGION"

# Set project
gcloud config set project "$PROJECT_ID"

# Build Docker image
print_status "Building Docker image..."
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"
docker build -f go/Dockerfile.prod -t "$IMAGE_NAME:latest" .

# Push to Artifact Registry
print_status "Pushing image to Artifact Registry..."
docker push "$IMAGE_NAME:latest"

# Prepare secrets
SECRET_ARGS="SERVICE_API_KEY=svc-api-key-prod:latest"

if gcloud secrets describe stripe-secret-key-prod --project="$PROJECT_ID" >/dev/null 2>&1; then
    SECRET_ARGS="$SECRET_ARGS,STRIPE_SECRET_KEY_PROD=stripe-secret-key-prod:latest"
else
    print_error "Secret stripe-secret-key-prod not found in $PROJECT_ID"
    exit 1
fi

if gcloud secrets describe stripe-secret-key-test --project="$PROJECT_ID" >/dev/null 2>&1; then
    SECRET_ARGS="$SECRET_ARGS,STRIPE_SECRET_KEY_TEST=stripe-secret-key-test:latest"
fi

# Deploy to Cloud Run
print_status "Deploying to Cloud Run..."
gcloud run deploy "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --image="$IMAGE_NAME:latest" \
    --platform=managed \
    --allow-unauthenticated \
    --port=8080 \
    --memory=1Gi \
    --cpu=2 \
    --min-instances=1 \
    --max-instances=20 \
    --concurrency=100 \
    --timeout=300 \
    --set-secrets="$SECRET_ARGS" \
    --set-env-vars="ENV=prod,LOG_LEVEL=info,CONFIG_DIR=/app/config,TEMPLATES_DIR=/app/templates" \
    --labels="env=prod,service=api-go,type=cloud-run"

# Get service URL
SERVICE_URL=$(gcloud run services describe "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --format="value(status.url)")

print_success "Deployment complete!"
print_status "Service URL: $SERVICE_URL"
print_status "Health check: $SERVICE_URL/api/health"
print_status "Stripe endpoint: $SERVICE_URL/api/stripe/deposit/calculate"

