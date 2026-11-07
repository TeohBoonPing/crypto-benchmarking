package main

type BenchmarkResult struct {
	Target         string `json:"target"`
	Operation      string `json:"operation"`
	InputSize      int    `json:"inputSize"`
	BatchSize      int    `json:"batchSize"`
	WarmupRuns     int    `json:"warmupRuns"`
	ExecutionOrder string `json:"executionOrder"`

	BaselineSamples  []BenchmarkSample `json:"baselineSamples"`
	FrameworkSamples []BenchmarkSample `json:"frameworkSamples"`

	BaselineStatistics  Statistics `json:"baselineStatistics"`
	FrameworkStatistics Statistics `json:"frameworkStatistics"`
}
