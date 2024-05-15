package main

import (
	"math"
	"testing"
)

var values100001 []float32 = make([]float32, 100001)

func init() {
	for i := 0; i < len(values100001); i++ {
		values100001[i] = float32(i) / 10.0
	}
}

func Benchmark_linearProc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range values100001 {
			linearProc(v)
		}
	}
}

func Benchmark_binaryProc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range values100001 {
			binaryProc(v)
		}
	}
}

func Benchmark_reluProc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range values100001 {
			reluProc(v)
		}
	}
}

func Benchmark_lreluProc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range values100001 {
			lreluProc(v)
		}
	}
}

func Benchmark_softsignProc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range values100001 {
			softsignProc(v)
		}
	}
}

// (-*,*)
func linearProc(in float32) float32 {

	return in
}

// {0,1}
func binaryProc(in float32) float32 {

	if in < 0 {
		return 0.0
	}

	return 1.0
}

// [0,*)
func reluProc(in float32) float32 {

	if in < 0 {
		return 0.0
	}

	return in
}

// (-*,*)
func lreluProc(in float32) float32 {

	if in < 0 {
		return in * 0.01
	}

	return in
}

// (-1,1)
func softsignProc(in float32) float32 {

	return in / (1.0 + float32(math.Abs(float64(in))))
}
