package main

// use pool in exec

import (
	"math/rand/v2"
	"sync"
	"testing"
)

var r []*rand.Rand
var inputRandom *rand.Rand
var ni int = 10
var ne int = 100
var nn int = 10
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
		b := make([]float64, 0, nn)
		return &b
	},
}
var outPool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn)
		return &b
	},
}
var inArrayPool = &sync.Pool{
	New: func() any {
		b := make([]array, 0, 7)
		return &b
	},
}
var outArrayPool = &sync.Pool{
	New: func() any {
		b := make([]array, 0, 7)
		return &b
	},
}
var in6Pool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn)
		return &b
	},
}
var out6Pool = &sync.Pool{
	New: func() any {
		b := make([]float64, 0, nn)
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
		matrix2: []web02{
			{[][2]int{{0, nn}}, float64(nn), nn, true},
			{[][2]int{{0, nn}}, float64(nn), nn, false},
			{[][2]int{{0, nn}}, float64(nn), nn, true},
			{[][2]int{{0, nn}}, float64(nn), nn, false},
			{[][2]int{{0, nn}, {nn, nn * 2}}, float64(nn) * 2.0, nn, false},
			{[][2]int{{nn, nn * 2}}, float64(nn), nn, true},
		},
	}
	for i := range p01.origin.matrix {
		p01.origin.matrix[i].stlink = i * nn * nn
	}
	p01.origin.link = nn * nn * len(p01.origin.matrix)
	p01.pop = make([]*entity01, ne)

}

func Benchmark_06(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for ii := range p01.pop {
			wg.Add(1)
			go spawn05(p01.origin, r[ii], wg, ii)
		}
		wg.Wait()
		for ii := range p01.pop {
			wg.Add(1)
			go execGroup03(ii, wg)
		}
		wg.Wait()
	}
}

func Benchmark_07(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for ii := range p01.pop {
			wg.Add(1)
			go spawn06(p01.origin, r[ii], wg, ii)
		}
		wg.Wait()
		for ii := range p01.pop {
			wg.Add(1)
			go execGroup04(ii, wg)
		}
		wg.Wait()
	}
}

func Benchmark_08(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for ii := range p01.pop {
			wg.Add(1)
			go spawn07(p01.origin, r[ii], wg, ii)
		}
		wg.Wait()
		for ii := range p01.pop {
			wg.Add(1)
			go execGroup05(ii, wg)
		}
		wg.Wait()
	}
}

func execGroup03(ii int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range input {
		p01.pop[ii].exec4(v)
	}
	p01.pop[ii].link = p01.pop[ii].link[:0]
	lPool.Put(p01.pop[ii].bufPtr)
}

func execGroup04(ii int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range input {
		p01.pop[ii].exec5(v)
	}
	p01.pop[ii].link = p01.pop[ii].link[:0]
	lPool.Put(p01.pop[ii].bufPtr)
	for i := range p01.pop[ii].indata {
		p01.pop[ii].indata[i].data = p01.pop[ii].indata[i].data[:0]
		inPool.Put(p01.pop[ii].indata[i].buf)
	}
	p01.pop[ii].indata = p01.pop[ii].indata[:0]
	inArrayPool.Put(p01.pop[ii].indPtr)
	for i := range p01.pop[ii].outdata {
		p01.pop[ii].outdata[i].data = p01.pop[ii].outdata[i].data[:0]
		outPool.Put(p01.pop[ii].outdata[i].buf)
	}
	p01.pop[ii].outdata = p01.pop[ii].outdata[:0]
	outArrayPool.Put(p01.pop[ii].outdPtr)

}

