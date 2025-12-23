@echo off
echo Opening test page...
start http://localhost:8080/test-final-invoice.html
timeout /t 2 /nobreak >nul
echo.
echo If page doesn't load correctly, the server needs to be restarted.
echo Please restart the server with: .\start-local.ps1
pause

