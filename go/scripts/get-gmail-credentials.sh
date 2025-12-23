#!/bin/bash
# Get Gmail credentials from GCP Secret Manager for local development

set -e

EMAIL_PROJECT="bizops360-email-dev"
SECRET_NAME="gmail-credentials-json"
GMAIL_FROM_DEFAULT="team@stlpartyhelpers.com"

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

print_info "Fetching Gmail credentials from Secret Manager..."
print_info "Project: $EMAIL_PROJECT"
print_info "Secret: $SECRET_NAME"

# Check if secret exists
if ! gcloud secrets describe "$SECRET_NAME" --project="$EMAIL_PROJECT" >/dev/null 2>&1; then
    print_error "Secret '$SECRET_NAME' not found in project '$EMAIL_PROJECT'"
    print_warning "Make sure the secret exists and you have access to it"
    exit 1
fi

# Get secret value
print_info "Retrieving secret value..."
CREDENTIALS=$(gcloud secrets versions access latest --secret="$SECRET_NAME" --project="$EMAIL_PROJECT" 2>&1)

if [ $? -ne 0 ]; then
    print_error "Failed to retrieve secret value"
    echo "$CREDENTIALS"
    exit 1
fi

# Determine temp directory (works on Windows Git Bash and Unix)
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$(uname -s)" == MINGW* ]] || [[ "$(uname -s)" == MSYS* ]]; then
    # Git Bash on Windows - ALWAYS use Windows temp directory
    # Try multiple common Windows temp locations
    if [ -d "/c/Users/$USER/AppData/Local/Temp" ]; then
        TEMP_DIR="/c/Users/$USER/AppData/Local/Temp"
    elif [ -d "/c/Windows/Temp" ]; then
        TEMP_DIR="/c/Windows/Temp"
    elif [ -n "$TEMP" ] && [[ "$TEMP" != /tmp* ]]; then
        TEMP_DIR="$TEMP"
        # Convert Windows path to Git Bash format if needed
        if [[ "$TEMP_DIR" == C:/* ]] || [[ "$TEMP_DIR" == C:\\* ]]; then
            TEMP_DIR="/c${TEMP_DIR#C:}"
            TEMP_DIR="${TEMP_DIR//\\//}"
        fi
    else
        # Fallback: use Windows temp even if /tmp is set
        TEMP_DIR="/c/Users/$USER/AppData/Local/Temp"
    fi
    TEMP_DIR_FOR_FILE="$TEMP_DIR"
else
    # Unix/Linux/macOS
    if [ -n "$TMPDIR" ]; then
        TEMP_DIR_FOR_FILE="$TMPDIR"
    elif [ -d "/tmp" ]; then
        TEMP_DIR_FOR_FILE="/tmp"
    else
        TEMP_DIR_FOR_FILE="/tmp"
    fi
fi

# Create temp file with timestamp
TIMESTAMP=$(date +"%Y%m%d-%H%M%S")
TEMP_FILE="${TEMP_DIR_FOR_FILE}/gmail-credentials-${TIMESTAMP}.json"

# Save credentials to file
echo "$CREDENTIALS" > "$TEMP_FILE"

# Convert to Windows path format for .env file (if on Windows)
# Go on Windows can handle both /c/... and C:/... formats, but C:/... is more reliable
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
    # Git Bash on Windows - convert /c/Users to C:/Users for .env
    ENV_PATH="$TEMP_FILE"
    if [[ "$ENV_PATH" == /c/* ]]; then
        ENV_PATH="C:${ENV_PATH#/c}"
    fi
    # Ensure forward slashes (Go handles these on Windows)
    ENV_PATH="${ENV_PATH//\\//}"
else
    ENV_PATH="$TEMP_FILE"
fi

print_success "Credentials saved to: $TEMP_FILE"

# Update .env file if we're in the go directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="$GO_DIR/.env"

if [ -f "$ENV_FILE" ]; then
    # Update or add GMAIL_CREDENTIALS_JSON
    if grep -q "GMAIL_CREDENTIALS_JSON=" "$ENV_FILE"; then
        # Use sed with proper escaping for Windows paths
        if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
            # Windows Git Bash - use forward slashes
            sed -i "s|GMAIL_CREDENTIALS_JSON=.*|GMAIL_CREDENTIALS_JSON=$ENV_PATH|" "$ENV_FILE"
        else
            sed -i "s|GMAIL_CREDENTIALS_JSON=.*|GMAIL_CREDENTIALS_JSON=$ENV_PATH|" "$ENV_FILE"
        fi
    else
        echo "GMAIL_CREDENTIALS_JSON=$ENV_PATH" >> "$ENV_FILE"
    fi
    
    # Add or update GMAIL_FROM
    if grep -q "GMAIL_FROM=" "$ENV_FILE"; then
        sed -i "s|GMAIL_FROM=.*|GMAIL_FROM=$GMAIL_FROM_DEFAULT|" "$ENV_FILE"
    else
        echo "GMAIL_FROM=$GMAIL_FROM_DEFAULT" >> "$ENV_FILE"
    fi
    
    print_success "Updated .env file with credentials"
else
    # Create .env file
    echo "GMAIL_CREDENTIALS_JSON=$ENV_PATH" > "$ENV_FILE"
    echo "GMAIL_FROM=$GMAIL_FROM_DEFAULT" >> "$ENV_FILE"
    print_success "Created .env file with credentials"
fi

echo ""
print_info "Credentials are now available in:"
echo "  File: $TEMP_FILE"
echo "  .env: $ENV_FILE"
echo ""
print_info "GMAIL_FROM: $GMAIL_FROM_DEFAULT"
echo ""
print_success "You can now start the server - credentials will be loaded automatically!"

