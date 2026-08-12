"""
multiplexing_visualize.py

Reads h2load per-request log files (REST HTTP/1.1 and HTTP/2, both
JSON and Protobuf) and ghz JSON result files (gRPC) across five
concurrent-request levels (10, 50, 100, 500, 1000), and produces:
  1. Five response-time CDF charts (one per level), showing what
     fraction of requests completed within a given response time,
     one curve per protocol combination.
  2. One throughput bar chart, x-axis = protocol combination,
     grouped bars = one per concurrent-request level.
  3. One response-time summary table (mean/median/P90/P95/P99/min/max),
     saved as CSV.

REST combinations are read from h2load's --log-file output: tab-
separated, no header, columns are (start time in microseconds,
HTTP status code, duration in microseconds).

The gRPC combination is read from ghz's --format json output: a
JSON object whose "details" key holds a list of per-request records
with timestamp, latency (nanoseconds), status, and error fields.

Expects files named rest-h1-json-m{10,50,100,500,1000}.log etc. and
grpc-m{10,50,100,500,1000}.json in ../../results/multiplexing/.

Run with: python3 multiplexing_visualize.py
"""

import json
import sys
from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

SCRIPT_DIR = Path(__file__).parent
RESULTS_DIR = SCRIPT_DIR.parent.parent / "results" / "multiplexing"
OUTPUT_DIR = SCRIPT_DIR

sys.path.insert(0, str(SCRIPT_DIR.parent))
from style import COMBO_COLORS, COMBO_ORDER, COMBO_FILE_PREFIX, place_legend_outside  # noqa: E402

# Connections are fixed at 1 for every point in this scenario, per the
# agreed configuration; only m (concurrent requests allowed per
# connection) varies. Depth is fixed at Level 0, element count at 100
# (Hundred), matching the same isolation reasoning used in the
# concurrency scenario.
POINTS = [10, 50, 100, 500, 1000]

MULTIPLEX_COLORS = {
    "10 Permintaan": "#D8C4E8",
    "50 Permintaan": "#B79FD1",
    "100 Permintaan": "#9678BA",
    "500 Permintaan": "#74519F",
    "1000 Permintaan": "#4F2D7A",
}

SCENARIO_CONTEXT = r"$\bf{Kedalaman\ Data}$: Level 0 - $\bf{Jumlah\ Elemen}$: 100 - $\bf{Koneksi}$: 1 - $\bf{Total\ Permintaan}$: 5000 (tetap di seluruh titik pengujian)"

def load_rest_durations_ms(combo: str, m: int) -> pd.Series:
    """Load one h2load --log-file result and return its per-request durations in ms.
    - The file has no header row; columns are start time (us), HTTP status, duration (us).
    - REST HTTP/1.1 combinations use h2load's pipelining sense of -m rather than genuine
      multiplexing, included deliberately as a contrast against HTTP/2 and gRPC.
    """
    prefix = COMBO_FILE_PREFIX[combo]
    path = RESULTS_DIR / f"{prefix}-m{m}.log"
    df = pd.read_csv(path, sep="\t", header=None, names=["start_us", "status", "duration_us"])
    return df["duration_us"] / 1000


def load_grpc_durations_ms(m: int) -> pd.Series:
    """Load one ghz --format json result and return its per-request durations in ms.
    - latency is in nanoseconds in ghz's JSON output, unlike h2load's microseconds.
    """
    path = RESULTS_DIR / f"grpc-m{m}.json"
    with open(path) as f:
        data = json.load(f)
    latencies_ns = [d["latency"] for d in data["details"]]
    return pd.Series(latencies_ns) / 1_000_000


def load_durations_ms(combo: str, m: int) -> pd.Series:
    """Dispatch to the correct reader depending on which tool produced this combination's data."""
    if combo == "gRPC":
        return load_grpc_durations_ms(m)
    return load_rest_durations_ms(combo, m)


