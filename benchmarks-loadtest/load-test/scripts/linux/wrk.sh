#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
THREADS="${THREADS:-12}"
CONNECTIONS="${CONNECTIONS:-400}"
DURATION="${DURATION:-30s}"

echo "============================================"
echo "  wrk - HTTP Load Test (via WSL)"
echo "============================================"
echo "Target:      $BASE_URL"
echo "Threads:     $THREADS"
echo "Connections: $CONNECTIONS"
echo "Duration:    $DURATION"
echo "============================================"

wrk_cmd() {
  wsl -d Ubuntu-24.04 -- wrk "$@"
}

tmpfile=$(mktemp /tmp/wrk_post.lua.XXXXXX)
cat > "$tmpfile" <<'LUA'
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{"name":"WrkUser","email":"wrk@example.com","password":"wrkpass123"}'
LUA

tmpfile_put=$(mktemp /tmp/wrk_put.lua.XXXXXX)
cat > "$tmpfile_put" <<'LUA'
wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{"name":"WrkUpdated","email":"wrk.updated@example.com","password":"newwrk123"}'
LUA

echo ""
echo "--- GET /health ---"
wrk_cmd -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" "$BASE_URL/health"

echo ""
echo "--- GET /api/v1/users ---"
wrk_cmd -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" "$BASE_URL/api/v1/users"

echo ""
echo "--- GET /api/v1/users/1 ---"
wrk_cmd -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" "$BASE_URL/api/v1/users/1"

echo ""
echo "--- POST /api/v1/users ---"
wsl -d Ubuntu-24.04 -- bash -c "wrk -t$THREADS -c$CONNECTIONS -d$DURATION -s /tmp/wrk_post.lua $BASE_URL/api/v1/users"

echo ""
echo "--- PUT /api/v1/users/1 ---"
wsl -d Ubuntu-24.04 -- bash -c "wrk -t$THREADS -c$CONNECTIONS -d$DURATION -s /tmp/wrk_put.lua $BASE_URL/api/v1/users/1"

rm -f "$tmpfile" "$tmpfile_put"
echo ""
echo "Done."
