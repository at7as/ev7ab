package app

import (
	"fmt"
	"strings"
)

type position struct {
	x, y int
}

type kv struct {
	key, value string
}

func space(v string, left, right int) string {

	if left < 0 {
		left = 0
	}
	if right < 0 {
		right = 0
	}

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", left), v, strings.Repeat(" ", right))
}

func untrimRight(v string, n int) string {

	return " " + v + strings.Repeat(" ", n-len(v)-1)
}

func untrimLeft(v string, n int) string {

	return strings.Repeat(" ", n-len(v)-1) + v + " "
}

func toEven(v int) int {

	if v%2 == 0 {
		return v
	}
	return v + 1

}
