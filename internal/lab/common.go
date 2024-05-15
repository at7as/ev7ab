package lab

// Producer ...
type Producer interface {
	Produce(Next) []float64
	Compare([]float64, []float64) bool
}

// Next ...
type Next func([]float64) []float64

type undefined any

func sum[V int | float64](v []V) V {

	var r V
	for _, vv := range v {
		r += vv
	}

	return r
}

func multiply[V int | float64](s []V) V {

	var r V
	if len(s) > 0 {
		r = s[0]
		for i := 1; i < len(s); i++ {
			r *= s[i]
		}
	}

	return r
}

func between[V int | float64](v0, v1, v V) bool {

	return min(v0, v1) <= v && v <= max(v0, v1)
}
