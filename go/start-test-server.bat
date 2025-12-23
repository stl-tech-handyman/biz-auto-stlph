@echo off
echo Starting simple HTTP server for test page...
echo.
echo Server will be available at: http://localhost:8081/test-final-invoice.html
echo.
echo Press Ctrl+C to stop the server
echo.
cd /d %~dp0
python -m http.server 8081
pause

