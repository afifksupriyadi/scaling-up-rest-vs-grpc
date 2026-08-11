#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/concurrency"

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

echo "=== Skenario Tingkat Konkurensi — 1 VU (dataset: Hundred, total=500) ==="
run_rest json     http1 order/hundred 1 500 "$RESULTS_DIR/rest-h1-json-vu1.csv"
run_rest protobuf http1 order/hundred 1 500 "$RESULTS_DIR/rest-h1-protobuf-vu1.csv"
run_rest json     http2 order/hundred 1 500 "$RESULTS_DIR/rest-h2-json-vu1.csv"
run_rest protobuf http2 order/hundred 1 500 "$RESULTS_DIR/rest-h2-protobuf-vu1.csv"
run_grpc  GetOrderHundred                1 500 "$RESULTS_DIR/grpc-vu1.csv"

echo "=== Skenario Tingkat Konkurensi — 10 VU (dataset: Hundred, total=500) ==="
run_rest json     http1 order/hundred 10 50 "$RESULTS_DIR/rest-h1-json-vu10.csv"
run_rest protobuf http1 order/hundred 10 50 "$RESULTS_DIR/rest-h1-protobuf-vu10.csv"
run_rest json     http2 order/hundred 10 50 "$RESULTS_DIR/rest-h2-json-vu10.csv"
run_rest protobuf http2 order/hundred 10 50 "$RESULTS_DIR/rest-h2-protobuf-vu10.csv"
run_grpc  GetOrderHundred                10 50 "$RESULTS_DIR/grpc-vu10.csv"

echo "=== Skenario Tingkat Konkurensi — 100 VU (dataset: Hundred, total=500) ==="
run_rest json     http1 order/hundred 100 5 "$RESULTS_DIR/rest-h1-json-vu100.csv"
run_rest protobuf http1 order/hundred 100 5 "$RESULTS_DIR/rest-h1-protobuf-vu100.csv"
run_rest json     http2 order/hundred 100 5 "$RESULTS_DIR/rest-h2-json-vu100.csv"
run_rest protobuf http2 order/hundred 100 5 "$RESULTS_DIR/rest-h2-protobuf-vu100.csv"
run_grpc  GetOrderHundred                100 5 "$RESULTS_DIR/grpc-vu100.csv"

echo "Skenario Tingkat Konkurensi selesai — 15 file dihasilkan."