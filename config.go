package main

type BenchmarkConfig struct {
	WarmupRuns int
	Samples    int
	BatchSize  int
}

func DefaultBenchmarkConfig() BenchmarkConfig {
	return BenchmarkConfig{
		WarmupRuns: 10,
		Samples:    30,
		BatchSize:  1,
	}
}
