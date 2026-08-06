"""
vps_visualize-jmeter-ghz-combined.py

Reads two separate raw data sources and produces combined REST+gRPC
visualizations and summary tables:
  1. REST (16 Thread Groups via JMeter): rest-only-report_vps.csv in
     ../results/, covering 4 protocol/format combinations x 3 data
     sizes (GET) plus 4 combinations for POST.
  2. gRPC (4 separate runs via ghz, used instead of JMeter's gRPC
     plugin due to severe, well-documented measurement overhead found
     during earlier diagnosis): JSON report files in ../results/, one
     per data size (GET) plus one for POST.

Run with: python3 vps_visualize-jmeter-ghz-combined.py
"""

import json
from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

SCRIPT_DIR = Path(__file__).parent
RESULTS_DIR = SCRIPT_DIR.parent / "results"
REST_RESULTS_CSV = RESULTS_DIR / "rest-only-report_vps.csv"
OUTPUT_DIR = SCRIPT_DIR

# Shown under every chart's main title, identifying which scenario this
# specific chart was generated from, since all three scenarios share
# identical output filenames and must be told apart when read in isolation.
SCENARIO_LABEL = "Lingkungan Pengujian: VPS"
SERVER_LABEL = "Server (REST API dan gRPC): e2-standard-8, 8 vCPU, 32 GB RAM"
CLIENT_LABEL = "Client REST (JMeter): e2-standard-4, 4 vCPU, 16 GB RAM  |  Client gRPC (ghz): e2-standard-4, 4 vCPU, 16 GB RAM"

REST_COMBO_NAMES = {
    "rest-http1-json": "REST HTTP/1.1 + JSON",
    "rest-http1-protobuf": "REST HTTP/1.1 + Protobuf",
    "rest-http2-json": "REST HTTP/2 + JSON",
    "rest-http2-protobuf": "REST HTTP/2 + Protobuf",
}

LABEL_MAP = {
    "rest-http1-json-small": ("REST HTTP/1.1 + JSON", "small"),
    "rest-http1-json-medium": ("REST HTTP/1.1 + JSON", "medium"),
    "rest-http1-json-large": ("REST HTTP/1.1 + JSON", "large"),
    "rest-http1-protobuf-small": ("REST HTTP/1.1 + Protobuf", "small"),
    "rest-http1-protobuf-medium": ("REST HTTP/1.1 + Protobuf", "medium"),
    "rest-http1-protobuf-large": ("REST HTTP/1.1 + Protobuf", "large"),
    "rest-http2-json-small": ("REST HTTP/2 + JSON", "small"),
    "rest-http2-json-medium": ("REST HTTP/2 + JSON", "medium"),
    "rest-http2-json-large": ("REST HTTP/2 + JSON", "large"),
    "rest-http2-protobuf-small": ("REST HTTP/2 + Protobuf", "small"),
    "rest-http2-protobuf-medium": ("REST HTTP/2 + Protobuf", "medium"),
    "rest-http2-protobuf-large": ("REST HTTP/2 + Protobuf", "large"),
}

POST_LABEL_MAP = {
    "rest-http1-json-post": "REST HTTP/1.1 + JSON",
    "rest-http1-protobuf-post": "REST HTTP/1.1 + Protobuf",
    "rest-http2-json-post": "REST HTTP/2 + JSON",
    "rest-http2-protobuf-post": "REST HTTP/2 + Protobuf",
}

# ghz JSON reports live in the same results/ folder as the REST CSV now,
# not a separate ghz/ folder, since this scenario keeps everything for
# one test run together.
GHZ_GET_FILES = {
    "small": RESULTS_DIR / "grpc-small-report.json",
    "medium": RESULTS_DIR / "grpc-medium-report.json",
    "large": RESULTS_DIR / "grpc-large-report.json",
}
GHZ_POST_FILE = RESULTS_DIR / "grpc-post-report.json"

COMBO_COLORS = {
    "REST HTTP/1.1 + JSON": "#A9764F",
    "REST HTTP/1.1 + Protobuf": "#7B6FA6",
    "REST HTTP/2 + JSON": "#C74B5C",
    "REST HTTP/2 + Protobuf": "#C9A227",
    "gRPC": "#4F9D8C",
}

SIZE_LABELS = {"small": "1 Entri", "medium": "100 Entri", "large": "1000 Entri"}
SIZE_ORDER = ["small", "medium", "large"]
SIZE_COLORS = {"small": "#AEBFCF", "medium": "#7089A3", "large": "#3D5872"}
ALL_COMBO_ORDER = list(REST_COMBO_NAMES.values()) + ["gRPC"]


