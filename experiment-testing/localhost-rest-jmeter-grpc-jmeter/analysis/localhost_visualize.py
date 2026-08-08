"""
localhost_visualize.py (shape experiment)

Reads raw JMeter sample results (shape-experiment-report_localhost.csv,
containing one row per request across all 45 Thread Groups: 40 structural
shape/size-tier combinations + 5 Student baseline combinations, GET only)
and produces:
  1. Per-shape response time overlay charts (4 for compact, 5 for large,
     including students-large), one line per protocol/format combination.
  2. Two throughput bar charts (compact, large), x-axis = protocol/format
     combination, grouped bars = one per structural shape variant.
  3. Two response time summary tables (compact, large), saved as CSV.
  4. Two "gradation" charts (response time, throughput), showing gRPC's
     ratio against REST HTTP/2+Protobuf across the nesting-depth spectrum,
     for both compact and large tiers on the same axes.

Run with: python3 localhost_visualize.py
Expects shape-experiment-report_localhost.csv to sit in ../results/
relative to this script.
"""

from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

SCRIPT_DIR = Path(__file__).parent
RESULTS_CSV = SCRIPT_DIR.parent / "results" / "shape-experiment-report_localhost.csv"
OUTPUT_COMPACT_DIR = SCRIPT_DIR / "compact"
OUTPUT_LARGE_DIR = SCRIPT_DIR / "large"
OUTPUT_GRADATION_DIR = SCRIPT_DIR / "gradation"

SCENARIO_LABEL = "Lingkungan Pengujian: Localhost (Shape Experiment via JMeter)"

# ---------------------------------------------------------------------------
# Label vocabulary
# ---------------------------------------------------------------------------

# Same five combinations, same colors as the Student dataset's own scripts,
# so the two sets of charts read consistently when placed side by side in
# the thesis document.
COMBO_NAMES = {
    "rest-http1-json": "REST HTTP/1.1 + JSON",
    "rest-http1-protobuf": "REST HTTP/1.1 + Protobuf",
    "rest-http2-json": "REST HTTP/2 + JSON",
    "rest-http2-protobuf": "REST HTTP/2 + Protobuf",
    "grpc": "gRPC",
}

COMBO_COLORS = {
    "REST HTTP/1.1 + JSON": "#A9764F",
    "REST HTTP/1.1 + Protobuf": "#7B6FA6",
    "REST HTTP/2 + JSON": "#C74B5C",
    "REST HTTP/2 + Protobuf": "#C9A227",
    "gRPC": "#4F9D8C",
}

# The four structural variants, ordered by nesting depth (used for compact,
# and as the first four positions in large/gradation charts). "students" is
# always handled as a separate external baseline, never merged into this
# ordering, and only ever appears at the large tier.
SHAPE_ORDER_COMPACT = ["depth0", "depth1-wide", "depth3-narrow", "depth4-wide"]
SHAPE_ORDER_LARGE = ["depth0", "depth1-wide", "depth3-narrow", "depth4-wide", "students"]

SHAPE_LABELS = {
    "depth0": "Depth 0 (Flat)",
    "depth1-wide": "Depth 1 (Wide)",
    "depth3-narrow": "Depth 3 (Narrow)",
    "depth4-wide": "Depth 4 (Wide)",
    "students": "Student (Baseline Luar)",
}

# x-axis positions for the gradation charts: depth0..depth4-wide sit at 0-3,
# students-large sits at 5 (skipping 4), so the dashed break drawn at x=4
# visually separates it from the depth progression rather than implying
# it is a fifth step in the same sequence.
GRADATION_X_POS = {
    "depth0": 0, "depth1-wide": 1, "depth3-narrow": 2, "depth4-wide": 3, "students": 5,
}

SHAPE_BAR_COLORS = {
    "depth0": "#AEBFCF", "depth1-wide": "#8AA0B8", "depth3-narrow": "#5C7A96",
    "depth4-wide": "#3D5872", "students": "#C9A227",
}


def parse_label(label: str) -> tuple[str, str, str]:
    """Split a sampler label into (combination name, shape variant, size tier). Relies on COMBO_NAMES as the single source of truth for known prefixes, so build_label() below can never drift out of sync with this function."""
    for prefix, combo in COMBO_NAMES.items():
        if label.startswith(prefix + "-"):
            remainder = label[len(prefix) + 1:]
            if remainder.endswith("-compact"):
                return combo, remainder[: -len("-compact")], "compact"
            if remainder.endswith("-large"):
                return combo, remainder[: -len("-large")], "large"
    raise ValueError(f"Label tidak dikenali: {label!r}. Cek COMBO_NAMES di atas.")


def build_label(prefix: str, shape: str, tier: str) -> str:
    """Reconstruct a sampler label from its parts, the inverse of parse_label()."""
    return f"{prefix}-{shape}-{tier}"


