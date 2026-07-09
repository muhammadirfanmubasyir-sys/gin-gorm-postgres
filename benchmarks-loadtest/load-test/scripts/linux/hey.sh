#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
REQUESTS="${REQUESTS:-1000}"
CONCURRENCY="${CONCURRENCY:-50}"

echo "============================================"
echo "  hey - HTTP Load Test"
echo "============================================"
echo "Target:       $BASE_URL"
echo "Requests:     $REQUESTS"
echo "Concurrency:  $CONCURRENCY"
echo "============================================"

echo ""
echo "--- GET /health ---"
hey -n "$REQUESTS" -c "$CONCURRENCY" "$BASE_URL/health"

echo ""
echo "--- POST /api/v1/users ---"
hey -n "$REQUESTS" -c "$CONCURRENCY" \
  -m POST \
  -H "Content-Type: application/json" \
  -d '{"name":"LoadUser","email":"load@example.com","password":"loadtest123"}' \
  "$BASE_URL/api/v1/users"

echo ""
echo "--- GET /api/v1/users ---"
hey -n "$REQUESTS" -c "$CONCURRENCY" "$BASE_URL/api/v1/users"

echo ""
echo "--- GET /api/v1/users/1 ---"
hey -n "$REQUESTS" -c "$CONCURRENCY" "$BASE_URL/api/v1/users/1"

echo ""
echo "--- PUT /api/v1/users/1 ---"
hey -n "$REQUESTS" -c "$CONCURRENCY" \
  -m PUT \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated","email":"updated@example.com","password":"newpass123"}' \
  "$BASE_URL/api/v1/users/1"

echo ""
echo "Done."
