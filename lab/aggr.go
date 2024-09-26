// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lab

import (
	"slices"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// Aggregator provides function that aggregates values from out of the node.
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

// SetAggregator adds or sets custom aggregation function to functions map.
// To activate custom aggregation function need to set value of code to Config.Aggr and pass Config to Lab.Setup.
// May be used to overwrite builtin functions.
func SetAggregator(code string, aggr Aggregator) {

	aggrMap[code] = aggr

}
