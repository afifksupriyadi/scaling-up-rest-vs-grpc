#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/element-count"
REST_H2_BASE="http://10.184.0.2:8081/rest"

run_fortio() {
  local format="$1" path="$2" vus="$3" n="$4" outfile="$5"
  fortio load -c "${vus}" -qps 0 -n "${n}" -h2 -httpbufferkb 2048 \
    -json "${outfile}" \
    "${REST_H2_BASE}/${format}/http2/${path}"
}

echo "=== Skenario Jumlah Elemen — 1 Elemen (VU=10, total=500) ==="
run_fortio json     order/one 10 500 "$RESULTS_DIR/rest-h2-json-one.json"
run_fortio protobuf order/one 10 500 "$RESULTS_DIR/rest-h2-protobuf-one.json"

echo "=== Skenario Jumlah Elemen — 100 Elemen (VU=10, total=500) ==="
run_fortio json     order/hundred 10 500 "$RESULTS_DIR/rest-h2-json-hundred.json"
run_fortio protobuf order/hundred 10 500 "$RESULTS_DIR/rest-h2-protobuf-hundred.json"

echo "=== Skenario Jumlah Elemen — 1000 Elemen (VU=10, total=500) ==="
run_fortio json     order/thousand 10 500 "$RESULTS_DIR/rest-h2-json-thousand.json"
run_fortio protobuf order/thousand 10 500 "$RESULTS_DIR/rest-h2-protobuf-thousand.json"

echo "Skenario Jumlah Elemen (Fortio) selesai — 6 file dihasilkan."