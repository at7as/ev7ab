package lab

import "math"

type Aggregator interface {
	exec([]float64) float64
}

type sumAggr struct{}

func (aggr *sumAggr) exec(in []float64) float64 {

	r := 0.0
	for _, v := range in {
		r += v
	}

	return r
}

type avgAggr struct{}

func (aggr *avgAggr) exec(in []float64) float64 {

	r := 0.0
	for _, v := range in {
		r += v
	}

	return r / float64(len(in))
}

type minAggr struct{}

func (aggr *minAggr) exec(in []float64) float64 {

	r := 0.0
	for _, v := range in {
		r = math.Min(r, v)
	}

	return r
}

type maxAggr struct{}

func (aggr *maxAggr) exec(in []float64) float64 {

	r := 0.0
	for _, v := range in {
		r = math.Max(r, v)
	}

	return r
}
