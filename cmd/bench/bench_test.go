package main

import (
	"math/rand/v2"
	"testing"
)

var sampleBox []string
var sampleInt []int

func init() {

	sampleBox = make([]string, 1000)
	for i := range sampleBox {
		sampleBox[i] = randomValue(charBox)
	}
	sampleInt = make([]int, 1000)
	for i := range sampleBox {
		sampleInt[i] = randomValue(charInt)
	}

}

func randomValue[K string | int, V int | string](v map[K]V) K {
	n := rand.IntN(len(v))
	i := 0
	var r K
	for k := range v {
		if i == n {
			return k
		}
		i++
	}
	return r
}

// func BenchmarkVarBox(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range sampleBox {
// 			r := charBox[v]
// 			if r < 0 {
// 				print(r)
// 			}
// 		}
// 	}
// }

// func BenchmarkVarInt(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range sampleInt {
// 			r := charInt[v]
// 			if r == "*" {
// 				print(r)
// 			}
// 		}
// 	}
// }

// func BenchmarkFnBox(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range sampleBox {
// 			r := getCharBox(v)
// 			if r < 0 {
// 				print(r)
// 			}
// 		}
// 	}
// }

// func BenchmarkFnInt(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range sampleInt {
// 			r := getCharInt(v)
// 			if r == "*" {
// 				print(r)
// 			}
// 		}
// 	}
// }
// func BenchmarkKVBox(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range sampleBox {
// 			r := getKVCharBox(v)
// 			if r < 0 {
// 				print(r)
// 			}
// 		}
// 	}
// }

// func BenchmarkKVInt(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range sampleInt {
// 			r := getKVCharInt(v)
// 			if r == "*" {
// 				print(r)
// 			}
// 		}
// 	}
// }

func BenchmarkSwitchBox(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range sampleBox {
			r := getSwitchCharBox(v)
			if r < 0 {
				print(r)
			}
		}
	}
}

func BenchmarkSwitchInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range sampleInt {
			r := getSwitchCharInt(v)
			if r == "*" {
				print(r)
			}
		}
	}
}

var charBox map[string]int = map[string]int{
	" ": 0b0000,
	"└": 0b1100,
	"┌": 0b0110,
	"┐": 0b0011,
	"┘": 0b1001,
	"┴": 0b1101,
	"├": 0b1110,
	"┬": 0b0111,
	"┤": 0b1011,
	"│": 0b1010,
	"─": 0b0101,
	"┼": 0b1111,
}

var charInt map[int]string = map[int]string{
	0b0000: " ",
	0b1100: "└",
	0b0110: "┌",
	0b0011: "┐",
	0b1001: "┘",
	0b1101: "┴",
	0b1110: "├",
	0b0111: "┬",
	0b1011: "┤",
	0b1010: "│",
	0b0101: "─",
	0b1111: "┼",
}

func getCharBox(k string) int {

	return map[string]int{
		" ": 0b0000,
		"└": 0b1100,
		"┌": 0b0110,
		"┐": 0b0011,
		"┘": 0b1001,
		"┴": 0b1101,
		"├": 0b1110,
		"┬": 0b0111,
		"┤": 0b1011,
		"│": 0b1010,
		"─": 0b0101,
		"┼": 0b1111,
	}[k]
}

func getCharInt(k int) string {

	return map[int]string{
		0b0000: " ",
		0b1100: "└",
		0b0110: "┌",
		0b0011: "┐",
		0b1001: "┘",
		0b1101: "┴",
		0b1110: "├",
		0b0111: "┬",
		0b1011: "┤",
		0b1010: "│",
		0b0101: "─",
		0b1111: "┼",
	}[k]
}

type kvBox struct {
	key   string
	value int
}

type kvInt struct {
	key   int
	value string
}

var kvCharBox []kvBox = []kvBox{
	{" ", 0b0000},
	{"└", 0b1100},
	{"┌", 0b0110},
	{"┐", 0b0011},
	{"┘", 0b1001},
	{"┴", 0b1101},
	{"├", 0b1110},
	{"┬", 0b0111},
	{"┤", 0b1011},
	{"│", 0b1010},
	{"─", 0b0101},
	{"┼", 0b1111},
}

var kvCharInt []kvInt = []kvInt{
	{0b0000, " "},
	{0b1100, "└"},
	{0b0110, "┌"},
	{0b0011, "┐"},
	{0b1001, "┘"},
	{0b1101, "┴"},
	{0b1110, "├"},
	{0b0111, "┬"},
	{0b1011, "┤"},
	{0b1010, "│"},
	{0b0101, "─"},
	{0b1111, "┼"},
}

func getKVCharBox(k string) int {

	for _, v := range kvCharBox {
		if v.key == k {
			return v.value
		}
	}

	return 0b0000
}

func getKVCharInt(k int) string {

	for _, v := range kvCharInt {
		if v.key == k {
			return v.value
		}
	}

	return " "
}

func getSwitchCharBox(k string) int {

	switch k {
	case " ":
		return 0b0000
	case "└":
		return 0b1100
	case "┌":
		return 0b0110
	case "┐":
		return 0b0011
	case "┘":
		return 0b1001
	case "┴":
		return 0b1101
	case "├":
		return 0b1110
	case "┬":
		return 0b0111
	case "┤":
		return 0b1011
	case "│":
		return 0b1010
	case "─":
		return 0b0101
	case "┼":
		return 0b1111
	default:
		return 0b0000
	}

}

func getSwitchCharInt(k int) string {

	switch k {
	case 0b0000:
		return " "
	case 0b1100:
		return "└"
	case 0b0110:
		return "┌"
	case 0b0011:
		return "┐"
	case 0b1001:
		return "┘"
	case 0b1101:
		return "┴"
	case 0b1110:
		return "├"
	case 0b0111:
		return "┬"
	case 0b1011:
		return "┤"
	case 0b1010:
		return "│"
	case 0b0101:
		return "─"
	case 0b1111:
		return "┼"
	default:
		return " "
	}

}
