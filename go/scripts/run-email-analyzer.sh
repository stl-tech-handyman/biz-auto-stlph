#!/bin/bash

# Run Email Analyzer locally in Cursor
# This uses AI reasoning and can adapt as it processes

echo "========================================"
echo "Email Analyzer - Local Run in Cursor"
echo "========================================"
echo ""

cd ../cmd/email-analyzer

echo "Building..."
go build -o email-analyzer .

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo ""
echo "Running email analyzer..."
echo ""

# Run with default settings (50 emails)
./email-analyzer -max 50 -v

echo ""
echo "Done!"
