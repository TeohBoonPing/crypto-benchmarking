package main

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat/distuv"
)

type Statistics struct {
	MeanNS              float64
	MedianNS            float64
	StdDevNS            float64
	MinNS               float64
	MaxNS               float64
	CoVPercent          float64
	Confidence95LowerNS float64
	Confidence95UpperNS float64
	OutlierCount        int
}

func CalculateStatistics(
	samples []BenchmarkSample,
) Statistics {
	if len(samples) == 0 {
		return Statistics{}
	}

	values := make([]float64, len(samples))

	for i, sample := range samples {
		values[i] = sample.DurationNS
	}

	sort.Float64s(values)

	n := len(values)

	var sum float64

	for _, value := range values {
		sum += value
	}

	mean := sum / float64(n)
	medianValue := median(values)

	minValue := values[0]
	maxValue := values[n-1]

	var stdDev float64

	if n > 1 {
		var squaredDifferences float64

		for _, value := range values {
			difference := value - mean
			squaredDifferences += difference * difference
		}

		variance := squaredDifferences / float64(n-1)
		stdDev = math.Sqrt(variance)
	}

	var covPercent float64

	if mean != 0 {
		covPercent = (stdDev / mean) * 100
	}

	confidenceLower := mean
	confidenceUpper := mean

	if n > 1 {
		standardError := stdDev / math.Sqrt(float64(n))

		tDistribution := distuv.StudentsT{
			Mu:    0,
			Sigma: 1,
			Nu:    float64(n - 1),
		}

		tCritical := tDistribution.Quantile(0.975)

		marginOfError := tCritical * standardError

		confidenceLower = mean - marginOfError
		confidenceUpper = mean + marginOfError
	}

	outlierCount := countOutliers(values)

	return Statistics{
		MeanNS:              mean,
		MedianNS:            medianValue,
		StdDevNS:            stdDev,
		MinNS:               minValue,
		MaxNS:               maxValue,
		CoVPercent:          covPercent,
		Confidence95LowerNS: confidenceLower,
		Confidence95UpperNS: confidenceUpper,
		OutlierCount:        outlierCount,
	}
}

func median(values []float64) float64 {
	n := len(values)

	if n == 0 {
		return 0
	}

	middle := n / 2

	if n%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	}

	return values[middle]
}

func countOutliers(values []float64) int {
	if len(values) < 4 {
		return 0
	}

	middle := len(values) / 2

	lowerHalf := values[:middle]

	var upperHalf []float64

	if len(values)%2 == 0 {
		upperHalf = values[middle:]
	} else {
		upperHalf = values[middle+1:]
	}

	q1 := median(lowerHalf)
	q3 := median(upperHalf)

	iqr := q3 - q1

	lowerBound := q1 - (1.5 * iqr)
	upperBound := q3 + (1.5 * iqr)

	count := 0

	for _, value := range values {
		if value < lowerBound || value > upperBound {
			count++
		}
	}

	return count
}
