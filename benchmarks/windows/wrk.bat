@echo off
setlocal enabledelayedexpansion

set "BASE_URL=%BASE_URL%"
if "%BASE_URL%"=="" set "BASE_URL=http://localhost:8080"

set "THREADS=%THREADS%"
if "%THREADS%"=="" set "THREADS=12"

set "CONNECTIONS=%CONNECTIONS%"
if "%CONNECTIONS%"=="" set "CONNECTIONS=400"

set "DURATION=%DURATION%"
if "%DURATION%"=="" set "DURATION=30s"

echo ============================================
echo   wrk - HTTP Load Test (via WSL)
echo ============================================
echo Target:      %BASE_URL%
echo Threads:     %THREADS%
echo Connections: %CONNECTIONS%
echo Duration:    %DURATION%
echo ============================================

echo.
echo --- GET /health ---
wsl -d Ubuntu-24.04 -- wrk -t%THREADS% -c%CONNECTIONS% -d%DURATION% "%BASE_URL%/health"

echo.
echo --- GET /api/v1/users ---
wsl -d Ubuntu-24.04 -- wrk -t%THREADS% -c%CONNECTIONS% -d%DURATION% "%BASE_URL%/api/v1/users"

echo.
echo --- GET /api/v1/users/1 ---
wsl -d Ubuntu-24.04 -- wrk -t%THREADS% -c%CONNECTIONS% -d%DURATION% "%BASE_URL%/api/v1/users/1"

echo.
echo --- POST /api/v1/users ---
wsl -d Ubuntu-24.04 -- bash -c "cat > /tmp/wrk_post.lua << 'EOF'
wrk.method = \"POST\"
wrk.headers[\"Content-Type\"] = \"application/json\"
wrk.body = '{\"name\":\"WrkUser\",\"email\":\"wrk@example.com\",\"password\":\"wrkpass123\"}'
EOF
wrk -t%THREADS% -c%CONNECTIONS% -d%DURATION% -s /tmp/wrk_post.lua %BASE_URL%/api/v1/users"

echo.
echo --- PUT /api/v1/users/1 ---
wsl -d Ubuntu-24.04 -- bash -c "cat > /tmp/wrk_put.lua << 'EOF'
wrk.method = \"PUT\"
wrk.headers[\"Content-Type\"] = \"application/json\"
wrk.body = '{\"name\":\"WrkUpdated\",\"email\":\"wrk.updated@example.com\",\"password\":\"newwrk123\"}'
EOF
wrk -t%THREADS% -c%CONNECTIONS% -d%DURATION% -s /tmp/wrk_put.lua %BASE_URL%/api/v1/users/1"

echo.
echo Done.
endlocal
