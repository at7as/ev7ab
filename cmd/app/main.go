package main

import (
	"math"

	"github.com/at7as/ev7ab/internal/lab"
)

type producer struct{}

func (p *producer) Produce(n lab.Next) []float64 {
	r := n([]float64{10, 20})
	d := math.Abs((r[0] / 30.0) - 1.0)
	return []float64{r[0], d}
}

func (p *producer) Compare(a, b []float64) bool {
	return a[1] < b[1]
}

func main() {

	l := lab.New(&producer{})
	l.Configure([]int{2}, []int{1}, []string{}, []string{}, []float64{math.NaN(), 0.0}, []float64{0.0, 0.01}, true)
	l.Examine()

	// l := lab.NewLab(1000, 0.05, 0.2, 0.3, 0.2, 0.5, &producer{}, 1, true)
	// l.Add(2, "")
	// l.Add(10, "")
	// l.Add(1, "")
	// l.Examine(20)

}
