import glob
import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd


# Find the latest summary CSV
files = glob.glob("results/*-summary.csv")

if not files:
    raise FileNotFoundError("No summary CSV found in output/")

summary_file = max(files, key=os.path.getmtime)

print(f"Using: {summary_file}")

df = pd.read_csv(summary_file)

# Build benchmark labels
def benchmark_label(row):
    target = row["target"]
    operation = row["operation"]
    size = row["input_size_bytes"]

    if target == "RSA":
        return f"RSA\n{operation.capitalize()}"

    if size == 1024:
        size_label = "1 KB"
    elif size == 4096:
        size_label = "4 KB"
    elif size == 16384:
        size_label = "16 KB"
    else:
        size_label = f"{size} B"

    return f"{target}\n{size_label}"


df["label"] = df.apply(benchmark_label, axis=1)

# Keep benchmark order from the CSV
labels = df["label"].drop_duplicates().tolist()

baseline = (
    df[df["approach"].str.lower() == "baseline"]
    .set_index("label")
    .loc[labels]["cov_percent"]
    .values
)

framework = (
    df[df["approach"].str.lower() == "framework"]
    .set_index("label")
    .loc[labels]["cov_percent"]
    .values
)

x = np.arange(len(labels))
width = 0.38

fig, ax = plt.subplots(figsize=(14, 6))

ax.bar(x - width / 2, baseline, width, label="Baseline")
ax.bar(x + width / 2, framework, width, label="Framework")

# 10% CoV stability criterion
ax.axhline(
    y=10,
    linewidth=1.5,
    label="10% CoV stability criterion",
)

ax.set_yscale("log")

ax.set_title(
    "Overall CoV Comparison Between Baseline and Proposed Framework"
)
ax.set_xlabel("Benchmark Configuration")
ax.set_ylabel("Coefficient of Variation (%) - logarithmic scale")

ax.set_xticks(x)
ax.set_xticklabels(labels, fontsize=8)

ax.grid(axis="y", which="both", alpha=0.2)
ax.legend()

plt.tight_layout()

os.makedirs("charts", exist_ok=True)

output_file = "charts/overall-cov-comparison.png"
plt.savefig(output_file, dpi=300, bbox_inches="tight")

print(f"Chart generated: {output_file}")