def load_rest_samples() -> pd.DataFrame:
    """Load the raw per-request JMeter results (REST only, 16 Thread Groups) and print a quick sanity check of what was actually found."""
    if not REST_RESULTS_CSV.exists():
        raise FileNotFoundError(f"Tidak ketemu {REST_RESULTS_CSV}.")
    df = pd.read_csv(REST_RESULTS_CSV, low_memory=False)
    print(f"[REST] Loaded {len(df)} rows from {REST_RESULTS_CSV.name}")
    print(f"[REST] Labels found: {sorted(df['label'].unique())}")
    return df


def _exclude_negative_elapsed(group: pd.DataFrame, label: str) -> tuple[pd.Series, int]:
    """Shared helper: drop rows with a negative elapsed value (a known bzm-HTTP2 Sampler timing bug), printing how many rows were dropped."""
    total_count = len(group)
    valid = group[group["elapsed"] >= 0]
    excluded = total_count - len(valid)
    if excluded > 0:
        print(f"  [{label}] excluded {excluded}/{total_count} rows with negative elapsed")
    return valid["elapsed"], excluded


def compute_rest_get_summary(df: pd.DataFrame) -> pd.DataFrame:
    """Compute mean/p99/min/max/throughput/error_rate per REST GET label."""
    rows = []
    for label, group in df.groupby("label"):
        if label not in LABEL_MAP:
            continue
        combo, size = LABEL_MAP[label]
        elapsed, _ = _exclude_negative_elapsed(group, label)
        total_count = len(group)
        duration_s = max((group["timeStamp"] + group["elapsed"].clip(lower=0)).max() - group["timeStamp"].min(), 1) / 1000
        rows.append({
            "combo": combo, "size": size,
            "mean": elapsed.mean(), "p99": elapsed.quantile(0.99),
            "min": elapsed.min(), "max": elapsed.max(),
            "throughput": total_count / duration_s,
            "error_rate": (group["success"] == False).mean() * 100,  # noqa: E712
        })
    return pd.DataFrame(rows)


def compute_rest_post_summary(df: pd.DataFrame) -> pd.DataFrame:
    """Compute mean/p99/min/max/throughput/error_rate per REST POST label."""
    rows = []
    for label, group in df.groupby("label"):
        if label not in POST_LABEL_MAP:
            continue
        combo = POST_LABEL_MAP[label]
        elapsed, _ = _exclude_negative_elapsed(group, label)
        total_count = len(group)
        duration_s = max((group["timeStamp"] + group["elapsed"].clip(lower=0)).max() - group["timeStamp"].min(), 1) / 1000
        rows.append({
            "combo": combo,
            "mean": elapsed.mean(), "p99": elapsed.quantile(0.99),
            "min": elapsed.min(), "max": elapsed.max(),
            "throughput": total_count / duration_s,
            "error_rate": (group["success"] == False).mean() * 100,  # noqa: E712
        })
    return pd.DataFrame(rows)


def load_ghz_details(path: Path) -> pd.DataFrame:
    """Load a single ghz JSON report and return its per-request details as a DataFrame comparable to JMeter's (timestamp in ms since epoch, elapsed in ms, success)."""
    if not path.exists():
        raise FileNotFoundError(f"Tidak ketemu {path}. Pastikan berkas JSON hasil ghz sudah ada di results/.")
    with open(path) as f:
        report = json.load(f)
    rows = []
    for d in report["details"]:
        rows.append({
            "timeStamp": pd.Timestamp(d["timestamp"]).value // 1_000_000,
            "elapsed": d["latency"] / 1_000_000,
            "success": d["status"] == "OK",
        })
    df = pd.DataFrame(rows)
    print(f"[gRPC] Loaded {len(df)} rows from {path.name}")
    return df


def compute_ghz_summary(df: pd.DataFrame) -> dict:
    """Compute mean/p99/min/max/throughput/error_rate for a single ghz report's DataFrame."""
    elapsed = df["elapsed"]
    duration_s = max(df["timeStamp"].max() - df["timeStamp"].min(), 1) / 1000
    return {
        "mean": elapsed.mean(), "p99": elapsed.quantile(0.99),
        "min": elapsed.min(), "max": elapsed.max(),
        "throughput": len(df) / duration_s,
        "error_rate": (df["success"] == False).mean() * 100,  # noqa: E712
    }


def build_combined_get_summary(rest_summary: pd.DataFrame, ghz_dfs: dict) -> pd.DataFrame:
    """Combine REST's per-label summary with gRPC's per-size summary (from ghz) into one unified table covering all 5 combinations x 3 sizes."""
    rows = rest_summary.to_dict("records")
    for size, df in ghz_dfs.items():
        rows.append({"combo": "gRPC", "size": size, **compute_ghz_summary(df)})
    return pd.DataFrame(rows)


