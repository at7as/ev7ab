package lab

// Producer ...
type Producer interface {
	Produce(Next) []float64
	Compare([]float64, []float64) bool
	Validate([]float64) bool
}

// Next ...
type Next func([]float64) []float64

type undefined any

func sum(v []int) int {

	var s int
	for _, vv := range v {
		s += vv
	}

	return s
}

func qlinear(v, m float64) float64 {

	return max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
}

// func multiply[V int | float64](s []V) V {

// 	var r V
// 	if len(s) > 0 {
// 		r = s[0]
// 		for i := 1; i < len(s); i++ {
// 			r *= s[i]
// 		}
// 	}

// 	return r
// }

func between[V int | float64](v0, v1, v V) bool {

	return min(v0, v1) <= v && v <= max(v0, v1)
}
