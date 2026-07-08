@echo off
setlocal enabledelayedexpansion

set "BASE_URL=%BASE_URL%"
if "%BASE_URL%"=="" set "BASE_URL=http://localhost:8080"

echo ============================================
echo   k6 - HTTP Load Test
echo ============================================
echo Target: %BASE_URL%
echo ============================================

set "SCRIPT_DIR=%~dp0"

k6 run -e "BASE_URL=%BASE_URL%" --out json="%SCRIPT_DIR%k6-results.json" "%SCRIPT_DIR%k6.js"

echo.
echo Results saved to: %SCRIPT_DIR%k6-results.json
echo Done.
endlocal
