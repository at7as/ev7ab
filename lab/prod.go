// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lab

import "slices"

// Producer provides prepare and control lab examine.
type Producer interface {
	// Load used for loading and prepare data for producing. Settings map is optional.
	Load(map[string]string) error
	// Setup used for change settings of producer in runtime.
	Setup(string, string) error
	// Produce returns result for first entity.
	// Uses second entity like opponent and previous result when Config.Duel is true.
	Produce(Next, Next, []float64) []float64
	// Validate returns proper state of entity.
	Validate([]float64) bool
	// Compare return true if first result better than second. Used for sort and filter entities.
	Compare([]float64, []float64) bool
	// Best returns formatted string of best result.
	Best([]float64) string
	// Goal returns goal state. If true, lab is reached target.
	Goal([]float64) bool
}

// Next provides wrapper function that returns value of last node after execution.
// This is a service type, usually used in Producer.Produce.
type Next func([]float64) []float64

// IndexOfMax returns index of element has maximum value of slice.
func IndexOfMax(v []float64) int {

	index := 0
	if len(v) == 0 {
		return index
	}
	vmax := v[0]
	for i, vv := range v {
		if vv > vmax {
			vmax = vv
			index = i
		}
	}

	return index
}

// IndexOfMaxExcluded returns index of element has maximum value of slice, excluding some indexes.
func IndexOfMaxExcluded(v []float64, excl []int) int {

	index := 0
	if len(v) == 0 {
		return index
	}
	for i := range v {
		if !slices.Contains(excl, i) {
			index = i
			break
		}
	}
	vmax := v[index]
	for i, vv := range v {
		if vv > vmax && !slices.Contains(excl, i) {
			vmax = vv
			index = i
		}
	}

	return index
}
