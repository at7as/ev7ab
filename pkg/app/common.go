// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"strings"
)

type position struct {
	x, y int
}

func newPosition(x, y int) position {

	return position{x, y}
}

type dict struct {
	m map[string]*setupItem
	l []*setupItem
}

type setupItem struct {
	key, value string
}

func newDict(l []*setupItem) *dict {

	d := &dict{m: make(map[string]*setupItem), l: make([]*setupItem, len(l))}
	for i, v := range l {
		d.m[v.key] = v
		d.l[i] = v
	}

	return d
}

func (d *dict) get(key string) string {

	return d.m[key].value
}

func (d *dict) set(key string, value string) {

	v, ok := d.m[key]

	if ok {
		d.m[key].value = value
		return
	}

	v = &setupItem{key, value}
	d.m[key] = v
	d.l = append(d.l, v)

}

var defaultSetup []*setupItem = []*setupItem{
	{"LabFile", "./ev.lab"},
	{"Size", "1000"},
	{"Aggr", "avg"},
	{"Proc", "linear"},
	{"Goal", "false"},
	{"Duel", "false"},
	{"In", "2"},
	{"Out", "2"},
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