func execGroup05(ii int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range input {
		p01.pop[ii].exec6(v)
	}
	p01.pop[ii].link = p01.pop[ii].link[:0]
	lPool.Put(p01.pop[ii].bufPtr)
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
	matrix  []web01
	matrix2 []web02
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

type web02 struct {
	in     [][2]int
	inlen  float64
	outlen int
	end    bool
}

func spawn05(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	bufPtr := lPool.Get().(*[]float64)
	e := entity01{
		origin: o,
		link:   *bufPtr,
		bufPtr: bufPtr,
		in:     make([][]float64, len(o.in)),
		out:    make([][]float64, len(o.out)),
	}
	for range o.link {
		e.link = append(e.link, rr.Float64())
	}
	for i, v := range o.in {
		e.in[i] = make([]float64, v)
	}
	for i, v := range o.out {
		e.out[i] = make([]float64, v[0])
	}
	p01.pop[iii] = &e
}

func spawn06(o *model01, rr *rand.Rand, wg *sync.WaitGroup, iii int) {
	defer wg.Done()
	bufPtr := lPool.Get().(*[]float64)
	indPtr := inArrayPool.Get().(*[]array)
	outdPtr := outArrayPool.Get().(*[]array)
	e := entity01{
		origin:  o,
		link:    *bufPtr,
		bufPtr:  bufPtr,
		indata:  *indPtr,
		outdata: *outdPtr,
		indPtr:  indPtr,
		outdPtr: outdPtr,
	}
	for range o.link {
		e.link = append(e.link, rr.Float64())
	}
	for i, v := range o.in {
		inPtr := inPool.Get().(*[]float64)
		e.indata = append(e.indata, array{data: *inPtr, buf: inPtr})
		for range v {
			e.indata[i].data = append(e.indata[i].data, 0.0)
		}
	}
	for i, v := range o.out {
		outPtr := outPool.Get().(*[]float64)
		e.outdata = append(e.outdata, array{data: *outPtr, buf: outPtr})
		for range v[0] {
			e.outdata[i].data = append(e.outdata[i].data, 0.0)
		}
	}
	p01.pop[iii] = &e
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

type array struct {
	data []float64
	buf  *[]float64
}

type entity01 struct {
	origin  *model01
	link    []float64
	in      [][]float64
	out     [][]float64
	indata  []array
	outdata []array
	bufPtr  *[]float64
	indPtr  *[]array
	outdPtr *[]array
}

func (e *entity01) exec4(in []float64) {
	o := 0
	e.in[0] = in
	for _, w := range e.origin.matrix {
		for outi := range w.outlen {
			e.out[w.out][outi] = 0
		}
		for ini, inv := range e.in[w.in] {
			o = w.stlink + ini*w.outlen
			for outi := range w.outlen {
				lv := e.link[o+outi]
				e.out[w.out][outi] += qlinear(inv, lv)
			}
		}
		if w.end {
			for i, v := range e.out[w.out] {
				e.in[w.outin][i] = v / float64(e.origin.out[w.out][1])
			}
		}
	}
}

func (e *entity01) exec5(in []float64) {
	o := 0
	e.indata[0].data = in
	for _, w := range e.origin.matrix {
		for outi := range w.outlen {
			e.outdata[w.out].data[outi] = 0
		}
		for ini, inv := range e.indata[w.in].data {
			o = w.stlink + ini*w.outlen
			for outi := range w.outlen {
				lv := e.link[o+outi]
				e.outdata[w.out].data[outi] += qlinear(inv, lv)
			}
		}
		if w.end {
			for i, v := range e.outdata[w.out].data {
				e.indata[w.outin].data[i] = v / float64(e.origin.out[w.out][1])
			}
		}
	}
}

func (e *entity01) exec6(in []float64) {
	indataPointer := in6Pool.Get().(*[]float64)
	indata := *indataPointer
	indata = append(indata, in...)
	outdataPointer := out6Pool.Get().(*[]float64)
	outdata := *outdataPointer
	l := 0
	for _, w := range e.origin.matrix2 {
		for range w.outlen {
			v := 0.0
			for _, ini := range w.in {
				for _, inv := range indata[ini[0]:ini[1]] {
					v += qlinear(inv, e.link[l])
					l++
				}
			}
			outdata = append(outdata, v/w.inlen)
		}
		if w.end {
			indata = indata[:0]
			indata = append(indata, outdata...)
			outdata = outdata[:0]
		}
	}
	indata = indata[:0]
	in6Pool.Put(indataPointer)
	outdata = outdata[:0]
	out6Pool.Put(outdataPointer)
}

func qlinear(v, m float64) float64 {

	return max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
}
