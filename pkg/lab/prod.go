package lab

import "slices"

// Producer ...
type Producer interface {
	Load(map[string]string) error
	Setup(string, string) error
	Produce(Next, Next, []float64) []float64
	Validate([]float64) bool
	Compare([]float64, []float64) bool
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

// IndexOfMaxExcluded ...
func IndexOfMaxExcluded(v []float64, excl []int) int {

	index := 0
	if len(v) == 0 {
		return index
	}
	for i := range v {
		if !slices.Contains(excl, i) {
			index = i
			break
		}
	}
	vmax := v[index]
	for i, vv := range v {
		if vv > vmax && !slices.Contains(excl, i) {
			vmax = vv
			index = i
		}
	}

	return index
}
