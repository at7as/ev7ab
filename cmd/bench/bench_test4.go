package main

import (
	"bytes"
	"io"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var valuesFloat []float64 = make([]float64, 10000)
var valuesUint8 []uint8 = make([]uint8, 10)
var valuesBytes []byte = make([]byte, 10)
var rb *bytes.Reader

// var valuesBytesA [1000001]byte

func init() {
	for i := 0; i < len(valuesFloat); i++ {
		valuesFloat[i] = float64(i) / float64(len(valuesFloat))
	}
	for i := 0; i < len(valuesUint8); i++ {
		valuesUint8[i] = uint8(i)
	}
	for i := 0; i < len(valuesBytes); i++ {
		valuesBytes[i] = byte(i)
	}
	rb = bytes.NewReader(valuesBytes)
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

func Benchmark_v_uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v uint8 = 0
		for ii := 0; ii < len(valuesUint8); ii++ {
			v = add(v, valuesUint8[ii])
		}
	}
}

func Benchmark_i_uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v uint8 = 0
		for ii := range valuesUint8 {
			v = add(v, valuesUint8[ii])
		}
	}
}

func Benchmark_r_uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v uint8 = 0
		for _, vv := range valuesUint8 {
			v = add(v, vv)
		}
	}
}

func Benchmark_r_bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v byte = 0
		for _, vv := range valuesBytes {
			v &= vv
		}
	}
}

func Benchmark_uint8s(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ii := 0
		var v uint8 = 0
		for {
			v &= uint8(ii)
			ii++
			if ii == 10 {
				break
			}
		}
	}
}

func Benchmark_uint8flat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v uint8 = 0
		v |= 0 //valuesUint8[0]
		v |= 1 //valuesUint8[1]
		v |= 2 //valuesUint8[2]
		v |= 3 //valuesUint8[3]
		v |= 4 //valuesUint8[4]
		v |= 5 //valuesUint8[5]
		v |= 6 //valuesUint8[6]
		v |= 7 //valuesUint8[7]
		v |= 8 //valuesUint8[8]
		v |= 9 //valuesUint8[9]
	}
}

func Benchmark_reader(b *testing.B) {
	for i := 0; i < b.N; i++ {

		rb.Reset(valuesBytes)
		var v uint8 = 0
		for {
			b, err := rb.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			v |= b
		}

	}
}

func add(a, b uint8) uint8 {

	return a & b
}

func sum(v []float64) float64 {

	var r float64
	for i := range v {
		r += v[i]
	}

	return r
}
