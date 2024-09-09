package lab

import (
	"slices"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// Aggregator ...
type Aggregator func([]float64) float64

func sumAggr(in []float64) float64 {

	return floats.Sum(in)
}

func avgAggr(in []float64) float64 {

	return stat.Mean(in, nil)
}

func minAggr(in []float64) float64 {

	return slices.Min(in)
}

func maxAggr(in []float64) float64 {

	return slices.Max(in)
}

var aggrMap map[string]Aggregator = map[string]Aggregator{
	"sum": sumAggr,
	"avg": avgAggr,
	"min": minAggr,
	"max": maxAggr,
}

// SetAggregator ...
func SetAggregator(code string, aggr Aggregator) {

	aggrMap[code] = aggr

}
