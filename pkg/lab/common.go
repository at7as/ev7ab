package lab

// Producer ...
type Producer interface {
	Produce(Next) []float64
	Compare([]float64, []float64) bool
	Validate([]float64) bool
}

// Next ...
type Next func([]float64) []float64
