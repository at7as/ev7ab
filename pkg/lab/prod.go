package lab

// Producer ...
type Producer interface {
	Load(map[string]string) error
	Produce(Next, Next) []float64
	Compare([]float64, []float64) bool
	Validate([]float64) bool
	Best([]float64) string
	Goal([]float64) bool
}

// Next ...
type Next func([]float64) []float64
