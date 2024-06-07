package main

import (
	"math/rand/v2"
	"sync"
	"testing"
)

var r []*rand.Rand
var inputRandom *rand.Rand
var ni int = 1
var ne int = 10
var nn int = 100
var input [][]float64
var p01 *project01 = &project01{}
var lPool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, 7*nn*nn)
		return &b
	},
}
var inPool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn*3)
		return &b
	},
}
var outPool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn*3)
		return &b
	},
}
var in6Pool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn*3)
		return &b
	},
}
var out6Pool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn*3)
		return &b
	},
}
var dataPool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn*3)
		return &b
	},
}

func init() {
	r = make([]*rand.Rand, ne)
	for i := range r {
		r[i] = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
	inputRandom = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	input = make([][]float64, ni)
	for i := range input {
		input[i] = make([]float64, nn)
		for ii := range nn {
			input[i][ii] = inputRandom.Float64()
		}
	}
	p01.origin = &model01{
		pt:   p01,
		link: nn * nn * 7,
		in:   []int{nn, nn, nn, nn, nn, nn, nn},
		out:  [][2]int{{nn, nn}, {nn, nn}, {nn, nn}, {nn, nn}, {nn, nn + nn}, {nn, nn}},
		matrix2: []web02{
			{[][2]int{{0, nn}}, float64(nn), nn, true},
			{[][2]int{{0, nn}}, float64(nn), nn, false},
			{[][2]int{{0, nn}}, float64(nn), nn, true},
			{[][2]int{{0, nn}}, float64(nn), nn, false},
			{[][2]int{{0, nn}, {nn, nn * 2}}, float64(nn) * 2.0, nn, false},
			{[][2]int{{nn, nn * 2}}, float64(nn), nn, true},
		},
	}
	p01.pop = make([]*entity01, ne)

	for range ne {
		buf := lPool.Get().(*[]float64)
		lPool.Put(buf)
	}

}

func Benchmark_10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for ii := range p01.pop {
			wg.Add(1)
			go spawn08(p01.origin, r[ii], wg, ii)
		}
		wg.Wait()
		for ii := range p01.pop {
			wg.Add(1)
			go execGroup07(ii, wg)
		}
		wg.Wait()
	}
}

func execGroup07(ii int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range input {
		p01.pop[ii].exec8(v)
	}
	p01.pop[ii].link = p01.pop[ii].link[:0]
	lPool.Put(p01.pop[ii].bufPtr)
	p01.pop[ii].in = p01.pop[ii].in[:0]
	dataPool.Put(p01.pop[ii].inPtr)
	p01.pop[ii].out = p01.pop[ii].out[:0]
	dataPool.Put(p01.pop[ii].outPtr)

}

type project01 struct {
	origin *model01
	pop    []*entity01
}

type model01 struct {
	pt      *project01
	link    int
	in      []int
	out     [][2]int
	matrix2 []web02
}

type web02 struct {
	in     [][2]int
	inlen  float64
	outlen int
	end    bool
}

func spawn07(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	bufPtr := lPool.Get().(*[]float64)
	e := entity01{
		origin: o,
		link:   *bufPtr,
		bufPtr: bufPtr,
	}
	for range o.link {
		e.link = append(e.link, rr.Float64())
	}
	p01.pop[iii] = &e
}

func spawn08(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	bufPtr := lPool.Get().(*[]float64)
	indataPointer := dataPool.Get().(*[]float64)
	outdataPointer := dataPool.Get().(*[]float64)
	e := entity01{
		origin: o,
		link:   *bufPtr,
		bufPtr: bufPtr,
		in:     *indataPointer,
		out:    *outdataPointer,
		inPtr:  indataPointer,
		outPtr: outdataPointer,
	}
	for range o.link {
		e.link = append(e.link, rr.Float64())
	}
	p01.pop[iii] = &e
}

type entity01 struct {
	origin *model01
	link   []float64
	in     []float64
	out    []float64
	bufPtr *[]float64
	inPtr  *[]float64
	outPtr *[]float64
}

func (e *entity01) exec8(in []float64) {
	e.in = append(e.in, in...)
	l := 0
	for _, w := range e.origin.matrix2 {
		for range w.outlen {
			v := 0.0
			for _, ini := range w.in {
				for _, inv := range e.in[ini[0]:ini[1]] {
					v += qlinear(inv, e.link[l])
					l++
				}
			}
			e.out = append(e.out, v/w.inlen)
		}
		if w.end {
			e.in = e.in[:0]
			e.in = append(e.in, e.out...)
			e.out = e.out[:0]
		}
	}
}

func qlinear(v, m float64) float64 {

	return max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
}
