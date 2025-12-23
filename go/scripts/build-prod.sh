#!/bin/bash
# Build Go API for production environment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

PROJECT_ID="bizops360-prod"
SERVICE_NAME="bizops360-api-go-prod"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

echo "Building Go API for PROD environment..."
echo "Image: $IMAGE_NAME"

cd "$PROJECT_ROOT"

# Build using prod Dockerfile (from project root)
docker build -f go/Dockerfile.prod -t "$IMAGE_NAME:latest" .
docker tag "$IMAGE_NAME:latest" "$IMAGE_NAME:$(date +%Y%m%d-%H%M%S)"

echo "Build complete!"
echo "To push: docker push $IMAGE_NAME:latest"

