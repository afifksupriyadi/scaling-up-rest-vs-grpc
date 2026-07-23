"""
check_anomaly.py

One-off diagnostic: how many rows have a negative "elapsed" value for
the REST HTTP/2 + Protobuf combination, and what do those rows look
like? Negative elapsed time is not physically possible (finish time
before start time), so this checks whether it's a rare outlier or a
widespread issue with the bzm-HTTP2 Sampler's timing.
"""

from pathlib import Path

import pandas as pd

SCRIPT_DIR = Path(__file__).parent
RESULTS_CSV = SCRIPT_DIR / "results" / "results-aggregate.csv"

TARGET_LABELS = ["rest-http2-protobuf-small", "rest-http2-protobuf-large"]

df = pd.read_csv(RESULTS_CSV, low_memory=False)

for label in TARGET_LABELS:
    subset = df[df["label"] == label]
    total = len(subset)
    negative = subset[subset["elapsed"] < 0]

    print(f"\n=== {label} ===")
    print(f"Total baris: {total}")
    print(f"Baris dengan elapsed negatif: {len(negative)} ({len(negative) / total * 100:.2f}%)")

    if not negative.empty:
        print("Contoh 5 baris pertama yang negatif:")
        print(negative[["timeStamp", "elapsed", "Latency", "Connect", "responseCode", "success"]].head())