package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ReportInfo struct {
	RunID       string
	GeneratedAt string
}

type BenchmarkReport struct {
	RunID       string            `json:"runId"`
	GeneratedAt string            `json:"generatedAt"`
	Metadata    RuntimeMetadata   `json:"metadata"`
	Results     []BenchmarkResult `json:"results"`
}

func NewReportInfo() ReportInfo {
	now := time.Now().UTC()

	return ReportInfo{
		RunID:       now.Format("20060102-150405"),
		GeneratedAt: now.Format(time.RFC3339),
	}
}

func SaveJSONReport(
	info ReportInfo,
	metadata RuntimeMetadata,
	results []BenchmarkResult,
) (string, error) {
	report := BenchmarkReport{
		RunID:       info.RunID,
		GeneratedAt: info.GeneratedAt,
		Metadata:    metadata,
		Results:     results,
	}

	if err := os.MkdirAll("results", 0755); err != nil {
		return "", err
	}

	path := filepath.Join(
		"results",
		fmt.Sprintf("benchmark-%s.json", info.RunID),
	)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	return path, nil
}

func SaveCSVReports(
	info ReportInfo,
	results []BenchmarkResult,
) (string, string, error) {
	if err := os.MkdirAll("results", 0755); err != nil {
		return "", "", err
	}

	summaryPath := filepath.Join(
		"results",
		fmt.Sprintf(
			"benchmark-%s-summary.csv",
			info.RunID,
		),
	)

	samplesPath := filepath.Join(
		"results",
		fmt.Sprintf(
			"benchmark-%s-samples.csv",
			info.RunID,
		),
	)

	if err := saveSummaryCSV(
		summaryPath,
		info,
		results,
	); err != nil {
		return "", "", err
	}

	if err := saveSamplesCSV(
		samplesPath,
		info,
		results,
	); err != nil {
		return "", "", err
	}

	return summaryPath, samplesPath, nil
}

func saveSummaryCSV(
	path string,
	info ReportInfo,
	results []BenchmarkResult,
) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"run_id",
		"generated_at",
		"target",
		"operation",
		"input_size_bytes",
		"approach",
		"batch_size",
		"warmup_runs",
		"execution_order",
		"sample_count",
		"mean_ns",
		"median_ns",
		"standard_deviation_ns",
		"min_ns",
		"max_ns",
		"cov_percent",
		"ci_95_lower_ns",
		"ci_95_upper_ns",
		"outlier_count",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		if err := writeSummaryRow(
			writer,
			info,
			result,
			"baseline",
			result.BaselineStatistics,
			len(result.BaselineSamples),
		); err != nil {
			return err
		}

		if err := writeSummaryRow(
			writer,
			info,
			result,
			"framework",
			result.FrameworkStatistics,
			len(result.FrameworkSamples),
		); err != nil {
			return err
		}
	}

	writer.Flush()

	return writer.Error()
}

func writeSummaryRow(
	writer *csv.Writer,
	info ReportInfo,
	result BenchmarkResult,
	approach string,
	stats Statistics,
	sampleCount int,
) error {
	batchSize := 1
	warmupRuns := 0

	if approach == "framework" {
		batchSize = result.BatchSize
		warmupRuns = result.WarmupRuns
	}

	row := []string{
		info.RunID,
		info.GeneratedAt,
		result.Target,
		result.Operation,
		strconv.Itoa(result.InputSize),
		approach,
		strconv.Itoa(batchSize),
		strconv.Itoa(warmupRuns),
		result.ExecutionOrder,
		strconv.Itoa(sampleCount),
		strconv.FormatFloat(stats.MeanNS, 'f', -1, 64),
		strconv.FormatFloat(stats.MedianNS, 'f', -1, 64),
		strconv.FormatFloat(stats.StdDevNS, 'f', -1, 64),
		strconv.FormatFloat(stats.MinNS, 'f', -1, 64),
		strconv.FormatFloat(stats.MaxNS, 'f', -1, 64),
		strconv.FormatFloat(stats.CoVPercent, 'f', -1, 64),
		strconv.FormatFloat(
			stats.Confidence95LowerNS,
			'f',
			-1,
			64,
		),
		strconv.FormatFloat(
			stats.Confidence95UpperNS,
			'f',
			-1,
			64,
		),
		strconv.Itoa(stats.OutlierCount),
	}

	return writer.Write(row)
}

func saveSamplesCSV(
	path string,
	info ReportInfo,
	results []BenchmarkResult,
) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"run_id",
		"generated_at",
		"target",
		"operation",
		"input_size_bytes",
		"approach",
		"batch_size",
		"execution_order",
		"warmup_runs",
		"sample_number",
		"duration_ns",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		for i, sample := range result.BaselineSamples {
			row := []string{
				info.RunID,
				info.GeneratedAt,
				result.Target,
				result.Operation,
				strconv.Itoa(result.InputSize),
				"baseline",
				"1",
				"0",
				result.ExecutionOrder,
				strconv.Itoa(i + 1),
				strconv.FormatFloat(
					sample.DurationNS,
					'f',
					-1,
					64,
				),
			}

			if err := writer.Write(row); err != nil {
				return err
			}
		}

		for i, sample := range result.FrameworkSamples {
			row := []string{
				info.RunID,
				info.GeneratedAt,
				result.Target,
				result.Operation,
				strconv.Itoa(result.InputSize),
				"framework",
				strconv.Itoa(result.BatchSize),
				strconv.Itoa(result.WarmupRuns),
				result.ExecutionOrder,
				strconv.Itoa(i + 1),
				strconv.FormatFloat(
					sample.DurationNS,
					'f',
					-1,
					64,
				),
			}

			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	writer.Flush()

	return writer.Error()
}
