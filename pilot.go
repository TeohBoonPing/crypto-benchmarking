package main

import (
	"fmt"
	"time"
)

func RunBatchSizePilot(
	newTarget func() BenchmarkTarget,
	batchSizes []int,
) error {
	const pilotRuns = 5

	runMeans := make(map[int][]BenchmarkSample)
	stableRunCount := make(map[int]int)

	// Preconditioning round.
	// Results are intentionally discarded.
	fmt.Println("Preconditioning for 2 seconds...")

	target := newTarget()

	if err := target.Setup(); err != nil {
		return err
	}

	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {
		if err := target.Run(); err != nil {
			return err
		}
	}

	fmt.Println("Preconditioning complete.")
	fmt.Println()

	fmt.Println("Run | Batch size | Mean (ns) | CoV")

	for run := 1; run <= pilotRuns; run++ {
		// Rotate the starting batch size for each run.
		for i := 0; i < len(batchSizes); i++ {
			index := (i + run - 1) % len(batchSizes)
			batchSize := batchSizes[index]

			target := newTarget()

			if err := target.Setup(); err != nil {
				return err
			}

			config := DefaultBenchmarkConfig()
			config.BatchSize = batchSize

			samples, err := RunFramework(target, config)
			if err != nil {
				return err
			}

			stats := CalculateStatistics(samples)

			fmt.Printf(
				"%3d | %10d | %9.2f | %.2f%%\n",
				run,
				batchSize,
				stats.MeanNS,
				stats.CoVPercent,
			)

			runMeans[batchSize] = append(
				runMeans[batchSize],
				BenchmarkSample{
					DurationNS: stats.MeanNS,
				},
			)

			if stats.CoVPercent < 10 {
				stableRunCount[batchSize]++
			}
		}

		fmt.Println()
	}

	fmt.Println("Summary")
	fmt.Println(
		"Batch size | Mean of runs (ns) | Run-mean CoV | Runs <10% | Stable",
	)

	for _, batchSize := range batchSizes {
		stats := CalculateStatistics(runMeans[batchSize])

		stable := "YES"

		if stats.CoVPercent >= 10 || stableRunCount[batchSize] < 4 {
			stable = "NO"
		}

		fmt.Printf(
			"%10d | %17.2f | %12.2f%% | %9d/5 | %s\n",
			batchSize,
			stats.MeanNS,
			stats.CoVPercent,
			stableRunCount[batchSize],
			stable,
		)
	}

	return nil
}
