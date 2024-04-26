package lab

// Producer ...
type Producer interface {
	Produce(Next) []float64
	Compare([]float64, []float64) bool
}

// Next ...
type Next func([]float64) []float64

type undefined any

func multiply(s []int) int {

	r := 0
	if len(s) > 0 {
		r = s[0]
		for i := 1; i < len(s); i++ {
			r *= s[i]
		}
	}

	return r
}

func between(v0, v1, v float64) bool {

	return min(v0, v1) <= v && v <= max(v0, v1)
}

// import (
// 	"math"
// 	"math/rand"
// 	"time"
// )

// type chaos struct {
// 	r *rand.Rand
// }

// func (c *chaos) rand101() float64 {
// 	return float64(c.r.Intn(1000001))/1000000.0*2.0 - 1.0
// }

// func (c *chaos) randUniqueInt(n int, l int) []int {
// 	a := make([]int, n)
// 	for i := range a {
// 		a[i] = i
// 	}
// 	c.r.Shuffle(n, func(i, j int) { a[i], a[j] = a[j], a[i] })
// 	return a[:l]
// }

// func (c *chaos) intn(n int) int {
// 	return c.r.Intn(n)
// }

// func clamp(v, lo, hi float64) float64 {
// 	return math.Min(math.Max(v, lo), hi)
// }

// func clamp01(v float64) float64 {
// 	return clamp(v, 0.0, 1.0)
// }

// func clamp101(v float64) float64 {
// 	return clamp(v, -1.0, 1.0)
// }

// func contains[E comparable](s []E, v E) bool {
// 	return index(s, v) >= 0
// }

// func index[E comparable](s []E, v E) int {
// 	for i, vs := range s {
// 		if v == vs {
// 			return i
// 		}
// 	}
// 	return -1
// }

// type producer interface {
// 	Produce(e *Entity) []float64
// 	Sort(a, b *Entity) bool
// }
