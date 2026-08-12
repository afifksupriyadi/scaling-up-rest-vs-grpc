"""
depth_visualize.py

Reads k6 CSV results (all 5 protocol combinations: REST HTTP/1.1
JSON/Protobuf, REST HTTP/2 JSON/Protobuf, gRPC) across three depth
levels (Level 0, Level 2, Level 4), and produces:
  1. Three response-time overlay charts (one per level), x-axis =
     time since test start, y-axis = response time, one line per
     protocol combination.
  2. One throughput bar chart, x-axis = protocol combination,
     grouped bars = one per depth level.
  3. One data_received summary table (average bytes per response),
     saved as CSV.

Expects files named rest-h1-json-depth-{zero,two,four}.csv etc. and
grpc-depth-{zero,two,four}.csv in ../../results/depth/.

Run with: python3 depth_visualize.py
"""

import sys
from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

SCRIPT_DIR = Path(__file__).parent
RESULTS_DIR = SCRIPT_DIR.parent.parent / "results" / "depth"
OUTPUT_DIR = SCRIPT_DIR

sys.path.insert(0, str(SCRIPT_DIR.parent))
from style import COMBO_COLORS, COMBO_ORDER, COMBO_FILE_PREFIX, place_legend_outside  # noqa: E402

# VUs are fixed at 10 for every point in this scenario, per the
# agreed configuration; only depth (structural shape) varies.
VUS = 10

LEVELS = [
    ("zero", "Level 0"),
    ("two", "Level 2"),
    ("four", "Level 4"),
]

DEPTH_COLORS = {
    "Level 0": "#8FB8DE",
    "Level 2": "#4A82B4",
    "Level 4": "#1F5788",
}

SCENARIO_CONTEXT = r"$\bf{Jumlah\ Elemen}$: 1 - $\bf{Koneksi}$: 10 - $\bf{Total\ Permintaan}$: 500 (tetap di seluruh titik pengujian)"

def metric_name_for(combo: str) -> str:
    """Return the k6 metric_name that holds this combination's response time (gRPC and REST log it under different names)."""
    return "grpc_req_duration" if combo == "gRPC" else "http_req_duration"


def load_csv(prefix: str, level_suffix: str) -> pd.DataFrame:
    """Load one raw k6 CSV result file for this scenario."""
    path = RESULTS_DIR / f"{prefix}-depth-{level_suffix}.csv"
    return pd.read_csv(path, low_memory=False)


