package lib

import (
	"fmt"
	"math"

	"github.com/at7as/ev7ab/pkg/lab"
)

type ExampleDigits struct{}

func (p *ExampleDigits) Load(setup map[string]string) error {

	return nil
}

func (p *ExampleDigits) Setup(key, value string) error {

	return nil
}

func (p *ExampleDigits) Produce(n lab.Next, op lab.Next) []float64 {
	r := n([]float64{0.1, 0.2})
	d := math.Abs((r[0] / 0.3) - 1.0)
	return []float64{d, r[0]}
}

func (p *ExampleDigits) Compare(a, b []float64) bool {
	return a[0] < b[0]
}

func (p *ExampleDigits) Validate(r []float64) bool {

	if r[0] > 1.0 {
		return false
	}

	return true
}

func (p *ExampleDigits) Best(v []float64) string {

	best := ""
	if len(v) > 0 {
		best = fmt.Sprintf("%.2f", v[0])
	}

	return best
}

func (p *ExampleDigits) Goal(v []float64) bool {

	if v[0] < 0.01 {
		return true
	}

	return false
}
