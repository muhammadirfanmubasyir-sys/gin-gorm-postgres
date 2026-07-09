#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "============================================"
echo "  k6 - HTTP Load Test"
echo "============================================"
echo "Target: $BASE_URL"
echo "============================================"

k6 run \
  -e "BASE_URL=$BASE_URL" \
  --out json="$SCRIPT_DIR/k6-results.json" \
  "$SCRIPT_DIR/k6.js"

echo ""
echo "Results saved to: $SCRIPT_DIR/k6-results.json"
echo "Done."
