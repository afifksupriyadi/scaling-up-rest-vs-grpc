#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/concurrency"
REST_H2_BASE="http://10.184.0.2:8081/rest"

run_fortio() {
  local format="$1" path="$2" vus="$3" n="$4" outfile="$5"
  fortio load -c "${vus}" -qps 0 -n "${n}" -h2 -httpbufferkb 2048 \
    -json "${outfile}" \
    "${REST_H2_BASE}/${format}/http2/${path}"
}

echo "=== Skenario Tingkat Konkurensi — 1 VU (dataset: Hundred, total=500) ==="
run_fortio json     order/hundred 1 500 "$RESULTS_DIR/rest-h2-json-vu1.json"
run_fortio protobuf order/hundred 1 500 "$RESULTS_DIR/rest-h2-protobuf-vu1.json"

echo "=== Skenario Tingkat Konkurensi — 10 VU (dataset: Hundred, total=500) ==="
run_fortio json     order/hundred 10 500 "$RESULTS_DIR/rest-h2-json-vu10.json"
run_fortio protobuf order/hundred 10 500 "$RESULTS_DIR/rest-h2-protobuf-vu10.json"

echo "=== Skenario Tingkat Konkurensi — 100 VU (dataset: Hundred, total=500) ==="
run_fortio json     order/hundred 100 500 "$RESULTS_DIR/rest-h2-json-vu100.json"
run_fortio protobuf order/hundred 100 500 "$RESULTS_DIR/rest-h2-protobuf-vu100.json"

echo "Skenario Tingkat Konkurensi (Fortio) selesai — 6 file dihasilkan."