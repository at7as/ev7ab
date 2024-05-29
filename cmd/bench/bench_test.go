package main

import (
	"math/rand"
	rand2 "math/rand/v2"
	"sync"
	"testing"

	"gonum.org/v1/gonum/stat"
)

var r []*rand.Rand
var r2 []*rand2.Rand
var inputRandom *rand.Rand
var inputRandom2 *rand2.Rand
var ni int = 1
var ne int = 100
var nn int = 10
var input [][]float64
var p01 *project01 = &project01{}

func init() {
	r = make([]*rand.Rand, ne)
	for i := range r {
		r[i] = rand.New(rand.NewSource(rand.Int63()))
	}
	r2 = make([]*rand2.Rand, ne)
	for i := range r2 {
		r2[i] = rand2.New(rand2.NewPCG(rand2.Uint64(), rand2.Uint64()))
	}
	inputRandom = rand.New(rand.NewSource(rand.Int63()))
	inputRandom2 = rand2.New(rand2.NewPCG(rand2.Uint64(), rand2.Uint64()))
	input = make([][]float64, ni)
	for i := range input {
		input[i] = make([]float64, nn)
		for ii := range nn {
			input[i][ii] = inputRandom.Float64()
		}
	}
	p01.origin = &model01{
		pt:   p01,
		link: 0,
		in:   []int{nn, nn, nn, nn, nn, nn, nn},
		out:  [][2]int{{nn, nn}, {nn, nn}, {nn, nn}, {nn, nn}, {nn, nn + nn}, {nn, nn}},
		matrix: []web01{
			{0, 0, nn, 0, 1, 0, true},
			{1, 1, nn, 0, 2, 0, true},
			{1, 2, nn, 0, 3, 0, true},
			{2, 3, nn, 0, 4, 0, true},
			{2, 4, nn, 0, 5, 0, false},
			{3, 4, nn, nn, 5, 0, true},
			{3, 5, nn, 0, 6, 0, true},
		},
	}
	for i := range p01.origin.matrix {
		p01.origin.matrix[i].stlink = i * nn * nn
	}
	p01.origin.link += nn * nn
	p01.pop = make([]*entity01, ne)

}

// func Benchmark_04(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn02(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go execGroup(ii, wg)
// 		}
// 		wg.Wait()
// 	}
// }

// func Benchmark_01spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn01(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 	}
// }

// func Benchmark_02spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn02(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 	}
// }

func Benchmark_22spawn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for ii := range p01.pop {
			wg.Add(1)
			go spawn22(p01.origin, r2[ii], wg, ii)
		}
		wg.Wait()
	}
}

// func Benchmark_03spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn03(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 	}
// }

// func Benchmark_04spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn04(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 	}
// }

// func Benchmark_05spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn05(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 	}
// }

// func Benchmark_06spawn(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		wg := &sync.WaitGroup{}
// 		for ii := range p01.pop {
// 			wg.Add(1)
// 			go spawn06(p01.origin, r[ii], wg, ii)
// 		}
// 		wg.Wait()
// 	}
// }

// pool of []
// deep clone of origin with []
// one [] and exec without goroutines

func execGroup(ii int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range input {
		p01.pop[ii].exec3(v)
	}
}

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

func spawn01(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := entity01{
		origin: o,
		// link:   make([]float64, o.link),
		// in:     make([][]float64, len(o.in)),
		// out:    make([][][]float64, len(o.out)),
	}
	// for i := range o.link {
	// 	e.link[i] = rr.Float64()
	// }
	// for i, v := range o.in {
	// 	e.in[i] = make([]float64, v)
	// }
	// for i, v := range o.out {
	// 	e.out[i] = make([][]float64, v[0])
	// 	for ii := range e.out[i] {
	// 		e.out[i][ii] = make([]float64, v[1])
	// 	}
	// }
	p01.pop[iii] = &e
}

func spawn02(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := &entity01{
		origin: o,
		link:   make([]float64, o.link),
		in:     make([][]float64, len(o.in)),
		out:    make([][][]float64, len(o.out)),
	}
	for i := range o.link {
		e.link[i] = rr.Float64()
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
	p01.pop[iii] = e
}

func spawn22(o *model01, rr *rand2.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := &entity01{
		origin: o,
		link:   make([]float64, o.link),
		in:     make([][]float64, len(o.in)),
		out:    make([][][]float64, len(o.out)),
	}
	for i := range o.link {
		e.link[i] = rr.Float64()
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
	p01.pop[iii] = e
}

func spawn03(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := entity01{
		origin: o,
		link:   make([]float64, o.link),
		// in:     make([][]float64, len(o.in)),
		// out:    make([][][]float64, len(o.out)),
	}
	// for i := range o.link {
	// 	e.link[i] = rr.Float64()
	// }
	// for i, v := range o.in {
	// 	e.in[i] = make([]float64, v)
	// }
	// for i, v := range o.out {
	// 	e.out[i] = make([][]float64, v[0])
	// 	for ii := range e.out[i] {
	// 		e.out[i][ii] = make([]float64, v[1])
	// 	}
	// }
	p01.pop[iii] = &e
}

func spawn04(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := entity01{
		origin: o,
		link:   make([]float64, o.link),
		// in:     make([][]float64, len(o.in)),
		// out:    make([][][]float64, len(o.out)),
	}
	for i := range o.link {
		e.link[i] = rr.Float64()
	}
	// for i, v := range o.in {
	// 	e.in[i] = make([]float64, v)
	// }
	// for i, v := range o.out {
	// 	e.out[i] = make([][]float64, v[0])
	// 	for ii := range e.out[i] {
	// 		e.out[i][ii] = make([]float64, v[1])
	// 	}
	// }
	p01.pop[iii] = &e
}

func spawn05(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := entity01{
		origin: o,
		link:   make([]float64, o.link),
		in:     make([][]float64, len(o.in)),
		// out:    make([][][]float64, len(o.out)),
	}
	for i := range o.link {
		e.link[i] = rr.Float64()
	}
	// for i, v := range o.in {
	// 	e.in[i] = make([]float64, v)
	// }
	// for i, v := range o.out {
	// 	e.out[i] = make([][]float64, v[0])
	// 	for ii := range e.out[i] {
	// 		e.out[i][ii] = make([]float64, v[1])
	// 	}
	// }
	p01.pop[iii] = &e
}

func spawn06(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	e := entity01{
		origin: o,
		link:   make([]float64, o.link),
		in:     make([][]float64, len(o.in)),
		// out:    make([][][]float64, len(o.out)),
	}
	for i := range o.link {
		e.link[i] = rr.Float64()
	}
	for i, v := range o.in {
		e.in[i] = make([]float64, v)
	}
	// for i, v := range o.out {
	// 	e.out[i] = make([][]float64, v[0])
	// 	for ii := range e.out[i] {
	// 		e.out[i][ii] = make([]float64, v[1])
	// 	}
	// }
	p01.pop[iii] = &e
}

type entity01 struct {
	origin *model01
	link   []float64
	in     [][]float64
	out    [][][]float64
}

func (e *entity01) exec3(in []float64) {
	outini := 0
	o := 0
	e.in[0] = in
	for _, w := range e.origin.matrix {
		for ini, inv := range e.in[w.in] {
			outini = w.stin + ini
			o = w.stlink + ini*w.outlen
			for outi := range w.outlen {
				l := e.link[o+outi]
				e.out[w.out][outi][outini] = qlinear(inv, l)
				// += then div
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
