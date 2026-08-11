#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/element-count"

REST_TEMPLATE="$SCRIPT_DIR/rest-template.js"
GRPC_TEMPLATE="$SCRIPT_DIR/grpc-template.js"

run_rest() {
  local format="$1" protocol="$2" path="$3" vus="$4" iterations="$5" outfile="$6"
  local port=8080
  [ "$protocol" = "http2" ] && port=8081
  k6 run \
    -e TARGET_URL="https://10.184.0.2:${port}/rest/${format}/${protocol}/${path}" \
    -e VUS="${vus}" \
    -e ITERATIONS="${iterations}" \
    --out csv="${outfile}" \
    "${REST_TEMPLATE}"
}

run_grpc() {
  local method="$1" vus="$2" iterations="$3" outfile="$4"
  k6 run \
    -e GRPC_METHOD="${method}" \
    -e VUS="${vus}" \
    -e ITERATIONS="${iterations}" \
    --out csv="${outfile}" \
    "${GRPC_TEMPLATE}"
}

echo "=== Skenario Jumlah Elemen — 1 Elemen (VU=10, iterasi=50) ==="
run_rest json     http1 order/one 10 50 "$RESULTS_DIR/rest-h1-json-one.csv"
run_rest protobuf http1 order/one 10 50 "$RESULTS_DIR/rest-h1-protobuf-one.csv"
run_rest json     http2 order/one 10 50 "$RESULTS_DIR/rest-h2-json-one.csv"
run_rest protobuf http2 order/one 10 50 "$RESULTS_DIR/rest-h2-protobuf-one.csv"
run_grpc  GetOrderOne                10 50 "$RESULTS_DIR/grpc-one.csv"

echo "=== Skenario Jumlah Elemen — 100 Elemen (VU=10, iterasi=50) ==="
run_rest json     http1 order/hundred 10 50 "$RESULTS_DIR/rest-h1-json-hundred.csv"
run_rest protobuf http1 order/hundred 10 50 "$RESULTS_DIR/rest-h1-protobuf-hundred.csv"
run_rest json     http2 order/hundred 10 50 "$RESULTS_DIR/rest-h2-json-hundred.csv"
run_rest protobuf http2 order/hundred 10 50 "$RESULTS_DIR/rest-h2-protobuf-hundred.csv"
run_grpc  GetOrderHundred                10 50 "$RESULTS_DIR/grpc-hundred.csv"

echo "=== Skenario Jumlah Elemen — 1000 Elemen (VU=10, iterasi=50) ==="
run_rest json     http1 order/thousand 10 50 "$RESULTS_DIR/rest-h1-json-thousand.csv"
run_rest protobuf http1 order/thousand 10 50 "$RESULTS_DIR/rest-h1-protobuf-thousand.csv"
run_rest json     http2 order/thousand 10 50 "$RESULTS_DIR/rest-h2-json-thousand.csv"
run_rest protobuf http2 order/thousand 10 50 "$RESULTS_DIR/rest-h2-protobuf-thousand.csv"
run_grpc  GetOrderThousand                10 50 "$RESULTS_DIR/grpc-thousand.csv"

echo "Skenario Jumlah Elemen selesai — 15 file dihasilkan."