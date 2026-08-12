#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/multiplexing"
PROTO_PATH="$SCRIPT_DIR/../../proto/order.proto"

REST_H1_URL_JSON="https://10.184.0.2:8080/rest/json/http1/order/hundred"
REST_H1_URL_PROTOBUF="https://10.184.0.2:8080/rest/protobuf/http1/order/hundred"
REST_H2_URL_JSON="https://10.184.0.2:8081/rest/json/http2/order/hundred"
REST_H2_URL_PROTOBUF="https://10.184.0.2:8081/rest/protobuf/http2/order/hundred"
GRPC_TARGET="10.184.0.2:50051"
GRPC_CALL="orderexperiment.OrderExperimentService/GetOrderHundred"

run_h2load_h1() {
  local url="$1" m="$2" outfile="$3"
  h2load --h1 -c 1 -n 5000 -m "${m}" --log-file="${outfile}" "${url}"
}

run_h2load_h2() {
  local url="$1" m="$2" outfile="$3"
  h2load -c 1 -n 5000 -m "${m}" --log-file="${outfile}" "${url}"
}

run_ghz() {
  local m="$1" outfile="$2"
  ghz --skipTLS \
    --proto "${PROTO_PATH}" \
    --call "${GRPC_CALL}" \
    --connections 1 \
    --concurrency "${m}" \
    --total 5000 \
    --format csv \
    --output "${outfile}" \
    "${GRPC_TARGET}"
}

for m in 10 50 100 500 1000; do
  echo "=== Skenario Tingkat Multiplexing — ${m} Permintaan ==="
  run_h2load_h1 "${REST_H1_URL_JSON}"     "${m}" "${RESULTS_DIR}/rest-h1-json-m${m}.log"
  run_h2load_h1 "${REST_H1_URL_PROTOBUF}" "${m}" "${RESULTS_DIR}/rest-h1-protobuf-m${m}.log"
  run_h2load_h2 "${REST_H2_URL_JSON}"     "${m}" "${RESULTS_DIR}/rest-h2-json-m${m}.log"
  run_h2load_h2 "${REST_H2_URL_PROTOBUF}" "${m}" "${RESULTS_DIR}/rest-h2-protobuf-m${m}.log"
  run_ghz "${m}" "${RESULTS_DIR}/grpc-m${m}.csv"
done

echo "Skenario Tingkat Multiplexing Selesai — 25 file dihasilkan."