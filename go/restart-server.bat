@echo off
echo ========================================
echo   Restarting Go API Server
echo ========================================
echo.
echo This will:
echo   1. Stop any running server (if possible)
echo   2. Start the server with new code
echo.
echo Press Ctrl+C to stop the server when done
echo.
pause

cd /d %~dp0
powershell -File start-local.ps1

