package lib

import (
	"math"

	"github.com/at7as/ev7ab/pkg/lab"
)

type Example1 struct{}

func (p *Example1) Produce(n lab.Next) []float64 {
	r := n([]float64{10, 20})
	d := math.Abs((r[0] / 30.0) - 1.0)
	return []float64{d, r[0]}
}

func (p *Example1) Compare(a, b []float64) bool {
	return a[0] < b[0]
}

func (p *Example1) Validate(v []float64) bool {

	if v[0] == 0 || v[1] == 0 {
		return false
	}

	return true
}
