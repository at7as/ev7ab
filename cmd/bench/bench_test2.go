package main

import (
	"math/rand"
	"testing"
)

var arr map[uint8]float32 = make(map[uint8]float32)

func init() {
	var i uint8 = 0
	for i <= 255 {
		arr[i] = calc(i)
		if i == 255 {
			break
		}
		i++
	}
}

// 1 from index
// 2 from calc
// 3 from float

func calc(v uint8) float32 {
	return float32(v) / 255.0
}

func BenchmarkProduceIndex(b *testing.B) {
	// var arr map[uint8]float32 = make(map[uint8]float32)
	// var i uint8 = 0
	// for i <= 255 {
	// 	arr[i] = calc(i)
	// 	if i == 255 {
	// 		break
	// 	}
	// 	i++
	// }
	// b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := uint8(rand.Int31n(256))
		in := rand.Float32()
		v := arr[idx] * rand.Float32()
		if v < 0.0 {
			println(idx, in)
		}
	}
}

func BenchmarkProduceCalc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idx := uint8(rand.Int31n(256))
		in := rand.Float32()
		v := calc(idx) * rand.Float32()
		if v < 0.0 {
			println(idx, in)
		}
	}
}

func BenchmarkProduceFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idx := uint8(rand.Int31n(256))
		in := rand.Float32()
		v := in * rand.Float32()
		if v < 0.0 {
			println(idx, in)
		}
	}
}

// mutate
// produce