def load_raw_samples() -> pd.DataFrame:
    """Load the raw per-request JMeter results and print a quick sanity check of what was actually found, so mismatches are caught early instead of failing silently deeper in the script."""
    if not RESULTS_CSV.exists():
        raise FileNotFoundError(
            f"Tidak ketemu {RESULTS_CSV}. Pastikan shape-experiment-report_localhost.csv "
            f"sudah ada di folder results/, sejajar dengan folder analysis/ ini."
        )
    df = pd.read_csv(RESULTS_CSV, low_memory=False)
    print(f"Loaded {len(df)} rows from {RESULTS_CSV.name}")
    print(f"Labels found ({df['label'].nunique()}): {sorted(df['label'].unique())}")
    return df


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
    """Compute mean/median/p90/p95/p99/min/max/throughput per label, across all 45 Thread Groups."""
    rows = []
    for label, group in df.groupby("label"):
        combo, shape, tier = parse_label(label)
        elapsed, excluded = _exclude_negative_elapsed(group, label)
        total_count = len(group)

        duration_s = (group["timeStamp"] + group["elapsed"].clip(lower=0)).max() - group["timeStamp"].min()
        duration_s = max(duration_s, 1) / 1000
        throughput = total_count / duration_s

        rows.append({
            "label": label, "combo": combo, "shape": shape, "tier": tier,
            "mean": elapsed.mean(), "median": elapsed.median(),
            "p90": elapsed.quantile(0.90), "p95": elapsed.quantile(0.95), "p99": elapsed.quantile(0.99),
            "min": elapsed.min(), "max": elapsed.max(), "throughput": throughput,
            "error_rate": (group["success"] == False).mean() * 100,  # noqa: E712
            "excluded_rows": excluded,
        })
    return pd.DataFrame(rows)


