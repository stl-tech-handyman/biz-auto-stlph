#!/bin/bash
# Update .env file with Gmail credentials

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

# Find latest credential file
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    TEMP_DIR=$(cmd.exe /c "echo %TEMP%" | tr -d '\r')
    TEMP_DIR=$(cygpath -u "$TEMP_DIR")
else
    TEMP_DIR="/tmp"
fi

CRED_FILE=$(ls -t "$TEMP_DIR"/gmail-credentials-*.json 2>/dev/null | head -1)

if [ -n "$CRED_FILE" ]; then
    echo "Found credential file: $CRED_FILE"
    
    if [ -f "$ENV_FILE" ] && grep -q "GMAIL_CREDENTIALS_JSON=" "$ENV_FILE"; then
        sed -i "s|GMAIL_CREDENTIALS_JSON=.*|GMAIL_CREDENTIALS_JSON=$CRED_FILE|" "$ENV_FILE"
        if ! grep -q "GMAIL_FROM=" "$ENV_FILE"; then
            echo "GMAIL_FROM=team@stlpartyhelpers.com" >> "$ENV_FILE"
        fi
        echo "Updated Gmail credentials in .env file"
    else
        echo "GMAIL_CREDENTIALS_JSON=$CRED_FILE" >> "$ENV_FILE"
        echo "GMAIL_FROM=team@stlpartyhelpers.com" >> "$ENV_FILE"
        echo "Added Gmail credentials to .env file"
    fi
else
    echo "No credential file found. Run: bash scripts/get-gmail-credentials.sh"
    exit 1
fi

