package lab

// Producer ...
type Producer interface {
	Load(map[string]string) error
	Setup(string, string) error
	Produce(Next, Next) []float64
	Compare([]float64, []float64) bool
	Validate([]float64) bool
	Best([]float64) string
	Goal([]float64) bool
}

// Next ...
type Next func([]float64) []float64

// IndexOfMax ...
func IndexOfMax(v []float64) int {

	index := 0
	if len(v) == 0 {
		return index
	}
	vmax := v[0]
	for i, vv := range v {
		if vv > vmax {
			vmax = vv
			index = i
		}
	}

	return index
}
