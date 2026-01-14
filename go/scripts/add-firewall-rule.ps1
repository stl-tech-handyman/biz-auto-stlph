# PowerShell script to add Windows Firewall rule for API port
# Run this script as Administrator: Right-click PowerShell -> Run as administrator
# This creates a port-based rule, so it works regardless of where the executable is built

$ErrorActionPreference = 'Stop'

# Default API port (can be overridden)
$ApiPort = 8080

Write-Host "Adding Windows Firewall rule for STL Party Helpers API on port $ApiPort" -ForegroundColor Cyan
Write-Host ""
Write-Host "Note: This creates a port-based rule, so it works regardless of where api.exe is built." -ForegroundColor Yellow
Write-Host ""

try {
    # Remove existing rules if they exist
    Remove-NetFirewallRule -DisplayName "STL Party Helpers API (Port $ApiPort)" -ErrorAction SilentlyContinue
    Remove-NetFirewallRule -DisplayName "STL Party Helpers API (Port $ApiPort - Outbound)" -ErrorAction SilentlyContinue
    
    # Add new inbound rule for the port
    $inboundDisplayName = "STL Party Helpers API (Port $ApiPort)"
    $inboundDescription = "Allow STL Party Helpers API to accept incoming connections on port $ApiPort"
    
    New-NetFirewallRule -DisplayName $inboundDisplayName -Direction Inbound -LocalPort $ApiPort -Protocol TCP -Action Allow -Profile Domain,Private,Public -EdgeTraversalPolicy Allow -Description $inboundDescription
    
    Write-Host "✓ Inbound firewall rule added successfully for port $ApiPort" -ForegroundColor Green
    
    # Add outbound rule (usually not needed, but some apps require it)
    $outboundDisplayName = "STL Party Helpers API (Port $ApiPort - Outbound)"
    $outboundDescription = "Allow STL Party Helpers API to make outgoing connections on port $ApiPort"
    
    New-NetFirewallRule -DisplayName $outboundDisplayName -Direction Outbound -LocalPort $ApiPort -Protocol TCP -Action Allow -Profile Domain,Private,Public -Description $outboundDescription
    
    Write-Host "✓ Outbound firewall rule added successfully for port $ApiPort" -ForegroundColor Green
    Write-Host ""
    Write-Host "Firewall rules configured successfully!" -ForegroundColor Green
    Write-Host "The API should no longer be blocked by Windows Firewall on port $ApiPort." -ForegroundColor Green
    Write-Host ""
    Write-Host "You can verify the rules in:" -ForegroundColor Cyan
    Write-Host "  Windows Defender Firewall -> Advanced settings -> Inbound Rules" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "To use a different port, edit this script and change `$ApiPort variable." -ForegroundColor Gray
    
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Make sure you're running PowerShell as Administrator!" -ForegroundColor Yellow
    exit 1
}
