package main

import "time"

type BenchmarkTarget interface {
	Name() string
	Operation() string
	InputSize() int
	Setup() error
	Run() error
}

type BenchmarkSample struct {
	DurationNS float64
}

func RunBaseline(target BenchmarkTarget, samples int) ([]BenchmarkSample, error) {
	results := make([]BenchmarkSample, 0, samples)

	for i := 0; i < samples; i++ {
		start := time.Now()

		if err := target.Run(); err != nil {
			return nil, err
		}

		elapsed := time.Since(start)

		results = append(results, BenchmarkSample{
			DurationNS: float64(elapsed.Nanoseconds()),
		})
	}
	return results, nil
}

func RunFramework(target BenchmarkTarget, config BenchmarkConfig) ([]BenchmarkSample, error) {
	// Warm-up runs are not recorded.
	for i := 0; i < config.WarmupRuns; i++ {
		if err := target.Run(); err != nil {
			return nil, err
		}
	}

	results := make([]BenchmarkSample, 0, config.Samples)

	for i := 0; i < config.Samples; i++ {
		start := time.Now()

		for j := 0; j < config.BatchSize; j++ {
			if err := target.Run(); err != nil {
				return nil, err
			}
		}

		elapsed := time.Since(start)

		durationPerOperation := float64(elapsed.Nanoseconds()) / float64(config.BatchSize)

		results = append(results, BenchmarkSample{
			DurationNS: durationPerOperation,
		})
	}

	return results, nil
}
