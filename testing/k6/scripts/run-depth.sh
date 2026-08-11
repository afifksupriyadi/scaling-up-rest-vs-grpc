#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/depth"

REST_TEMPLATE="$SCRIPT_DIR/rest-template.js"
GRPC_TEMPLATE="$SCRIPT_DIR/grpc-template.js"
REST_H1_BASE="http://10.184.0.2:8080/rest"

run_rest() {
  local format="$1" path="$2" vus="$3" iterations="$4" outfile="$5"
  k6 run \
    -e TARGET_URL="${REST_H1_BASE}/${format}/http1/${path}" \
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
run_rest json     order/depth-zero 10 50 "$RESULTS_DIR/rest-h1-json-depth-zero.csv"
run_rest protobuf order/depth-zero 10 50 "$RESULTS_DIR/rest-h1-protobuf-depth-zero.csv"
run_grpc  GetOrderDepthZero        10 50 "$RESULTS_DIR/grpc-depth-zero.csv"

echo "=== Skenario Kedalaman Data — Level 2 (VU=10, iterasi=50) ==="
run_rest json     order/depth-two 10 50 "$RESULTS_DIR/rest-h1-json-depth-two.csv"
run_rest protobuf order/depth-two 10 50 "$RESULTS_DIR/rest-h1-protobuf-depth-two.csv"
run_grpc  GetOrderDepthTwo        10 50 "$RESULTS_DIR/grpc-depth-two.csv"

echo "=== Skenario Kedalaman Data — Level 4 (VU=10, iterasi=50) ==="
run_rest json     order/depth-four 10 50 "$RESULTS_DIR/rest-h1-json-depth-four.csv"
run_rest protobuf order/depth-four 10 50 "$RESULTS_DIR/rest-h1-protobuf-depth-four.csv"
run_grpc  GetOrderDepthFour        10 50 "$RESULTS_DIR/grpc-depth-four.csv"

echo "Skenario Kedalaman Data (k6) selesai — 9 file dihasilkan."