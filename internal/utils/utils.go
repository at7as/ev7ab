package utils

import (
	"math"
	"math/rand"
	"time"
)

// Rad ...
func Rad(a float64) float64 {
	return a * (math.Pi / 180.0)
}

// Deg ...
func Deg(a float64) float64 {
	return a * (180.0 / math.Pi)
}

// RandUniqueInt ...
func RandUniqueInt(n int, l int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a[:l]
}

// ShuffleInts ...
func ShuffleInts(n []int, l int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := make([]int, len(n))
	for i, v := range n {
		a[i] = v
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a[:l]
}

// Clamp ...
func Clamp(v, lo, hi float64) float64 {
	return math.Min(math.Max(v, lo), hi)
}

// Max ...
func Max(s []float64) float64 {
	r := 0.0
	for i, v := range s {
		if i == 0 || v > r {
			r = v
		}
	}
	return r
}

// Norm ...
func Norm(s []float64) []float64 {
	mx := Max(s)
	ns := make([]float64, len(s))
	for i, v := range s {
		ns[i] = v / mx
	}
	return ns
}

// Point ...
type Point struct {
	x float64
	y float64
}

// Clone ...
func (p0 Point) Clone() Point {
	return Point{p0.x, p0.y}
}

// Distance ...
func (p0 Point) Distance(p1 Point) float64 {
	return math.Sqrt(math.Pow(p1.x-p0.x, 2) + math.Pow(p1.y-p0.y, 2))
}

// Rotate ...
func (p0 *Point) Rotate(a float64, p1 Point) *Point {
	s := math.Sin(a)
	c := math.Cos(a)
	p0.x -= p1.x
	p0.y -= p1.y
	x := p0.x*c - p0.y*s
	y := p0.x*s - p0.y*c
	p0.x = x + p1.x
	p0.y = y + p1.y
	return p0
}

// Lerp ...
func (p0 *Point) Lerp(p1 Point, l float64) *Point {
	p0.x = p0.x + (p1.x-p0.x)*l
	p0.y = p0.y + (p1.y-p0.y)*l
	return p0
}

// Towards ...
func (p0 *Point) Towards(p1 Point, l float64) *Point {
	d := l / p0.Distance(p1)
	p0.x = p0.x + (p1.x-p0.x)*d
	p0.y = p0.y + (p1.y-p0.y)*d
	return p0
}

// Middle ...
func (p0 *Point) Middle(p1 Point) *Point {
	return p0.Lerp(p1, 0.5)
}

// Angle2P ...
func (p0 Point) Angle2P(p1 Point, p2 Point) float64 {
	x1 := p0.x - p1.x
	x2 := p2.x - p1.x
	y1 := p0.y - p1.y
	y2 := p2.y - p1.y
	return math.Acos((x1*x2 + y1*y2) / (math.Sqrt(x1*x1+y1*y1) * math.Sqrt(x2*x2+y2*y2)))
}