def plot_response_time_cdf(level_suffix: str, level_label: str):
    """CDF chart for one depth level: one curve per protocol combination,
    showing what fraction of requests completed within a given response
    time. X-axis is zoomed to the 99th percentile across all
    combinations, not the raw max, so rare outliers don't squash the
    part of the curve where most data actually sits.
    - Each curve gets a small dot at its P50 point, matched by color to
      a boxed summary in the bottom-right corner - that corner is
      reliably empty for a CDF shape, since cumulative percentage is
      already near 100% by the time x is large.
    - The summary box's width is measured from the actual rendered text
      width of its widest line, not a hardcoded guess, so it always fits
      snugly regardless of how long the combination names or values are.
    """
    fig, ax = plt.subplots(figsize=(11, 6))

    all_p99 = []
    p50_data = []
    for combo in COMBO_ORDER:
        prefix = COMBO_FILE_PREFIX[combo]
        df = load_csv(prefix, level_suffix)
        vals = df[df["metric_name"] == metric_name_for(combo)]["metric_value"].sort_values().values
        cumulative_pct = (pd.Series(range(1, len(vals) + 1)) / len(vals)) * 100
        ax.plot(vals, cumulative_pct, label=combo, color=COMBO_COLORS[combo], linewidth=1.8)
        all_p99.append(pd.Series(vals).quantile(0.99))
        p50 = pd.Series(vals).quantile(0.50)
        p50_data.append((combo, p50))
        ax.plot(p50, 50, "o", color=COMBO_COLORS[combo], markersize=6,
                markeredgecolor="white", markeredgewidth=0.8, zorder=5)

    ax.set_xlim(0, max(all_p99) * 1.05)
    ax.set_xlabel("Response Time (ms)")
    ax.set_ylabel("Persentase Permintaan Kumulatif (%)")
    ax.grid(True, alpha=0.3)
    place_legend_outside(ax)
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)

    fastest_combo = min(p50_data, key=lambda cp: cp[1])[0]
    title_text = "Median Response Time per Kombinasi (P50)"
    name_texts = [f"• {combo}" for combo, p50 in p50_data]
    value_texts = []
    for combo, p50 in p50_data:
        is_fastest = combo == fastest_combo
        value_texts.append(f": {p50:.2f} ms" + ("*" if is_fastest else ""))
    footnote_text = "*: tercepat (P50 terendah)"

    pad_x = 0.014
    pad_top = 0.022
    title_gap = 0.020
    line_step = 0.030
    footnote_gap = 0.020
    tab_gap = 0.006
    box_right = 0.98
    box_y_top_border = 0.34
    SAFETY_MARGIN = 1.15

    fig.tight_layout()
    fig.canvas.draw()
    renderer = fig.canvas.get_renderer()
    ax_width_px = ax.get_window_extent(renderer=renderer).width

    def text_width_fraction(s, fontsize, weight="normal", style="normal"):
        t = ax.text(0, 0, s, fontsize=fontsize, fontweight=weight, style=style, alpha=0)
        fig.canvas.draw()
        bbox = t.get_window_extent(renderer=renderer)
        t.remove()
        return bbox.width / ax_width_px

    name_widths = [text_width_fraction(name_texts[i], 8, weight="bold" if p50_data[i][0] == fastest_combo else "normal") for i in range(len(p50_data))]
    value_widths = [text_width_fraction(value_texts[i], 8, weight="bold" if p50_data[i][0] == fastest_combo else "normal") for i in range(len(p50_data))]
    max_name_w = max(name_widths)
    max_value_w = max(value_widths)
    title_w = text_width_fraction(title_text, 7.5, weight="bold")
    footnote_w = text_width_fraction(footnote_text, 6.5, style="italic")

    entry_full_width = max_name_w + tab_gap + max_value_w
    box_width = max(title_w, entry_full_width, footnote_w) * SAFETY_MARGIN + pad_x * 2
    box_left = box_right - box_width
    text_x = box_left + pad_x
    value_col_x = text_x + max_name_w + tab_gap

    n_entries = len(p50_data)
    box_height = pad_top + 0.026 + title_gap + line_step * n_entries + footnote_gap + 0.018 + pad_x
    box_bottom = box_y_top_border - box_height

    ax.add_patch(plt.Rectangle(
        (box_left, box_bottom), box_width, box_height,
        transform=ax.transAxes, facecolor="white", edgecolor="#CCCCCC",
        linewidth=0.8, zorder=6, clip_on=False,
    ))

    y = box_y_top_border - pad_top
    ax.text(text_x, y, title_text, transform=ax.transAxes, fontsize=7.5, fontweight="bold",
            color="black", ha="left", va="top", zorder=7)
    y -= title_gap + 0.026

    for i, (combo, p50) in enumerate(p50_data):
        is_fastest = combo == fastest_combo
        color = COMBO_COLORS[combo]
        weight = "bold" if is_fastest else "normal"
        ax.text(text_x, y, name_texts[i], transform=ax.transAxes, fontsize=8, color=color,
                fontweight=weight, ha="left", va="top", zorder=7)
        ax.text(value_col_x, y, value_texts[i], transform=ax.transAxes, fontsize=8, color=color,
                fontweight=weight, ha="left", va="top", zorder=7)
        y -= line_step

    y -= footnote_gap
    ax.text(text_x, y, footnote_text, transform=ax.transAxes, fontsize=6.5, color="black",
            style="italic", ha="left", va="top", zorder=7)

    fig.suptitle(f"Distribusi Response Time — Kedalaman Data {level_label}", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_CONTEXT, ha="center", fontsize=9, color="#555555", style="italic")
    fig.tight_layout(rect=[0, 0, 1, 0.93])
    
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    filename = f"response_time_cdf_level_{level_suffix}.png"
    fig.savefig(OUTPUT_DIR / filename, dpi=150, bbox_inches="tight")
    plt.close(fig)
    print(f"Saved {OUTPUT_DIR / filename}")


