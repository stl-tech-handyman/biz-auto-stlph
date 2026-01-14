# Quick script to check if firewall rules exist
# Run: powershell.exe -ExecutionPolicy Bypass -File scripts/check-firewall-rules.ps1

Write-Host "Checking for STL Party Helpers API firewall rules..." -ForegroundColor Cyan
Write-Host ""

$rules = Get-NetFirewallRule -DisplayName "*STL Party Helpers*" -ErrorAction SilentlyContinue

if ($rules) {
    Write-Host "Found firewall rules:" -ForegroundColor Green
    $rules | Format-Table DisplayName, Enabled, Direction, Action, Profile -AutoSize
    
    $inbound = $rules | Where-Object { $_.Direction -eq "Inbound" -and $_.Enabled -eq $true }
    if ($inbound) {
        Write-Host "✓ Inbound rule is enabled" -ForegroundColor Green
    } else {
        Write-Host "✗ Inbound rule is missing or disabled" -ForegroundColor Red
    }
} else {
    Write-Host "✗ No firewall rules found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Run the add-firewall-rule.ps1 script as Administrator to create the rules." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "To add/update rules, run as Administrator:" -ForegroundColor Cyan
Write-Host '  .\scripts\add-firewall-rule.ps1' -ForegroundColor White

