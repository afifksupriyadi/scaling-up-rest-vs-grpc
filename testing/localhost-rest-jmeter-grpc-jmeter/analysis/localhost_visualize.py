"""
localhost_visualize.py

Reads raw JMeter sample results (rest-grpc-report_localhost.csv,
containing one row per request across all 20 Thread Groups: 15 GET +
5 POST) and produces:
  1. GET: a throughput bar chart (x-axis: protocol/format combination,
     grouped by data size: small/medium/large).
  2. GET: a response time summary table (mean, p99, min, max) per
     combination and data size.
  3. GET: response time over-time overlay charts (one per data size, one
     line per protocol/format combination).
  4. POST: a throughput bar chart (one bar per combination, no size
     grouping, since POST always sends exactly one entity).
  5. POST: a response time summary table per combination.
  6. POST: a response time over-time overlay chart (one line per
     combination).

Run with: python3 localhost_visualize.py
Expects rest-grpc-report_localhost.csv to sit in ../results/ relative
to this script.
"""

from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

# Paths are resolved relative to this script's own location, not the
# current working directory, so this can be run from anywhere. The script
# and its outputs live together in analysis/, while the raw JMeter output
# lives in the sibling results/ folder.
SCRIPT_DIR = Path(__file__).parent
RESULTS_CSV = SCRIPT_DIR.parent / "results" / "rest-grpc-report_localhost.csv"
OUTPUT_DIR = SCRIPT_DIR

# Shown under every chart's main title, identifying which scenario this
# specific chart was generated from, since the three scenarios (localhost,
# vps-jmeter-only, vps-jmeter-rest-ghz-grpc) share identical output
# filenames and must be told apart when read in isolation (e.g. pasted
# into the thesis document one scenario at a time).
SCENARIO_LABEL = "Lingkungan Pengujian: Localhost (REST & gRPC via JMeter)"

# ---------------------------------------------------------------------------
# GET (15 Thread Groups: 5 combinations x 3 data sizes)
# ---------------------------------------------------------------------------

# Human-readable combination names, used for chart titles/legends.
COMBO_NAMES = {
    "rest-http1-json": "REST HTTP/1.1 + JSON",
    "rest-http1-protobuf": "REST HTTP/1.1 + Protobuf",
    "rest-http2-json": "REST HTTP/2 + JSON",
    "rest-http2-protobuf": "REST HTTP/2 + Protobuf",
    "grpc": "gRPC",
}

# Maps the EXACT label string that appears in the CSV to
# (human-readable combination name, data size).
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
    "grpc-small": ("gRPC", "small"),
    "grpc-medium": ("gRPC", "medium"),
    "grpc-large": ("gRPC", "large"),
}

# ---------------------------------------------------------------------------
# POST (5 Thread Groups: 5 combinations, no data size variation)
# ---------------------------------------------------------------------------

# Maps the EXACT POST label string to its human-readable combination name.
# No size axis here, since POST always sends exactly one Student entity.
POST_LABEL_MAP = {
    "rest-http1-json-post": "REST HTTP/1.1 + JSON",
    "rest-http1-protobuf-post": "REST HTTP/1.1 + Protobuf",
    "rest-http2-json-post": "REST HTTP/2 + JSON",
    "rest-http2-protobuf-post": "REST HTTP/2 + Protobuf",
    "grpc-post": "gRPC",
}

# ---------------------------------------------------------------------------
# Shared styling
# ---------------------------------------------------------------------------

# Mutually distinct palette for combination lines/bars, chosen so no two
# combinations read as visually similar even when placed side by side
# (an earlier palette had three warm brownish tones that blurred together).
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

def load_raw_samples() -> pd.DataFrame:
    """Load the raw per-request JMeter results and print a quick sanity check of what was actually found, so mismatches are caught early instead of failing silently deeper in the script."""
    if not RESULTS_CSV.exists():
        raise FileNotFoundError(
            f"Tidak ketemu {RESULTS_CSV}. Pastikan rest-grpc-report_localhost.csv "
            f"sudah ada di folder results/, sejajar dengan folder analysis/ ini."
        )

    df = pd.read_csv(RESULTS_CSV, low_memory=False)
    print(f"Loaded {len(df)} rows from {RESULTS_CSV.name}")
    print(f"Columns found: {list(df.columns)}")
    print(f"Labels found: {sorted(df['label'].unique())}")
    return df


