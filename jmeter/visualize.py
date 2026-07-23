"""
visualize.py

Reads raw JMeter sample results (results-aggregate.csv, containing one row
per request across all 10 Thread Groups) and produces:
  1. A throughput bar chart (x-axis: protocol/format combination,
     grouped by data size).
  2. A response time summary table (mean, p99, min, max) per combination.
  3. Response time over-time overlay charts (one per data size, one line
     per protocol/format combination).

Run with: python3 visualize.py
Expects results-aggregate.csv to sit in ./results/ next to this script.
"""

from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

# Paths are resolved relative to this script's own location, not the
# current working directory, so this can be run from anywhere.
SCRIPT_DIR = Path(__file__).parent
RESULTS_CSV = SCRIPT_DIR / "results" / "results-aggregate.csv"
OUTPUT_DIR = SCRIPT_DIR / "output"

# Human-readable combination names, used for chart titles/legends.
COMBO_NAMES = {
    "rest-http1-json": "REST HTTP/1.1 + JSON",
    "rest-http1-protobuf": "REST HTTP/1.1 + Protobuf",
    "rest-http2-json": "REST HTTP/2 + JSON",
    "rest-http2-protobuf": "REST HTTP/2 + Protobuf",
    "grpc": "gRPC",
}

# Maps the EXACT label string that appears in the CSV to
# (human-readable combination name, data size). Using an explicit map
# instead of a regex pattern because the two gRPC samplers ended up
# named "Get Student"/"Get Students" (auto-named by the plugin when the
# method was picked from the "Listing..." dropdown) instead of the
# planned "grpc-small"/"grpc-large".
LABEL_MAP = {
    "rest-http1-json-small": ("REST HTTP/1.1 + JSON", "small"),
    "rest-http1-json-large": ("REST HTTP/1.1 + JSON", "large"),
    "rest-http1-protobuf-small": ("REST HTTP/1.1 + Protobuf", "small"),
    "rest-http1-protobuf-large": ("REST HTTP/1.1 + Protobuf", "large"),
    "rest-http2-json-small": ("REST HTTP/2 + JSON", "small"),
    "rest-http2-json-large": ("REST HTTP/2 + JSON", "large"),
    "rest-http2-protobuf-small": ("REST HTTP/2 + Protobuf", "small"),
    "rest-http2-protobuf-large": ("REST HTTP/2 + Protobuf", "large"),
    "Get Student": ("gRPC", "small"),
    "Get Students": ("gRPC", "large"),
}

# Soft, muted colors (used consistently across every chart in this script).
COLOR_SMALL = "#3B82A0"  # dipakai utk seri "1 Entri"
COLOR_LARGE = "#E0526E"  # dipakai utk seri "100 Entri"
COMBO_COLORS = {
    "REST HTTP/1.1 + JSON": "#3B82A0",
    "REST HTTP/1.1 + Protobuf": "#7C6BC4",
    "REST HTTP/2 + JSON": "#E0526E",
    "REST HTTP/2 + Protobuf": "#F0954A",
    "gRPC": "#2FA88A",
}

SIZE_LABELS = {"small": "1 Entri", "large": "100 Entri"}


def load_raw_samples() -> pd.DataFrame:
    """Load the raw per-request JMeter results and print a quick sanity
    check of what was actually found, so mismatches are caught early
    instead of failing silently deeper in the script."""
    if not RESULTS_CSV.exists():
        raise FileNotFoundError(
            f"Tidak ketemu {RESULTS_CSV}. Pastikan file sudah di-copy ke "
            f"folder jmeter/results/ sesuai struktur yang disepakati."
        )

    df = pd.read_csv(RESULTS_CSV, low_memory=False)
    print(f"Loaded {len(df)} rows from {RESULTS_CSV.name}")
    print(f"Columns found: {list(df.columns)}")
    print(f"Labels found: {sorted(df['label'].unique())}")
    return df


def parse_label(label: str) -> tuple[str, str]:
    """Look up a sampler label's (combination name, data size)."""
    if label not in LABEL_MAP:
        raise ValueError(f"Label tidak dikenali: {label!r}. Cek LABEL_MAP di atas.")
    return LABEL_MAP[label]


def compute_summary(df: pd.DataFrame) -> pd.DataFrame:
    """Compute mean/median/p90/p95/p99/min/max/throughput per label.

    Rows with a negative "elapsed" value are excluded from the response
    time statistics (mean/median/percentiles/min/max), since a negative
    elapsed time is not physically possible and reflects a timing bug in
    the bzm-HTTP2 Sampler plugin, not an actual slow/fast response. These
    rows are still counted toward throughput, since the request itself
    completed successfully (responseCode 200); only its recorded timing
    value is unusable.
    """
    rows = []
    for label, group in df.groupby("label"):
        combo, size = parse_label(label)

        total_count = len(group)
        valid = group[group["elapsed"] >= 0]
        excluded = total_count - len(valid)
        if excluded > 0:
            print(f"  [{label}] excluded {excluded}/{total_count} rows with negative elapsed "
                  f"({excluded / total_count * 100:.2f}%) from response time stats")

        elapsed = valid["elapsed"]

        duration_s = (group["timeStamp"] + group["elapsed"].clip(lower=0)).max() - group["timeStamp"].min()
        duration_s = max(duration_s, 1) / 1000
        throughput = total_count / duration_s

        rows.append({
            "label": label,
            "combo": combo,
            "size": size,
            "mean": elapsed.mean(),
            "median": elapsed.median(),
            "p90": elapsed.quantile(0.90),
            "p95": elapsed.quantile(0.95),
            "p99": elapsed.quantile(0.99),
            "min": elapsed.min(),
            "max": elapsed.max(),
            "throughput": throughput,
            "error_rate": (group["success"] == False).mean() * 100,  # noqa: E712
            "excluded_rows": excluded,
        })
    return pd.DataFrame(rows)

