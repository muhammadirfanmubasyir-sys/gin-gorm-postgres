@echo off
setlocal enabledelayedexpansion

set "BASE_URL=%BASE_URL%"
if "%BASE_URL%"=="" set "BASE_URL=http://localhost:8080"

set "RATE=%RATE%"
if "%RATE%"=="" set "RATE=50"

set "DURATION=%DURATION%"
if "%DURATION%"=="" set "DURATION=20s"

echo ============================================
echo   vegeta - HTTP Load Test
echo ============================================
echo Target:    %BASE_URL%
echo Rate:      %RATE% req/s
echo Duration:  %DURATION%
echo ============================================

echo.
echo --- All Endpoints ---
(
echo GET %BASE_URL%/health
echo GET %BASE_URL%/api/v1/users
echo GET %BASE_URL%/api/v1/users/1
echo POST %BASE_URL%/api/v1/users
echo Content-Type: application/json
echo {"name":"VegetaUser","email":"vegeta@example.com","password":"vegeta123"}
echo PUT %BASE_URL%/api/v1/users/1
echo Content-Type: application/json
echo {"name":"VegetaUpdated","email":"vegeta.updated@example.com","password":"newveg123"}
) | vegeta attack -duration=%DURATION% -rate=%RATE% | vegeta report

echo.
echo --- GET /health only ---
echo GET %BASE_URL%/health | vegeta attack -duration=%DURATION% -rate=%RATE% | vegeta report

echo.
echo --- POST /api/v1/users only ---
(
echo POST %BASE_URL%/api/v1/users
echo Content-Type: application/json
echo {"name":"VegetaUser","email":"vegeta@example.com","password":"vegeta123"}
) | vegeta attack -duration=%DURATION% -rate=%RATE% | vegeta report

echo.
echo Done.
endlocal
