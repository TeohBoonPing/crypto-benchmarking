# Go-Based Performance Evaluation Framework for Cryptographic Algorithms

This repository contains the implementation and experimental materials for the research project:

**A Go-Based Performance Evaluation Framework for Cryptographic Algorithms**

The project was developed as part of the **Master of Software Engineering (Software Technology)** programme at **Universiti Malaya**.

## Overview

This project provides a reusable Go-based software benchmarking framework for controlled execution-time benchmarking of cryptographic algorithms.

The framework provides a structured benchmarking process by integrating benchmark configuration, controlled execution, warm-up runs, repeated measurement, configurable batch execution, statistical analysis, runtime metadata collection, and structured result reporting.

The main focus of the study is **execution-time measurement variability**, rather than determining which cryptographic algorithm is the fastest.

## Benchmark Targets

The framework was evaluated using five benchmark targets:

- AES-GCM encryption
- ChaCha20-Poly1305 encryption
- SHA-256 hashing
- RSA signing
- RSA verification

AES-GCM, ChaCha20-Poly1305 and SHA-256 were evaluated using:

- 1 KB
- 4 KB
- 16 KB

RSA signing and verification were evaluated using a fixed SHA-256 digest input.

## Framework Features

The framework provides:

- Separation of benchmark setup from measured execution
- Configurable warm-up runs
- Repeated execution-time measurement
- Configurable batch sizes
- Pilot-based batch-size calibration
- Statistical analysis of timing measurements
- Mean and median execution time
- Standard deviation
- Coefficient of Variation (CoV)
- 95% confidence intervals
- Outlier detection using Tukey's 1.5 × IQR rule
- Runtime and platform metadata collection
- JSON output
- Summary CSV output
- Sample-level CSV output
- Reproducible CoV comparison chart generation

## Experimental Evaluation

The proposed framework was compared with a baseline direct-measurement approach.

The formal evaluation consisted of:

- 5 benchmark targets
- 11 benchmark configurations
- 30 measured samples per approach for each configuration
- 660 formal timing samples in total
- 10 warm-up runs for the proposed framework
- Pilot-selected batch sizes
- Alternated baseline-first and framework-first execution order

Both approaches were evaluated using the same cryptographic implementation, input configuration, hardware, operating system, Go version, and runtime configuration.

The **Coefficient of Variation (CoV)** was used as the primary measure of execution-time measurement variability.

## Experimental Environment

The formal experiment was performed using:

- **Go:** 1.27.0
- **Operating System:** macOS (darwin)
- **Architecture:** ARM64
- **CPU Count:** 8
- **GOMAXPROCS:** 8

Benchmark results may differ when the framework is executed on different hardware, operating systems, or runtime configurations.

## Results

Across the 11 evaluated benchmark configurations:

- The framework produced lower CoV in **8 out of 11 configurations**
- The baseline produced lower CoV in **3 out of 11 configurations**
- The framework achieved CoV below 10% in **9 out of 11 configurations**
- The baseline achieved CoV below 10% in **6 out of 11 configurations**
- Framework median CoV: **3.67%**
- Baseline median CoV: **6.90%**

Overall, the proposed framework generally produced lower execution-time measurement variability under the evaluated conditions.

However, the framework did not produce lower variability in every benchmark configuration. The results should therefore be interpreted together with the benchmark workload and execution environment.

## Requirements

### Benchmark Framework

- Go 1.27.0 or compatible version
- Make

### Chart Generation

The chart-generation script requires Python 3 with:

- pandas
- numpy
- matplotlib

The required Python packages can be installed using:

```bash
pip3 install pandas numpy matplotlib
```

## Run

Run the benchmarking framework:

```bash
make run
```

This executes the cryptographic benchmarks and generates the benchmark result files.

## Generate Chart

Generate the overall CoV comparison chart:

```bash
make chart
```

The chart is generated from the latest benchmark summary CSV produced by the framework.

The generated chart is saved as:

```text
charts/overall-cov-comparison.png
```

The chart compares the Coefficient of Variation (CoV) between the baseline and proposed framework across the 11 benchmark configurations.

It also includes the **10% CoV stability criterion** used in the study.

## Output

Each benchmark execution generates:

- JSON report
- Summary CSV
- Sample-level CSV

### JSON Report

Contains:

- Benchmark configuration
- Execution order
- Individual timing samples
- Statistical results
- Runtime and platform metadata

### Summary CSV

Contains the aggregated statistical results for each benchmark configuration and measurement approach, including:

- Mean
- Median
- Standard deviation
- Minimum and maximum
- Coefficient of Variation (CoV)
- 95% confidence interval
- Outlier count

The summary CSV is also used to generate the overall CoV comparison chart.

### Sample-Level CSV

Contains the individual execution-time measurements used to calculate the statistical results.

## Chart Reproduction

The overall CoV comparison figure used in the research can be reproduced directly from the benchmark output.

First, run the benchmark:

```bash
make run
```

Then generate the chart:

```bash
make chart
```

The process provides a direct path from:

```text
Benchmark execution
        ↓
JSON / CSV results
        ↓
Summary CSV
        ↓
CoV comparison chart
```

This allows the experimental results and the corresponding chart to be reproduced from the source code and generated benchmark data.

## Interpretation of Results

The framework should not be interpreted as eliminating all execution-time variability.

Execution-time measurements may still be affected by factors such as:

- Garbage collection
- Memory allocation
- Operating-system scheduling
- Cache behaviour
- Processor-frequency changes
- Background system activity
- Benchmark workload
- Execution order

The framework is intended to provide a more structured and controlled benchmarking process while recording statistical and runtime information to support interpretation of the results.

## Limitations

The experimental results were obtained using:

- One machine
- One operating system
- One Go/runtime configuration
- One controlled execution session

The evaluation was also limited to the selected benchmark targets and input sizes.

The results should therefore not be generalized automatically to all cryptographic algorithms, hardware platforms, runtime environments, or operating systems.

## Research Use

This repository accompanies the dissertation:

**A Go-Based Performance Evaluation Framework for Cryptographic Algorithms**

The source code, benchmark outputs, and chart-generation scripts are provided to support transparency and reproducibility of the benchmarking procedure and experimental analysis.

## Author

Teoh Boon Ping