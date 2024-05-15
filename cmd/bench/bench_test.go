package main

import (
	"math/big"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var valuesFloat []float64 = make([]float64, 100001)
var float big.Float

func init() {
	for i := 0; i < len(valuesFloat); i++ {
		valuesFloat[i] = float64(i) / 10000.0
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

func Benchmark_floatBig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumBig(valuesFloat)
	}
}

func sum(v []float64) float64 {

	var r float64
	for i := range v {
		r += v[i]
	}

	return r
}

func sumBig(v []float64) float64 {

	r := big.NewFloat(0)
	f := big.NewFloat(0)
	for i := range v {
		f.SetFloat64(v[i])
		r.Add(r, f)
		// r.SetFloat64(float.Add(&r, vv))
	}
	res, _ := r.Float64()
	return res
}
