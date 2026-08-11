#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/depth"

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

echo "=== Skenario Kedalaman Data — Level 0 (VU=10, iterasi=50) ==="
run_rest json     http1 order/depth-zero 10 50 "$RESULTS_DIR/rest-h1-json-depth-zero.csv"
run_rest protobuf http1 order/depth-zero 10 50 "$RESULTS_DIR/rest-h1-protobuf-depth-zero.csv"
run_rest json     http2 order/depth-zero 10 50 "$RESULTS_DIR/rest-h2-json-depth-zero.csv"
run_rest protobuf http2 order/depth-zero 10 50 "$RESULTS_DIR/rest-h2-protobuf-depth-zero.csv"
run_grpc  GetOrderDepthZero                10 50 "$RESULTS_DIR/grpc-depth-zero.csv"

echo "=== Skenario Kedalaman Data — Level 2 (VU=10, iterasi=50) ==="
run_rest json     http1 order/depth-two 10 50 "$RESULTS_DIR/rest-h1-json-depth-two.csv"
run_rest protobuf http1 order/depth-two 10 50 "$RESULTS_DIR/rest-h1-protobuf-depth-two.csv"
run_rest json     http2 order/depth-two 10 50 "$RESULTS_DIR/rest-h2-json-depth-two.csv"
run_rest protobuf http2 order/depth-two 10 50 "$RESULTS_DIR/rest-h2-protobuf-depth-two.csv"
run_grpc  GetOrderDepthTwo                10 50 "$RESULTS_DIR/grpc-depth-two.csv"

echo "=== Skenario Kedalaman Data — Level 4 (VU=10, iterasi=50) ==="
run_rest json     http1 order/depth-four 10 50 "$RESULTS_DIR/rest-h1-json-depth-four.csv"
run_rest protobuf http1 order/depth-four 10 50 "$RESULTS_DIR/rest-h1-protobuf-depth-four.csv"
run_rest json     http2 order/depth-four 10 50 "$RESULTS_DIR/rest-h2-json-depth-four.csv"
run_rest protobuf http2 order/depth-four 10 50 "$RESULTS_DIR/rest-h2-protobuf-depth-four.csv"
run_grpc  GetOrderDepthFour                10 50 "$RESULTS_DIR/grpc-depth-four.csv"

echo "Skenario Kedalaman Data selesai — 15 file dihasilkan."