def plot_combined_throughput_bar(summary: pd.DataFrame):
    """GET bar chart: x-axis = combination (4 REST + gRPC), grouped into three bars per combination, one per data size."""
    combos = ALL_COMBO_ORDER
    x = range(len(combos))
    width = 0.25
    offsets = {"small": -width, "medium": 0, "large": width}

    fig, ax = plt.subplots(figsize=(12, 6))
    for size in SIZE_ORDER:
        vals = [summary[(summary.combo == c) & (summary["size"] == size)]["throughput"].sum() for c in combos]
        positions = [i + offsets[size] for i in x]
        bars = ax.bar(positions, vals, width, label=SIZE_LABELS[size], color=SIZE_COLORS[size],
                       alpha=0.88, edgecolor="#E5E5E5", linewidth=0.8)
        ax.bar_label(bars, fmt="%.1f", fontsize=8, padding=3, rotation=90)

    ax.set_xticks(list(x))
    ax.set_xticklabels(combos, rotation=20, ha="right")
    ax.set_ylabel("Throughput (request/detik)")
    fig.suptitle("GET - Throughput", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.935, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    fig.text(0.5, 0.905, SERVER_LABEL, fontsize=8, color="black", ha="center")
    fig.text(0.5, 0.875, CLIENT_LABEL, fontsize=8, color="black", ha="center")
    ax.legend(title="Ukuran Data", loc="upper left", bbox_to_anchor=(1.01, 1.0))
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout(rect=[0, 0, 1, 0.85])
    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "get_throughput.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'get_throughput.png'}")


