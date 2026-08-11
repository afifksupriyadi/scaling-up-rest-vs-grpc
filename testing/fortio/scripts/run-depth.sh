#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/depth"
REST_H2_BASE="http://10.184.0.2:8081/rest"

run_fortio() {
  local format="$1" path="$2" vus="$3" n="$4" outfile="$5"
  fortio load -c "${vus}" -qps 0 -n "${n}" -h2 -httpbufferkb 2048 \
    -json "${outfile}" \
    "${REST_H2_BASE}/${format}/http2/${path}"
}

echo "=== Skenario Kedalaman Data — Level 0 (VU=10, total=500) ==="
run_fortio json     order/depth-zero 10 500 "$RESULTS_DIR/rest-h2-json-depth-zero.json"
run_fortio protobuf order/depth-zero 10 500 "$RESULTS_DIR/rest-h2-protobuf-depth-zero.json"

echo "=== Skenario Kedalaman Data — Level 2 (VU=10, total=500) ==="
run_fortio json     order/depth-two 10 500 "$RESULTS_DIR/rest-h2-json-depth-two.json"
run_fortio protobuf order/depth-two 10 500 "$RESULTS_DIR/rest-h2-protobuf-depth-two.json"

echo "=== Skenario Kedalaman Data — Level 4 (VU=10, total=500) ==="
run_fortio json     order/depth-four 10 500 "$RESULTS_DIR/rest-h2-json-depth-four.json"
run_fortio protobuf order/depth-four 10 500 "$RESULTS_DIR/rest-h2-protobuf-depth-four.json"

echo "Skenario Kedalaman Data (Fortio) selesai — 6 file dihasilkan."