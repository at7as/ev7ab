package lab

import (
	"slices"
)

// Aggregator ...
type Aggregator func([]float64) float64

func sumAggr(in []float64) float64 {

	return sum(in)
}

func avgAggr(in []float64) float64 {

	return sum(in) / float64(len(in))
}

func minAggr(in []float64) float64 {

	return slices.Min(in)
}

func maxAggr(in []float64) float64 {

	return slices.Max(in)
}
