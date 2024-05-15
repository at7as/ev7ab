package main

import (
	"testing"

	"gonum.org/v1/gonum/floats"
)

var valuesFloat []float64 = make([]float64, 65536)
var valuesUint16 []uint16 = make([]uint16, 65536)

func init() {
	for i := 0; i < len(valuesFloat); i++ {
		valuesFloat[i] = float64(i) / 65535.0
	}
	for i := 0; i < len(valuesUint16); i++ {
		valuesUint16[i] = uint16(i)
	}
}

func Benchmark_float(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum(valuesFloat)
	}
}

func Benchmark_floatGonum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		floats.Sum(valuesFloat)
	}
}

func Benchmark_uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := valuesUint16[0]
		for ii := 1; ii < len(valuesUint16); ii++ {
			v = v & valuesUint16[ii]
		}
	}
}

func sum(v []float64) float64 {

	var r float64
	for i := range v {
		r += v[i]
	}

	return r
}
