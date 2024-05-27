package main

import (
	"math"

	"github.com/at7as/ev7ab/internal/lab"
)

type producer struct{}

func (p *producer) Produce(n lab.Next) []float64 {
	r := n([]float64{10, 20})
	d := math.Abs((r[0] / 30.0) - 1.0)
	return []float64{d, r[0]}
}

func (p *producer) Compare(a, b []float64) bool {
	return a[0] < b[0]
}

func (p *producer) Validate(v []float64) bool {

	if v[0] == 0 || v[1] == 0 {
		return false
	}

	return true
}

func main() {

	cfg := lab.Config{
		In:     2,
		Out:    1,
		Target: []float64{0.0},
		Limit:  []float64{0.01},
		Goal:   true,
		Size:   1000,
	}
	l := lab.New(cfg, &producer{})
	l.Examine()

	// l := lab.NewLab(1000, 0.05, 0.2, 0.3, 0.2, 0.5, &producer{}, 1, true)
	// l.Add(2, "")
	// l.Add(10, "")
	// l.Add(1, "")
	// l.Examine(20)

}
