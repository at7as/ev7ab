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

func space(v any, left, right int) string {

	if left < 0 {
		left = 0
	}
	if right < 0 {
		right = 0
	}

	return fmt.Sprintf("%s%v%s", strings.Repeat(" ", left), v, strings.Repeat(" ", right))
}

func lead(v any, width, right int) string {

	return space(v, width-len(fmt.Sprint(v))-right, right)
}

func trail(v any, width, left int) string {

	return space(v, left, width-len(fmt.Sprint(v))-left)
}