def plot_throughput_bar(summary: pd.DataFrame):
    """Bar chart: x-axis = protocol/format combination, grouped by data size."""
    combos = list(COMBO_NAMES.values())
    x = range(len(combos))
    width = 0.35

    small_vals = [summary[(summary.combo == c) & (summary["size"] == "small")]["throughput"].sum() for c in combos]
    large_vals = [summary[(summary.combo == c) & (summary["size"] == "large")]["throughput"].sum() for c in combos]

    fig, ax = plt.subplots(figsize=(11, 6))
    bars_small = ax.bar([i - width / 2 for i in x], small_vals, width, label="1 Entri", color=COLOR_SMALL, edgecolor="white")
    bars_large = ax.bar([i + width / 2 for i in x], large_vals, width, label="100 Entri", color=COLOR_LARGE, edgecolor="white")

    ax.bar_label(bars_small, fmt="%.1f", fontsize=9, padding=3)
    ax.bar_label(bars_large, fmt="%.1f", fontsize=9, padding=3)

    ax.set_xticks(list(x))
    ax.set_xticklabels(combos, rotation=20, ha="right")
    ax.set_ylabel("Throughput (request/detik)")
    ax.set_title("Perbandingan Throughput per Kombinasi Protokol/Format")
    ax.set_ylim(top=max(small_vals + large_vals) * 1.12)

    # Legend dipindah ke luar area plot (kanan atas, di luar batas grafik),
    # biar tidak menutupi batang atau angka label manapun.
    ax.legend(title="Ukuran Data", loc="upper left", bbox_to_anchor=(1.01, 1.0))

    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "throughput.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'throughput.png'}")


def print_response_time_table(summary: pd.DataFrame):
    """Print (and save as CSV) the response time summary table."""
    table = summary[["combo", "size", "mean", "p99", "min", "max"]].copy()
    table["size"] = table["size"].map(SIZE_LABELS)
    table = table.sort_values(["combo", "size"])
    table.columns = ["Kombinasi", "Ukuran Data", "Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]
    for col in ["Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]:
        table[col] = table[col].round(1)

    print("\nRingkasan Response Time:")
    print(table.to_string(index=False))

    OUTPUT_DIR.mkdir(exist_ok=True)
    table.to_csv(OUTPUT_DIR / "response_time_summary.csv", index=False)
    print(f"Saved {OUTPUT_DIR / 'response_time_summary.csv'}")


def plot_response_time_overlay(df: pd.DataFrame, size: str):
    """Overlay chart: one line per combination, x-axis = time since test
    start (bucketed into 1-second windows), y-axis = mean response time."""
    fig, ax = plt.subplots(figsize=(10, 6))

    for prefix, combo in COMBO_NAMES.items():
        # Find whichever label(s) in LABEL_MAP resolve to this (combo, size),
        # since the gRPC labels don't follow the "<prefix>-<size>" pattern.
        matching_labels = [lbl for lbl, (c, s) in LABEL_MAP.items() if c == combo and s == size]
        if not matching_labels:
            continue

        subset = df[df["label"].isin(matching_labels)].copy()
        subset = subset[subset["elapsed"] >= 0]
        if subset.empty:
            continue

        start = subset["timeStamp"].min()
        subset["bucket_s"] = (subset["timeStamp"] - start) // 1000
        bucketed = subset.groupby("bucket_s")["elapsed"].mean()

        ax.plot(bucketed.index, bucketed.values, label=combo, color=COMBO_COLORS[combo], linewidth=2)

    ax.set_xlabel("Waktu Sejak Pengujian Dimulai (detik)")
    ax.set_ylabel("Response Time Rata-rata (ms)")
    size_label = "1 Entri" if size == "small" else "100 Entri"
    ax.set_title(f"Response Time Sepanjang Waktu Pengujian ({size_label})")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_DIR.mkdir(exist_ok=True)
    filename = f"response_time_overlay_{size}.png"
    fig.savefig(OUTPUT_DIR / filename, dpi=150)
    print(f"Saved {OUTPUT_DIR / filename}")


def main():
    df = load_raw_samples()
    summary = compute_summary(df)

    plot_throughput_bar(summary)
    print_response_time_table(summary)
    plot_response_time_overlay(df, "small")
    plot_response_time_overlay(df, "large")


if __name__ == "__main__":
    main()