def plot_response_time_cdf(m: int):
    """CDF chart for one concurrent-request level: one curve per protocol
    combination, showing what fraction of requests completed within a
    given response time. X-axis is zoomed to the 99th percentile
    across all combinations, not the raw max, so rare outliers don't
    squash the part of the curve where most data actually sits.
    - Each curve gets a small dot at its P50 point, matched by color to
      a boxed summary in the bottom-right corner - that corner is
      reliably empty for a CDF shape, since cumulative percentage is
      already near 100% by the time x is large.
    - The summary box's width is measured from the actual rendered text
      width of its widest line, not a hardcoded guess, so it always fits
      snugly regardless of how long the combination names or values are.
    - Combination names and their values are drawn as two separate
      columns (tab-stop style), so every colon lines up regardless of
      how long each combination's name is.
    """
    fig, ax = plt.subplots(figsize=(11, 6))

    all_p99 = []
    p50_data = []
    for combo in COMBO_ORDER:
        vals = load_durations_ms(combo, m).sort_values().values
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

    fig.suptitle(f"Distribusi Response Time — Tingkat Multiplexing {m} Permintaan", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_CONTEXT, ha="center", fontsize=9, color="#555555", style="italic")
    fig.tight_layout(rect=[0, 0, 1, 0.93])

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    filename = f"response_time_cdf_multiplexing_m{m}.png"
    fig.savefig(OUTPUT_DIR / filename, dpi=150, bbox_inches="tight")
    plt.close(fig)
    print(f"Saved {OUTPUT_DIR / filename}")


def print_response_time_summary_table():
    """Print (and save as CSV) the response-time summary statistics table across all multiplexing levels and combinations."""
    rows = []
    for m in POINTS:
        for combo in COMBO_ORDER:
            vals = load_durations_ms(combo, m)
            rows.append({
                "Titik": f"{m} Permintaan",
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
    print("\nRingkasan Statistik Response Time - Skenario Tingkat Multiplexing:")
    print(table.to_string(index=False))

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    out_path = OUTPUT_DIR / "response_time_summary.csv"
    table.to_csv(out_path, index=False)
    print(f"Saved {out_path}")


def compute_throughput(combo: str, m: int) -> float:
    """Compute throughput as m / mean_latency_seconds.
    - This treats m (max concurrent requests per connection) the same way earlier
      scenarios treated VUS: since connections are fixed at 1 and m requests are kept
      continuously in flight, m/mean_latency_seconds approximates the sustained rate,
      analogous to VUS/mean_latency_seconds for the per-vu-iterations closed model.
    """
    mean_ms = load_durations_ms(combo, m).mean()
    return m / (mean_ms / 1000)


def plot_throughput():
    """Bar chart: x-axis = protocol combination, grouped bars = one per concurrent-request level."""
    fig, ax = plt.subplots(figsize=(13, 6))
    x = range(len(COMBO_ORDER))
    n = len(POINTS)
    width = 0.8 / n
    offsets = {m: (i - (n - 1) / 2) * width for i, m in enumerate(POINTS)}

    for m in POINTS:
        label = f"{m} Permintaan"
        vals = [compute_throughput(combo, m) for combo in COMBO_ORDER]
        positions = [xi + offsets[m] for xi in x]
        bars = ax.bar(positions, vals, width=width, label=label,
                      color=MULTIPLEX_COLORS[label], alpha=0.9, edgecolor="#E5E5E5", linewidth=0.6)
        ax.bar_label(bars, fmt="%.0f", fontsize=6, padding=2)

    ax.set_xticks(list(x))
    ax.set_xticklabels(COMBO_ORDER, rotation=15, ha="right")
    ax.set_ylabel("Throughput (request/detik)")
    place_legend_outside(ax, title="Tingkat Multiplexing")
    ax.spines["top"].set_visible(False)
    ax.spines["right"].set_visible(False)

    fig.suptitle("Throughput — Skenario Tingkat Multiplexing", fontsize=13, fontweight="bold", y=0.98)
    fig.text(0.5, 0.925, SCENARIO_CONTEXT, ha="center", fontsize=9, color="#555555", style="italic")
    fig.tight_layout(rect=[0, 0, 1, 0.93])

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    fig.savefig(OUTPUT_DIR / "throughput_multiplexing.png", dpi=150, bbox_inches="tight")
    plt.close(fig)
    print(f"Saved {OUTPUT_DIR / 'throughput_multiplexing.png'}")


def main():
    for m in POINTS:
        plot_response_time_cdf(m)
    plot_throughput()
    print_response_time_summary_table()


if __name__ == "__main__":
    main()