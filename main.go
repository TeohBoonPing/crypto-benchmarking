package main

import (
	"fmt"
	"log"
)

type BenchmarkCase struct {
	NewTarget func() BenchmarkTarget
	BatchSize int
}

func main() {
	metadata := GetRuntimeMetadata()

	fmt.Println("Runtime metadata")
	fmt.Printf("Go version: %s\n", metadata.GoVersion)
	fmt.Printf("OS: %s\n", metadata.OS)
	fmt.Printf("Architecture: %s\n", metadata.Architecture)
	fmt.Printf("CPU count: %d\n", metadata.CPUCount)
	fmt.Printf("GOMAXPROCS: %d\n", metadata.GOMAXPROCS)
	fmt.Println()

	rsaData, err := NewRSAData()
	if err != nil {
		log.Fatal(err)
	}

	cases := []BenchmarkCase{
		// SHA-256
		{
			NewTarget: func() BenchmarkTarget {
				return NewSHA256Target(1024)
			},
			BatchSize: 1,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewSHA256Target(4 * 1024)
			},
			BatchSize: 1,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewSHA256Target(16 * 1024)
			},
			BatchSize: 1000,
		},

		// AES-GCM
		{
			NewTarget: func() BenchmarkTarget {
				return NewAESGCMTarget(1024)
			},
			BatchSize: 1,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewAESGCMTarget(4 * 1024)
			},
			BatchSize: 1,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewAESGCMTarget(16 * 1024)
			},
			BatchSize: 10,
		},

		// ChaCha20-Poly1305
		{
			NewTarget: func() BenchmarkTarget {
				return NewChaCha20Poly1305Target(1024)
			},
			BatchSize: 1,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewChaCha20Poly1305Target(4 * 1024)
			},
			BatchSize: 10,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewChaCha20Poly1305Target(16 * 1024)
			},
			BatchSize: 1,
		},

		// RSA
		{
			NewTarget: func() BenchmarkTarget {
				return NewRSASignTarget(rsaData)
			},
			BatchSize: 1,
		},
		{
			NewTarget: func() BenchmarkTarget {
				return NewRSAVerifyTarget(rsaData)
			},
			BatchSize: 2,
		},
	}

	var results []BenchmarkResult

	for i, benchmarkCase := range cases {
		frameworkFirst := i%2 == 1

		result, err := runBenchmarkCase(
			benchmarkCase,
			frameworkFirst,
		)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, *result)

		fmt.Println()
	}

	fmt.Printf("Completed %d benchmark cases.\n", len(results))

	reportInfo := NewReportInfo()

	jsonPath, err := SaveJSONReport(
		reportInfo,
		metadata,
		results,
	)
	if err != nil {
		log.Fatal(err)
	}

	summaryPath, samplesPath, err := SaveCSVReports(
		reportInfo,
		results,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Run ID: %s\n", reportInfo.RunID)
	fmt.Printf("JSON report saved: %s\n", jsonPath)
	fmt.Printf("Summary CSV saved: %s\n", summaryPath)
	fmt.Printf("Samples CSV saved: %s\n", samplesPath)
}

func runBenchmarkCase(
	benchmarkCase BenchmarkCase,
	frameworkFirst bool,
) (*BenchmarkResult, error) {
	config := DefaultBenchmarkConfig()
	config.BatchSize = benchmarkCase.BatchSize

	var baselineSamples []BenchmarkSample
	var frameworkSamples []BenchmarkSample

	var baselineTarget BenchmarkTarget
	var frameworkTarget BenchmarkTarget

	var err error

	if frameworkFirst {
		frameworkTarget = benchmarkCase.NewTarget()

		if err := frameworkTarget.Setup(); err != nil {
			return nil, err
		}

		frameworkSamples, err = RunFramework(
			frameworkTarget,
			config,
		)
		if err != nil {
			return nil, err
		}

		baselineTarget = benchmarkCase.NewTarget()

		if err := baselineTarget.Setup(); err != nil {
			return nil, err
		}

		baselineSamples, err = RunBaseline(
			baselineTarget,
			config.Samples,
		)
		if err != nil {
			return nil, err
		}
	} else {
		baselineTarget = benchmarkCase.NewTarget()

		if err := baselineTarget.Setup(); err != nil {
			return nil, err
		}

		baselineSamples, err = RunBaseline(
			baselineTarget,
			config.Samples,
		)
		if err != nil {
			return nil, err
		}

		frameworkTarget = benchmarkCase.NewTarget()

		if err := frameworkTarget.Setup(); err != nil {
			return nil, err
		}

		frameworkSamples, err = RunFramework(
			frameworkTarget,
			config,
		)
		if err != nil {
			return nil, err
		}
	}

	baselineStats := CalculateStatistics(
		baselineSamples,
	)

	frameworkStats := CalculateStatistics(
		frameworkSamples,
	)

	executionOrder := "baseline-first"

	if frameworkFirst {
		executionOrder = "framework-first"
	}

	result := &BenchmarkResult{
		Target:              baselineTarget.Name(),
		Operation:           baselineTarget.Operation(),
		InputSize:           baselineTarget.InputSize(),
		BatchSize:           config.BatchSize,
		WarmupRuns:          config.WarmupRuns,
		ExecutionOrder:      executionOrder,
		BaselineSamples:     baselineSamples,
		FrameworkSamples:    frameworkSamples,
		BaselineStatistics:  baselineStats,
		FrameworkStatistics: frameworkStats,
	}

	fmt.Println("========================================")
	fmt.Printf("Benchmark target: %s\n", result.Target)
	fmt.Printf("Operation: %s\n", result.Operation)
	fmt.Printf("Input size: %d bytes\n", result.InputSize)
	fmt.Printf("Batch size: %d\n", result.BatchSize)
	fmt.Printf("Warm-up runs: %d\n", result.WarmupRuns)

	if frameworkFirst {
		fmt.Println("Execution order: Framework -> Baseline")
	} else {
		fmt.Println("Execution order: Baseline -> Framework")
	}

	fmt.Println()

	printStatistics(
		"Baseline",
		result.BaselineStatistics,
	)

	fmt.Println()

	printStatistics(
		"Framework",
		result.FrameworkStatistics,
	)

	return result, nil
}

func printStatistics(
	name string,
	stats Statistics,
) {
	fmt.Println(name)

	fmt.Printf(
		"Mean: %.2f ns\n",
		stats.MeanNS,
	)

	fmt.Printf(
		"Median: %.2f ns\n",
		stats.MedianNS,
	)

	fmt.Printf(
		"Standard deviation: %.2f ns\n",
		stats.StdDevNS,
	)

	fmt.Printf(
		"Minimum: %.2f ns\n",
		stats.MinNS,
	)

	fmt.Printf(
		"Maximum: %.2f ns\n",
		stats.MaxNS,
	)

	fmt.Printf(
		"CoV: %.2f%%\n",
		stats.CoVPercent,
	)

	fmt.Printf(
		"95%% CI: %.2f ns to %.2f ns\n",
		stats.Confidence95LowerNS,
		stats.Confidence95UpperNS,
	)

	fmt.Printf(
		"Outliers: %d\n",
		stats.OutlierCount,
	)
}