def parse_label(label: str) -> tuple[str, str]:
    """Look up a GET sampler label's (combination name, data size)."""
    if label not in LABEL_MAP:
        raise ValueError(f"Label GET tidak dikenali: {label!r}. Cek LABEL_MAP di atas.")
    return LABEL_MAP[label]


def _exclude_negative_elapsed(group: pd.DataFrame, label: str) -> tuple[pd.Series, int]:
    """Shared helper: drop rows with a negative elapsed value (a known bzm-HTTP2 Sampler timing bug, not a real negative response time), printing how many rows were dropped so it's visible, not silent."""
    total_count = len(group)
    valid = group[group["elapsed"] >= 0]
    excluded = total_count - len(valid)
    if excluded > 0:
        print(f"  [{label}] excluded {excluded}/{total_count} rows with negative elapsed "
              f"({excluded / total_count * 100:.2f}%) from response time stats")
    return valid["elapsed"], excluded


def compute_summary(df: pd.DataFrame) -> pd.DataFrame:
    """Compute mean/median/p90/p95/p99/min/max/throughput per GET label."""
    rows = []
    for label, group in df.groupby("label"):
        if label not in LABEL_MAP:
            continue  # skip POST labels, handled separately
        combo, size = parse_label(label)

        elapsed, excluded = _exclude_negative_elapsed(group, label)
        total_count = len(group)

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
    """GET bar chart: x-axis = protocol/format combination, grouped into three bars per combination, one per data size (small/medium/large)."""
    combos = list(COMBO_NAMES.values())
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
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")

    ax.legend(title="Ukuran Data", loc="upper left", bbox_to_anchor=(1.01, 1.0))
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "get_throughput.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'get_throughput.png'}")

