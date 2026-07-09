#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
RATE="${RATE:-50}"
DURATION="${DURATION:-20s}"

echo "============================================"
echo "  vegeta - HTTP Load Test"
echo "============================================"
echo "Target:    $BASE_URL"
echo "Rate:      $RATE req/s"
echo "Duration:  $DURATION"
echo "============================================"

TARGETS_FILE=$(mktemp /tmp/vegeta_targets.XXXXXX)
cat > "$TARGETS_FILE" <<EOF
GET $BASE_URL/health
GET $BASE_URL/api/v1/users
GET $BASE_URL/api/v1/users/1
POST $BASE_URL/api/v1/users
Content-Type: application/json
{"name":"VegetaUser","email":"vegeta@example.com","password":"vegeta123"}
PUT $BASE_URL/api/v1/users/1
Content-Type: application/json
{"name":"VegetaUpdated","email":"vegeta.updated@example.com","password":"newveg123"}
EOF

echo ""
echo "--- All Endpoints ---"
vegeta attack -duration="$DURATION" -rate="$RATE" -targets="$TARGETS_FILE" | vegeta report

echo ""
echo "--- GET /health only ---"
echo "GET $BASE_URL/health" | vegeta attack -duration="$DURATION" -rate="$RATE" | vegeta report

echo ""
echo "--- POST /api/v1/users only ---"
printf "POST $BASE_URL/api/v1/users\nContent-Type: application/json\n{\"name\":\"VegetaUser\",\"email\":\"vegeta@example.com\",\"password\":\"vegeta123\"}\n" | \
  vegeta attack -duration="$DURATION" -rate="$RATE" | vegeta report

rm -f "$TARGETS_FILE"
echo ""
echo "Done."
