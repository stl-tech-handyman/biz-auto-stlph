@echo off
REM Run Email Analyzer locally in Cursor
REM This uses AI reasoning and can adapt as it processes

echo ========================================
echo Email Analyzer - Local Run in Cursor
echo ========================================
echo.

cd ..\cmd\email-analyzer

echo Building...
go build -o email-analyzer.exe .

if errorlevel 1 (
    echo Build failed!
    pause
    exit /b 1
)

echo.
echo Running email analyzer...
echo.

REM Run with default settings (50 emails)
email-analyzer.exe -max 50 -v

echo.
echo Done!
pause