def print_response_time_table(summary: pd.DataFrame, tier: str):
    """Print (and save as CSV) the response time summary table for one size tier."""
    shapes = SHAPE_ORDER_COMPACT if tier == "compact" else SHAPE_ORDER_LARGE
    subset = summary[summary.tier == tier].copy()
    table = subset[["combo", "shape", "mean", "p99", "min", "max"]].copy()
    table["shape"] = table["shape"].map(SHAPE_LABELS)
    order_map = {s: i for i, s in enumerate(shapes)}
    table["_order"] = subset["shape"].map(order_map)
    table = table.sort_values(["combo", "_order"]).drop(columns="_order")
    table.columns = ["Kombinasi", "Variasi Struktur", "Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]
    for col in ["Mean (ms)", "P99 (ms)", "Min (ms)", "Max (ms)"]:
        table[col] = table[col].round(1)

    tier_title = "Compact (~100 KB)" if tier == "compact" else "Large (~500 KB)"
    print(f"\nRingkasan Response Time - Tingkat {tier_title}:")
    print(table.to_string(index=False))

    out_path = SCRIPT_DIR / f"response_time_summary_{tier}.csv"
    table.to_csv(out_path, index=False)
    print(f"Saved {out_path}")


def plot_response_time_overlay(df: pd.DataFrame, shape: str, tier: str):
    """Overlay chart for one (shape, tier) combination: one line per protocol/format combination, x-axis = time since that Thread Group's own execution started (500ms buckets), y-axis = mean response time. Same bucketing/style as the Student dataset's overlay charts."""
    fig, ax = plt.subplots(figsize=(10, 6))

    for prefix, combo in COMBO_NAMES.items():
        label = build_label(prefix, shape, tier)
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
    tier_title = "Compact (~100 KB)" if tier == "compact" else "Large (~500 KB)"
    fig.suptitle(f"Response Time - {SHAPE_LABELS[shape]} ({tier_title})", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    ax.legend(title="Kombinasi")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    out_dir = OUTPUT_COMPACT_DIR if tier == "compact" else OUTPUT_LARGE_DIR
    out_dir.mkdir(parents=True, exist_ok=True)
    filename = f"response_time_overlay_{shape}.png"
    fig.savefig(out_dir / filename, dpi=150)
    print(f"Saved {out_dir / filename}")


def plot_throughput_bar(summary: pd.DataFrame, tier: str):
    """Bar chart: x-axis = protocol/format combination, grouped bars = one per structural shape variant within this tier (4 bars for compact, 5 for large)."""
    shapes = SHAPE_ORDER_COMPACT if tier == "compact" else SHAPE_ORDER_LARGE
    combos = list(COMBO_NAMES.values())
    x = range(len(combos))
    n = len(shapes)
    width = 0.8 / n
    offsets = {shape: (i - (n - 1) / 2) * width for i, shape in enumerate(shapes)}

    fig, ax = plt.subplots(figsize=(13, 6))
    for shape in shapes:
        vals = [summary[(summary.combo == c) & (summary["shape"] == shape) & (summary.tier == tier)]["throughput"].sum()
                for c in combos]
        bars = ax.bar([xi + offsets[shape] for xi in x], vals, width=width,
                      label=SHAPE_LABELS[shape], color=SHAPE_BAR_COLORS[shape], alpha=0.88,
                      edgecolor="#E5E5E5", linewidth=0.6)
        ax.bar_label(bars, fmt="%.0f", fontsize=7, padding=2, rotation=90)

    ax.set_xticks(list(x))
    ax.set_xticklabels(combos, rotation=15, ha="right")
    ax.set_ylabel("Throughput (request/detik)")
    tier_title = "Compact (~100 KB)" if tier == "compact" else "Large (~500 KB)"
    fig.suptitle(f"Throughput - Tingkat {tier_title}", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    ax.legend(title="Variasi Struktur", fontsize=8)
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    out_dir = OUTPUT_COMPACT_DIR if tier == "compact" else OUTPUT_LARGE_DIR
    out_dir.mkdir(parents=True, exist_ok=True)
    fig.savefig(out_dir / "throughput.png", dpi=150)
    print(f"Saved {out_dir / 'throughput.png'}")


def _gradation_ratio(summary: pd.DataFrame, metric: str, tier: str) -> tuple[list[float], list[float]]:
    """Shared helper for both gradation charts: for a given tier, return (x positions, ratio values) of gRPC's metric divided by REST HTTP/2+Protobuf's metric, across every shape defined for that tier."""
    shapes = SHAPE_ORDER_COMPACT if tier == "compact" else SHAPE_ORDER_LARGE
    xs, ys = [], []
    for shape in shapes:
        grpc_row = summary[(summary.combo == "gRPC") & (summary["shape"] == shape) & (summary.tier == tier)]
        rest_row = summary[(summary.combo == "REST HTTP/2 + Protobuf") & (summary["shape"] == shape) & (summary.tier == tier)]
        if grpc_row.empty or rest_row.empty:
            continue
        ratio = grpc_row[metric].iloc[0] / rest_row[metric].iloc[0]
        xs.append(GRADATION_X_POS[shape])
        ys.append(ratio)
    return xs, ys


def _plot_gradation(summary: pd.DataFrame, metric: str, ylabel: str, title: str, filename: str):
    """Shared plotting logic for both gradation charts (response time, throughput). Draws two lines (compact, large) on the same axes, with students-large plotted as a separate marker past a dashed break at x=4."""
    fig, ax = plt.subplots(figsize=(11, 6))
    tier_styles = {"compact": ("#7089A3", "o"), "large": ("#3D5872", "s")}

    for tier, (color, marker) in tier_styles.items():
        xs, ys = _gradation_ratio(summary, metric, tier)
        main_xy = [(x, y) for x, y in zip(xs, ys) if x <= 3]
        ax.plot([x for x, _ in main_xy], [y for _, y in main_xy], label=f"Tingkat {tier.capitalize()}",
                color=color, marker=marker, markersize=6, linewidth=1.8)
        student_xy = [(x, y) for x, y in zip(xs, ys) if x > 3]
        if student_xy:
            ax.plot([x for x, _ in student_xy], [y for _, y in student_xy], color=color,
                    marker="D", markersize=8, linestyle="none")

    ax.axvline(x=4, color="gray", linestyle="--", linewidth=1, alpha=0.6)
    ax.set_xticks([0, 1, 2, 3, 5])
    ax.set_xticklabels([SHAPE_LABELS[s] for s in SHAPE_ORDER_LARGE], rotation=15, ha="right")
    ax.set_xlabel("Variasi Struktur (diurutkan berdasarkan kedalaman bersarang)")
    ax.set_ylabel(ylabel)
    fig.suptitle(title, fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_LABEL, fontsize=9, color="black", style="italic", ha="center")
    ax.legend(title="Tingkat Ukuran")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)
    fig.tight_layout()

    OUTPUT_GRADATION_DIR.mkdir(parents=True, exist_ok=True)
    fig.savefig(OUTPUT_GRADATION_DIR / filename, dpi=150)
    print(f"Saved {OUTPUT_GRADATION_DIR / filename}")


def plot_gradation_response_time(summary: pd.DataFrame):
    _plot_gradation(summary, metric="mean", ylabel="Rasio Mean Response Time (gRPC ÷ REST HTTP/2+Protobuf)",
                    title="Gradasi Kedalaman Struktur - Response Time", filename="response_time_gradation.png")


def plot_gradation_throughput(summary: pd.DataFrame):
    _plot_gradation(summary, metric="throughput", ylabel="Rasio Throughput (gRPC ÷ REST HTTP/2+Protobuf)",
                    title="Gradasi Kedalaman Struktur - Throughput", filename="throughput_gradation.png")


def main():
    df = load_raw_samples()
    summary = compute_summary(df)

    print_response_time_table(summary, "compact")
    print_response_time_table(summary, "large")

    for shape in SHAPE_ORDER_COMPACT:
        plot_response_time_overlay(df, shape, "compact")
    for shape in SHAPE_ORDER_LARGE:
        plot_response_time_overlay(df, shape, "large")

    plot_throughput_bar(summary, "compact")
    plot_throughput_bar(summary, "large")

    plot_gradation_response_time(summary)
    plot_gradation_throughput(summary)


if __name__ == "__main__":
    main()