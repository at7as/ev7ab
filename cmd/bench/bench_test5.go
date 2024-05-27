package main

import (
	"math"
	"math/rand"
	"testing"
)

var mod []float64 = []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.99}
var result float64 = 0.0

func Benchmark_Linear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = max(1-math.Abs(v-m)/0.5, 0.0)
		}
	}
}

func Benchmark_Quadratic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = max(-4.0*v*v+8.0*m*v-4.0*m*m+1.0, 0.0)
		}
	}
}

func Benchmark_Qlinear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
		}
	}
}

func Benchmark_Linear_Max(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = Max(1-math.Abs(v-m)/0.5, 0.0)
		}
	}
}

func Benchmark_Quadratic_Max(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = Max(-4.0*v*v+8.0*m*v-4.0*m*m+1.0, 0.0)
		}
	}
}

func Benchmark_Qlinear_Max(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = Max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
		}
	}
}

func Benchmark_Linear_mathMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = math.Max(1-math.Abs(v-m)/0.5, 0.0)
		}
	}
}

func Benchmark_Quadratic_mathMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = math.Max(-4.0*v*v+8.0*m*v-4.0*m*m+1.0, 0.0)
		}
	}
}

func Benchmark_Qlinear_mathMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		for _, m := range mod {
			result = math.Max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
		}
	}
}

func Max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}
