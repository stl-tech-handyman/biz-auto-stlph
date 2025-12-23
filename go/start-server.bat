@echo off
cd /d %~dp0
start "Go API Server" powershell -NoExit -File start-local.ps1
timeout /t 3 /nobreak >nul
echo Server should be starting in a new window...
echo Check the new PowerShell window for status.
pause