def print_combined_get_table(summary: pd.DataFrame):
    """Print (and save as CSV) the combined GET response time summary table."""
    table = summary[["combo", "size", "mean", "p99", "min", "max"]].copy()
    table["_order"] = table["size"].map({s: i for i, s in enumerate(SIZE_ORDER)})
    table = table.sort_values(["combo", "_order"])
    table["size"] = table["size"].map(SIZE_LABELS)
    table = table.drop(columns="_order")
    table.columns = ["Kombinasi", "Ukuran Data", "Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]
    for col in ["Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]:
        table[col] = table[col].round(1)
    print("\nRingkasan Response Time (GET):")
    print(table.to_string(index=False))
    OUTPUT_DIR.mkdir(exist_ok=True)
    table.to_csv(OUTPUT_DIR / "get_response_time_summary.csv", index=False)
    print(f"Saved {OUTPUT_DIR / 'get_response_time_summary.csv'}")


def plot_combined_response_time_overlay(rest_df: pd.DataFrame, ghz_dfs: dict, size: str):
    """GET overlay chart: one line per combination (4 REST + 1 gRPC)."""
    fig, ax = plt.subplots(figsize=(10, 6))

    for prefix, combo in REST_COMBO_NAMES.items():
        matching_labels = [lbl for lbl, (c, s) in LABEL_MAP.items() if c == combo and s == size]
        subset = rest_df[rest_df["label"].isin(matching_labels)].copy()
        subset = subset[subset["elapsed"] >= 0]
        if subset.empty:
            continue
        start = subset["timeStamp"].min()
        subset["bucket_s"] = (subset["timeStamp"] - start) // 1000
        bucketed = subset.groupby("bucket_s")["elapsed"].mean()
        ax.plot(bucketed.index, bucketed.values, label=combo, color=COMBO_COLORS[combo], linewidth=1.5, marker="o", markersize=4)

    ghz_df = ghz_dfs[size].copy()
    start = ghz_df["timeStamp"].min()
    ghz_df["bucket_s"] = (ghz_df["timeStamp"] - start) // 1000
    bucketed = ghz_df.groupby("bucket_s")["elapsed"].mean()
    ax.plot(bucketed.index, bucketed.values, label="gRPC", color=COMBO_COLORS["gRPC"], linewidth=1.5, marker="o", markersize=4)

    ax.set_xlabel("Waktu Sejak Pengujian Dimulai (detik)")
    ax.set_ylabel("Response Time Rata-rata (ms)")
    fig.suptitle(f"GET - Response Time ({SIZE_LABELS[size]})", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.935, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    fig.text(0.5, 0.905, SERVER_LABEL, fontsize=8, color="black", ha="center")
    fig.text(0.5, 0.875, CLIENT_LABEL, fontsize=8, color="black", ha="center")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout(rect=[0, 0, 1, 0.85])
    OUTPUT_DIR.mkdir(exist_ok=True)
    filename = f"get_response_time_overlay_{size}.png"
    fig.savefig(OUTPUT_DIR / filename, dpi=150)
    print(f"Saved {OUTPUT_DIR / filename}")


def build_combined_post_summary(rest_summary: pd.DataFrame, ghz_post_df: pd.DataFrame) -> pd.DataFrame:
    """Combine REST POST summary with gRPC POST summary (from ghz) into one unified table."""
    rows = rest_summary.to_dict("records")
    rows.append({"combo": "gRPC", **compute_ghz_summary(ghz_post_df)})
    return pd.DataFrame(rows)


def plot_combined_post_throughput_bar(summary: pd.DataFrame):
    """POST bar chart: one bar per combination (4 REST + 1 gRPC), no size grouping."""
    combos = ALL_COMBO_ORDER
    x = range(len(combos))
    vals = [summary[summary.combo == c]["throughput"].sum() for c in combos]
    colors = [COMBO_COLORS[c] for c in combos]

    fig, ax = plt.subplots(figsize=(9, 6))
    bars = ax.bar(x, vals, color=colors, alpha=0.88, edgecolor="#E5E5E5", linewidth=0.8)
    ax.bar_label(bars, fmt="%.1f", fontsize=9, padding=3)
    ax.set_xticks(list(x))
    ax.set_xticklabels(combos, rotation=20, ha="right")
    ax.set_ylabel("Throughput (request/detik)")
    fig.suptitle("POST - Throughput", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.935, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    fig.text(0.5, 0.905, SERVER_LABEL, fontsize=8, color="black", ha="center")
    fig.text(0.5, 0.875, CLIENT_LABEL, fontsize=8, color="black", ha="center")
    ax.set_ylim(top=max(vals) * 1.15)
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout(rect=[0, 0, 1, 0.85])
    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "post_throughput.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'post_throughput.png'}")


def print_combined_post_table(summary: pd.DataFrame):
    """Print (and save as CSV) the combined POST response time summary table."""
    table = summary[["combo", "mean", "p99", "min", "max"]].copy().sort_values("combo")
    table.columns = ["Kombinasi", "Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]
    for col in ["Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]:
        table[col] = table[col].round(1)
    print("\nRingkasan Response Time (POST):")
    print(table.to_string(index=False))
    OUTPUT_DIR.mkdir(exist_ok=True)
    table.to_csv(OUTPUT_DIR / "post_response_time_summary.csv", index=False)
    print(f"Saved {OUTPUT_DIR / 'post_response_time_summary.csv'}")


def plot_combined_post_response_time_overlay(rest_df: pd.DataFrame, ghz_post_df: pd.DataFrame):
    """POST overlay chart: one line per combination (4 REST + 1 gRPC)."""
    fig, ax = plt.subplots(figsize=(10, 6))

    for label, combo in POST_LABEL_MAP.items():
        subset = rest_df[rest_df["label"] == label].copy()
        subset = subset[subset["elapsed"] >= 0]
        if subset.empty:
            continue
        start = subset["timeStamp"].min()
        subset["bucket_s"] = (subset["timeStamp"] - start) // 1000
        bucketed = subset.groupby("bucket_s")["elapsed"].mean()
        ax.plot(bucketed.index, bucketed.values, label=combo, color=COMBO_COLORS[combo], linewidth=1.5, marker="o", markersize=4)

    ghz_post_df = ghz_post_df.copy()
    start = ghz_post_df["timeStamp"].min()
    ghz_post_df["bucket_s"] = (ghz_post_df["timeStamp"] - start) // 1000
    bucketed = ghz_post_df.groupby("bucket_s")["elapsed"].mean()
    ax.plot(bucketed.index, bucketed.values, label="gRPC", color=COMBO_COLORS["gRPC"], linewidth=1.5, marker="o", markersize=4)

    ax.set_xlabel("Waktu Sejak Pengujian Dimulai (detik)")
    ax.set_ylabel("Response Time Rata-rata (ms)")
    fig.suptitle("POST - Response Time", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.935, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    fig.text(0.5, 0.905, SERVER_LABEL, fontsize=8, color="black", ha="center")
    fig.text(0.5, 0.875, CLIENT_LABEL, fontsize=8, color="black", ha="center")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout(rect=[0, 0, 1, 0.85])
    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "post_response_time_overlay.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'post_response_time_overlay.png'}")


def main():
    rest_df = load_rest_samples()
    ghz_get_dfs = {size: load_ghz_details(path) for size, path in GHZ_GET_FILES.items()}
    ghz_post_df = load_ghz_details(GHZ_POST_FILE)

    rest_get_summary = compute_rest_get_summary(rest_df)
    combined_get_summary = build_combined_get_summary(rest_get_summary, ghz_get_dfs)
    plot_combined_throughput_bar(combined_get_summary)
    print_combined_get_table(combined_get_summary)
    for size in SIZE_ORDER:
        plot_combined_response_time_overlay(rest_df, ghz_get_dfs, size)

    rest_post_summary = compute_rest_post_summary(rest_df)
    combined_post_summary = build_combined_post_summary(rest_post_summary, ghz_post_df)
    plot_combined_post_throughput_bar(combined_post_summary)
    print_combined_post_table(combined_post_summary)
    plot_combined_post_response_time_overlay(rest_df, ghz_post_df)


if __name__ == "__main__":
    main()