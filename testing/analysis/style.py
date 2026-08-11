"""
style.py

Shared plotting constants and helpers used by all three scenario
visualization scripts (depth, element-count, concurrency), so the
five protocol combinations and their legend placement stay identical
across every chart instead of drifting between files.
"""

# Colors for the five protocol combinations, used consistently across
# all 12 charts (9 response-time overlays + 3 throughput bar charts).
COMBO_COLORS = {
    "REST HTTP/1.1 + JSON": "#A9764F",
    "REST HTTP/1.1 + Protobuf": "#7B6FA6",
    "REST HTTP/2 + JSON": "#C74B5C",
    "REST HTTP/2 + Protobuf": "#C9A227",
    "gRPC": "#4F9D8C",
}

# Fixed display order for the five combinations (legend entries,
# throughput x-axis ticks), so ordering doesn't depend on dict/file
# iteration order and drift between the three scenario scripts.
COMBO_ORDER = [
    "REST HTTP/1.1 + JSON",
    "REST HTTP/1.1 + Protobuf",
    "REST HTTP/2 + JSON",
    "REST HTTP/2 + Protobuf",
    "gRPC",
]

# Filename prefix per protocol combination, identical across all three
# scenarios since every run-*.sh script writes results using this same
# naming convention (only the suffix after this prefix differs).
COMBO_FILE_PREFIX = {
    "REST HTTP/1.1 + JSON": "rest-h1-json",
    "REST HTTP/1.1 + Protobuf": "rest-h1-protobuf",
    "REST HTTP/2 + JSON": "rest-h2-json",
    "REST HTTP/2 + Protobuf": "rest-h2-protobuf",
    "gRPC": "grpc",
}


def place_legend_outside(ax, title="Kombinasi"):
    """Place the legend outside the plot area, to the right, so it never overlaps the data regardless of how busy the data is.
    - The caller must leave enough right-side margin for it, e.g. by saving with bbox_inches='tight', otherwise the legend gets cut off at the image edge.
    """
    ax.legend(title=title, bbox_to_anchor=(1.02, 1), loc="upper left")