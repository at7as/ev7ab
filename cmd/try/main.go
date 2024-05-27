package main

import (
	"fmt"
	"math/rand"

	"gonum.org/v1/gonum/stat"
)

var r *rand.Rand
var n int = 1
var nn int = 2
var input [][]float64
var p01 *project01 = &project01{}

func init() {
	r = rand.New(rand.NewSource(717))
	input = make([][]float64, n)
	for i := range input {
		input[i] = make([]float64, nn)
		for ii := range nn {
			input[i][ii] = r.Float64()
		}
	}
	p01.origin = &model01{
		pt:   p01,
		link: 28,
		in:   []int{nn, nn, nn, nn, nn, nn, nn},
		out:  [][2]int{{nn, nn}, {nn, nn}, {nn, nn}, {nn, nn}, {nn, nn + nn}, {nn, nn}},
		matrix: []web01{
			{0, 0, nn, 0, 1, 0, true},
			{1, 1, nn, 0, 2, 4, true},
			{1, 2, nn, 0, 3, 8, true},
			{2, 3, nn, 0, 4, 12, true},
			{2, 4, nn, 0, 5, 16, false},
			{3, 4, nn, nn, 5, 20, true},
			{3, 5, nn, 0, 6, 24, true},

			// {0, 0, 2, 0, 1, 0},
			// {0, 1, 2, 0, 2, 4},
			// {1, 2, 2, 0, 3, 8},
			// {2, 2, 2, 2, 3, 12},
		},
	}
	p01.pop = make([]*entity01, n)
	for i := range p01.pop {
		p01.pop[i] = spawn01(p01.origin)
	}
}

// func Benchmark_01spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for ii := range p01.pop {
// 			p01.pop[ii] = spawn01(p01.origin)
// 		}
// 	}
// }

// func Benchmark_01exec(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for ii := range p01.pop {
// 			for _, v := range input {
// 				p01.pop[ii].exec(v)
// 			}
// 		}
// 	}
// }

// func Benchmark_01full(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for ii := range p01.pop {
// 			p01.pop[ii] = spawn01(p01.origin)
// 			for _, v := range input {
// 				p01.pop[ii].exec(v)
// 			}
// 		}
// 	}
// }

type project01 struct {
	origin *model01
	pop    []*entity01
}

type model01 struct {
	pt     *project01
	link   int
	in     []int
	out    [][2]int
	matrix []web01
}

type web01 struct {
	in     int
	out    int
	outlen int
	stin   int
	outin  int
	stlink int
	end    bool
}

func spawn01(o *model01) *entity01 {
	e := &entity01{
		origin: o,
		link:   make([]float64, o.link),
		in:     make([][]float64, len(o.in)),
		out:    make([][][]float64, len(o.in)),
	}
	for i := range o.link {
		e.link[i] = r.Float64()
	}
	for i, v := range o.in {
		e.in[i] = make([]float64, v)
	}
	for i, v := range o.out {
		e.out[i] = make([][]float64, v[0])
		for ii := range e.out[i] {
			e.out[i][ii] = make([]float64, v[1])
		}
	}
	return e
}

type entity01 struct {
	origin *model01
	link   []float64
	in     [][]float64
	out    [][][]float64
}

func (e *entity01) exec(in []float64) {
	outini := 0
	o := 0
	e.in[0] = in
	for _, w := range e.origin.matrix {
		for ini, inv := range e.in[w.in] {
			outini = w.stin + ini
			o = w.stlink + ini*w.outlen
			for outi := range w.outlen {
				e.out[w.out][outi][outini] = qlinear(inv, e.link[o+outi])
			}
		}
		if w.end {
			for i, v := range e.out[w.out] {
				e.in[w.outin][i] = stat.Mean(v, nil)
			}
		}
	}
}

func qlinear(v, m float64) float64 {

	return max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
}

func main() {

	for ii := range p01.pop {
		p01.pop[ii] = spawn01(p01.origin)
		for _, v := range input {
			p01.pop[ii].exec(v)
		}
	}

	e := p01.pop[0]
	fmt.Printf("e.origin: %v\n", e.origin)
	fmt.Printf("e.link: %v\n", e.link)
	for i, v := range e.in {
		fmt.Printf("v[%v]: %v\n", i, v)
	}
	for i, v := range e.out {
		for ii, vv := range v {
			fmt.Printf("vv[%v][%v]: %v\n", i, ii, vv)
		}
	}

}