def print_response_time_summary_table():
    """Print (and save as CSV) the response-time summary statistics table across all depth levels and combinations."""
    rows = []
    for level_suffix, level_label in LEVELS:
        for combo in COMBO_ORDER:
            prefix = COMBO_FILE_PREFIX[combo]
            df = load_csv(prefix, level_suffix)
            vals = df[df["metric_name"] == metric_name_for(combo)]["metric_value"]
            rows.append({
                "Titik": level_label,
                "Kombinasi": combo,
                "Mean (ms)": round(vals.mean(), 3),
                "Median (ms)": round(vals.median(), 3),
                "P90 (ms)": round(vals.quantile(0.90), 3),
                "P95 (ms)": round(vals.quantile(0.95), 3),
                "P99 (ms)": round(vals.quantile(0.99), 3),
                "Min (ms)": round(vals.min(), 3),
                "Max (ms)": round(vals.max(), 3),
            })
    table = pd.DataFrame(rows)
    print("\nRingkasan Statistik Response Time - Skenario Kedalaman Data:")
    print(table.to_string(index=False))

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    out_path = OUTPUT_DIR / "response_time_summary.csv"
    table.to_csv(out_path, index=False)
    print(f"Saved {out_path}")


def compute_throughput(prefix: str, level_suffix: str, metric_name: str) -> float:
    """Compute throughput as VUS / mean_latency_seconds, not from raw timestamp span.
    - These tests often complete in well under a second (e.g. Level 0's ~0.8ms mean
      latency x 50 iterations is ~40ms total), far too short for k6's whole-second
      timestamp resolution to measure duration accurately.
    - This assumes the per-vu-iterations closed-model executor used throughout this
      test suite: each VU runs its iterations strictly sequentially, so one VU's own
      rate is 1/mean_latency_seconds, and VUS identical parallel VUs multiply that.
    """
    df = load_csv(prefix, level_suffix)
    mean_ms = df[df["metric_name"] == metric_name]["metric_value"].mean()
    return VUS / (mean_ms / 1000)


def plot_throughput():
    """Bar chart: x-axis = protocol combination, grouped bars = one per depth level."""
    fig, ax = plt.subplots(figsize=(12, 6))
    x = range(len(COMBO_ORDER))
    n = len(LEVELS)
    width = 0.8 / n
    offsets = {label: (i - (n - 1) / 2) * width for i, (_, label) in enumerate(LEVELS)}

    for level_suffix, level_label in LEVELS:
        vals = [compute_throughput(COMBO_FILE_PREFIX[combo], level_suffix, metric_name_for(combo)) for combo in COMBO_ORDER]
        positions = [xi + offsets[level_label] for xi in x]
        bars = ax.bar(positions, vals, width=width, label=level_label,
                      color=DEPTH_COLORS[level_label], alpha=0.9, edgecolor="#E5E5E5", linewidth=0.6)
        ax.bar_label(bars, fmt="%.0f", fontsize=6, padding=2)

    ax.set_xticks(list(x))
    ax.set_xticklabels(COMBO_ORDER, rotation=15, ha="right")
    ax.set_ylabel("Throughput (request/detik)")
    place_legend_outside(ax, title="Kedalaman")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)

    fig.suptitle("Throughput — Skenario Kedalaman Data", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_CONTEXT, ha="center", fontsize=9, color="#555555", style="italic")
    fig.tight_layout(rect=[0, 0, 1, 0.93])

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    fig.savefig(OUTPUT_DIR / "throughput_depth.png", dpi=150, bbox_inches="tight")
    plt.close(fig)
    print(f"Saved {OUTPUT_DIR / 'throughput_depth.png'}")


def compute_data_received(prefix: str, level_suffix: str) -> float:
    """Compute the average data_received (bytes) per request for one combination/level."""
    df = load_csv(prefix, level_suffix)
    return df[df["metric_name"] == "data_received"]["metric_value"].mean()


def print_data_received_table():
    """Print (and save as CSV) the data_received summary table across all depth levels and combinations."""
    rows = []
    for level_suffix, level_label in LEVELS:
        row = {"Titik": level_label}
        for combo in COMBO_ORDER:
            row[combo] = round(compute_data_received(COMBO_FILE_PREFIX[combo], level_suffix), 1)
        rows.append(row)

    table = pd.DataFrame(rows)
    print("\nRingkasan data_received (byte) - Skenario Kedalaman Data:")
    print(table.to_string(index=False))

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    out_path = OUTPUT_DIR / "data_received_summary.csv"
    table.to_csv(out_path, index=False)
    print(f"Saved {out_path}")


def main():
    for level_suffix, level_label in LEVELS:
        plot_response_time_cdf(level_suffix, level_label)
    plot_throughput()
    print_response_time_summary_table()
    print_data_received_table()


if __name__ == "__main__":
    main()