def print_response_time_table(summary: pd.DataFrame):
    """Print (and save as CSV) the GET response time summary table."""
    table = summary[["combo", "size", "mean", "p99", "min", "max"]].copy()
    table["size"] = table["size"].map(SIZE_LABELS)
    table["_order"] = summary["size"].map({s: i for i, s in enumerate(SIZE_ORDER)})
    table = table.sort_values(["combo", "_order"]).drop(columns="_order")
    table.columns = ["Kombinasi", "Ukuran Data", "Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]
    for col in ["Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]:
        table[col] = table[col].round(1)

    print("\nRingkasan Response Time (GET):")
    print(table.to_string(index=False))

    OUTPUT_DIR.mkdir(exist_ok=True)
    table.to_csv(OUTPUT_DIR / "get_response_time_summary.csv", index=False)
    print(f"Saved {OUTPUT_DIR / 'get_response_time_summary.csv'}")


def plot_response_time_overlay(df: pd.DataFrame, size: str):
    """GET overlay chart: one line per combination, x-axis = time since test start (bucketed into 1-second windows), y-axis = mean response time. Markers are drawn on each point, matching the reference literature's chart style."""
    fig, ax = plt.subplots(figsize=(10, 6))

    for prefix, combo in COMBO_NAMES.items():
        matching_labels = [lbl for lbl, (c, s) in LABEL_MAP.items() if c == combo and s == size]
        if not matching_labels:
            continue

        subset = df[df["label"].isin(matching_labels)].copy()
        subset = subset[subset["elapsed"] >= 0]
        if subset.empty:
            continue

        start = subset["timeStamp"].min()
        subset["bucket_s"] = (subset["timeStamp"] - start) // 500 / 2
        bucketed = subset.groupby("bucket_s")["elapsed"].mean()

        ax.plot(bucketed.index, bucketed.values, label=combo, color=COMBO_COLORS[combo],
                linewidth=1.5, marker="o", markersize=4)

    ax.set_xlabel("Waktu Sejak Pengujian Dimulai (detik)")
    ax.set_ylabel("Response Time Rata-rata (ms)")
    fig.suptitle(f"GET - Response Time ({SIZE_LABELS[size]})", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_DIR.mkdir(exist_ok=True)
    filename = f"get_response_time_overlay_{size}.png"
    fig.savefig(OUTPUT_DIR / filename, dpi=150)
    print(f"Saved {OUTPUT_DIR / filename}")


def compute_post_summary(df: pd.DataFrame) -> pd.DataFrame:
    """Compute mean/median/p90/p95/p99/min/max/throughput per POST label. Same statistics as GET, minus the size axis, since POST always sends exactly one Student entity per request."""
    rows = []
    for label, group in df.groupby("label"):
        if label not in POST_LABEL_MAP:
            continue  # skip GET labels, handled separately
        combo = POST_LABEL_MAP[label]

        elapsed, excluded = _exclude_negative_elapsed(group, label)
        total_count = len(group)

        duration_s = (group["timeStamp"] + group["elapsed"].clip(lower=0)).max() - group["timeStamp"].min()
        duration_s = max(duration_s, 1) / 1000
        throughput = total_count / duration_s

        rows.append({
            "label": label,
            "combo": combo,
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


def plot_post_throughput_bar(summary: pd.DataFrame):
    """POST bar chart: one bar per combination, no size grouping, since there is only one condition (one entity per request)."""
    combos = list(POST_LABEL_MAP.values())
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
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    ax.set_ylim(top=max(vals) * 1.15)
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "post_throughput.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'post_throughput.png'}")


def print_post_response_time_table(summary: pd.DataFrame):
    """Print (and save as CSV) the POST response time summary table."""
    table = summary[["combo", "mean", "p99", "min", "max"]].copy()
    table = table.sort_values("combo")
    table.columns = ["Kombinasi", "Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]
    for col in ["Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]:
        table[col] = table[col].round(1)

    print("\nRingkasan Response Time (POST):")
    print(table.to_string(index=False))

    OUTPUT_DIR.mkdir(exist_ok=True)
    table.to_csv(OUTPUT_DIR / "post_response_time_summary.csv", index=False)
    print(f"Saved {OUTPUT_DIR / 'post_response_time_summary.csv'}")


def plot_post_response_time_overlay(df: pd.DataFrame):
    """POST overlay chart: one line per combination, with markers on each point. Only one chart is needed here (unlike GET's three), since POST has no data size axis."""
    fig, ax = plt.subplots(figsize=(10, 6))

    for label, combo in POST_LABEL_MAP.items():
        subset = df[df["label"] == label].copy()
        subset = subset[subset["elapsed"] >= 0]
        if subset.empty:
            continue

        start = subset["timeStamp"].min()
        subset["bucket_s"] = (subset["timeStamp"] - start) // 500 / 2
        bucketed = subset.groupby("bucket_s")["elapsed"].mean()

        ax.plot(bucketed.index, bucketed.values, label=combo, color=COMBO_COLORS[combo],
                linewidth=1.5, marker="o", markersize=4)

    ax.set_xlabel("Waktu Sejak Pengujian Dimulai (detik)")
    ax.set_ylabel("Response Time Rata-rata (ms)")
    fig.suptitle("POST - Response Time", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_DIR.mkdir(exist_ok=True)
    fig.savefig(OUTPUT_DIR / "post_response_time_overlay.png", dpi=150)
    print(f"Saved {OUTPUT_DIR / 'post_response_time_overlay.png'}")


def main():
    df = load_raw_samples()

    get_summary = compute_summary(df)
    plot_throughput_bar(get_summary)
    print_response_time_table(get_summary)
    for size in SIZE_ORDER:
        plot_response_time_overlay(df, size)

    post_summary = compute_post_summary(df)
    plot_post_throughput_bar(post_summary)
    print_post_response_time_table(post_summary)
    plot_post_response_time_overlay(df)


if __name__ == "__main__":
    main()