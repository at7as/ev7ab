package main

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var big []float64
var medium [][]float64
var small [][][]float64
var result2 float64

func init() {
	iiii := 0
	big = make([]float64, 6*4*40000)
	medium = make([][]float64, 6)
	small = make([][][]float64, 6)
	for i := range 6 {
		iiiii := 0
		medium[i] = make([]float64, 4*40000)
		small[i] = make([][]float64, 4)
		for ii := range 4 {
			small[i][ii] = make([]float64, 40000)
			for iii := range 40000 {
				v := rand.Float64()
				big[iiii] = v
				medium[i][iiiii] = v
				small[i][ii][iii] = v
				iiii++
				iiiii++
			}
		}
	}
}

func Benchmark_small(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		for _, v := range small {
			for _, vv := range v {
				result2 = floats.Sum(vv)
			}
		}
	}
}

func Benchmark_medium(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		for _, v := range medium {
			for i := range 4 {
				result2 = floats.Sum(v[i*40000 : (i+1)*40000])
			}
		}
	}
}

func Benchmark_big(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		for iii := range 24 {
			result2 = floats.Sum(big[iii*40000 : (iii+1)*40000])
		}
	}
}
