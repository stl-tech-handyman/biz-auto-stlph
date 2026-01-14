#!/bin/bash

# Script to add Windows Firewall rule for API port
# This allows the API to accept incoming connections without Windows Firewall blocking it
# Uses port-based rule, so it works regardless of where the executable is built

# Default API port (can be changed if needed)
API_PORT=8080

echo "Adding Windows Firewall rule for STL Party Helpers API on port $API_PORT"
echo ""
echo "Note: This creates a port-based rule, so it works regardless of where api.exe is built."
echo ""

# Use PowerShell to add firewall rule
# Run the PowerShell script directly (it will use default port 8080)
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "& { Set-Location '$PWD'; & './scripts/add-firewall-rule.ps1' }"

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Firewall rules added successfully for port $API_PORT!"
    echo "  You can verify the rules in Windows Defender Firewall settings"
    echo "  The API should no longer be blocked by Windows Firewall on port $API_PORT."
    echo ""
    echo "To use a different port, edit this script and change API_PORT variable."
else
    echo ""
    echo "✗ Failed to add firewall rules"
    echo ""
    echo "This script requires Administrator privileges."
    echo ""
    echo "To run as Administrator:"
    echo "  1. Right-click on Git Bash"
    echo "  2. Select 'Run as administrator'"
    echo "  3. Navigate to the go directory: cd /c/Users/Alexey/Code/biz-operating-system/stlph/go"
    echo "  4. Run: bash scripts/add-firewall-rule.sh"
    echo ""
    echo "Or run PowerShell as Administrator and execute:"
    echo "  New-NetFirewallRule -DisplayName 'STL Party Helpers API (Port 8080)' -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow -Profile Domain,Private,Public"
    exit 1
fi

