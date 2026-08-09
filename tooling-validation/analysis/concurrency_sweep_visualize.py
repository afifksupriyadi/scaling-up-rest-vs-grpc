"""
concurrency_sweep_visualize.py

Reads Fortio JSON results (REST HTTP/1.1) and k6 CSV results (gRPC) across
six concurrency levels (1, 5, 10, 25, 50, 100) for the same Depth0-Large
endpoint, and produces two comparison charts plus a summary table.

Expects files named:
  ../fortio/rest-http1-depth0-large-c{N}.json
  ../k6/grpc-depth0-large-vu{N}.csv
for N in CONCURRENCY_LEVELS.

Run with: python3 concurrency_sweep_visualize.py
"""

import json
from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

SCRIPT_DIR = Path(__file__).parent
FORTIO_DIR = SCRIPT_DIR.parent / "fortio"
K6_DIR = SCRIPT_DIR.parent / "k6"
OUTPUT_DIR = SCRIPT_DIR / "concurrency-sweep"

CONCURRENCY_LEVELS = [1, 5, 10, 25, 50, 100]


def load_fortio_point(c: int) -> dict:
    """Load one Fortio JSON result (REST HTTP/1.1) and extract mean response time (ms) and throughput (req/s)."""
    path = FORTIO_DIR / f"rest-http1-depth0-large-c{c}.json"
    with open(path) as f:
        data = json.load(f)
    return {
        "concurrency": c,
        "mean_ms": data["DurationHistogram"]["Avg"] * 1000,
        "throughput": data["ActualQPS"],
    }


def load_k6_point(vu: int) -> dict:
    """Load one k6 CSV result (gRPC) and compute mean response time (ms) and throughput (req/s).
    - Throughput is derived from raw whole-second timestamps, which can be imprecise at low VU
      counts where the whole test finishes in well under a second (notably vu=1) - treat that
      specific point with some caution."""
    path = K6_DIR / f"grpc-depth0-large-vu{vu}.csv"
    df = pd.read_csv(path, low_memory=False)
    dur = df[df["metric_name"] == "grpc_req_duration"]["metric_value"]
    iters = df[df["metric_name"] == "iterations"]
    total_iterations = iters["metric_value"].sum()
    duration_s = max(iters["timestamp"].max() - iters["timestamp"].min(), 1)
    return {
        "concurrency": vu,
        "mean_ms": dur.mean(),
        "throughput": total_iterations / duration_s,
    }


def print_table(fortio_points, k6_points):
    """Print (and save as CSV) the combined summary table across all concurrency levels."""
    rows = []
    for f, k in zip(fortio_points, k6_points):
        rows.append({
            "Konkurensi": f["concurrency"],
            "REST Mean (ms)": round(f["mean_ms"], 2),
            "REST Throughput (req/s)": round(f["throughput"], 1),
            "gRPC Mean (ms)": round(k["mean_ms"], 2),
            "gRPC Throughput (req/s)": round(k["throughput"], 1),
        })
    table = pd.DataFrame(rows)
    print(table.to_string(index=False))
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    table.to_csv(OUTPUT_DIR / "concurrency_sweep_summary.csv", index=False)
    print(f"Saved {OUTPUT_DIR / 'concurrency_sweep_summary.csv'}")


def plot_response_time(fortio_points, k6_points):
    """Line chart: mean response time vs concurrency, one line per tool/protocol, log-scale x-axis."""
    fig, ax = plt.subplots(figsize=(10, 6))
    ax.plot([p["concurrency"] for p in fortio_points], [p["mean_ms"] for p in fortio_points],
            label="REST HTTP/1.1 (Fortio)", color="#C9A227", marker="o", linewidth=1.8)
    ax.plot([p["concurrency"] for p in k6_points], [p["mean_ms"] for p in k6_points],
            label="gRPC (k6)", color="#4F9D8C", marker="s", linewidth=1.8)
    ax.set_xscale("log")
    ax.set_xticks(CONCURRENCY_LEVELS)
    ax.set_xticklabels(CONCURRENCY_LEVELS)
    ax.set_xlabel("Jumlah VU / Koneksi Konkuren")
    ax.set_ylabel("Response Time Rata-rata (ms)")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.suptitle("Response Time vs Konkurensi - Depth0-Large", fontsize=13, fontweight="bold")
    fig.tight_layout()
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    fig.savefig(OUTPUT_DIR / "response_time_vs_concurrency.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'response_time_vs_concurrency.png'}")


def plot_throughput(fortio_points, k6_points):
    """Line chart: throughput vs concurrency, one line per tool/protocol, log-scale x-axis."""
    fig, ax = plt.subplots(figsize=(10, 6))
    ax.plot([p["concurrency"] for p in fortio_points], [p["throughput"] for p in fortio_points],
            label="REST HTTP/1.1 (Fortio)", color="#C9A227", marker="o", linewidth=1.8)
    ax.plot([p["concurrency"] for p in k6_points], [p["throughput"] for p in k6_points],
            label="gRPC (k6)", color="#4F9D8C", marker="s", linewidth=1.8)
    ax.set_xscale("log")
    ax.set_xticks(CONCURRENCY_LEVELS)
    ax.set_xticklabels(CONCURRENCY_LEVELS)
    ax.set_xlabel("Jumlah VU / Koneksi Konkuren")
    ax.set_ylabel("Throughput (request/detik)")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.suptitle("Throughput vs Konkurensi - Depth0-Large", fontsize=13, fontweight="bold")
    fig.tight_layout()
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    fig.savefig(OUTPUT_DIR / "throughput_vs_concurrency.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'throughput_vs_concurrency.png'}")


def main():
    fortio_points = [load_fortio_point(c) for c in CONCURRENCY_LEVELS]
    k6_points = [load_k6_point(vu) for vu in CONCURRENCY_LEVELS]

    print_table(fortio_points, k6_points)
    plot_response_time(fortio_points, k6_points)
    plot_throughput(fortio_points, k6_points)


if __name__ == "__main__":
    main()
