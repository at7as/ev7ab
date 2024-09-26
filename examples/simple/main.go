// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"

	"github.com/at7as/ev7ab/app"
	"github.com/at7as/ev7ab/lab"
	"gonum.org/v1/gonum/floats"
)

const sizeSimple int = 1000

type ExampleSimple struct {
	t [][]float64
	r []int
	s float64
}

func (p *ExampleSimple) Load(setup map[string]string) error {

	p.t = make([][]float64, sizeSimple)
	p.r = make([]int, sizeSimple)
	p.s = float64(sizeSimple)

	for i := range sizeSimple {

		p.t[i] = []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		p.r[i] = int(math.Floor(floats.Sum(p.t[i])))

	}

	return nil
}

func (p *ExampleSimple) Setup(key, value string) error {

	return nil
}

func (p *ExampleSimple) Produce(n lab.Next, _ lab.Next, _ []float64) []float64 {

	total := 0
	for i, in := range p.t {
		res := 0.0
		index := 0
		for ii, v := range n(in) {
			if v > res {
				index = ii
				res = v
			}
		}
		if index == p.r[i] {
			total++
		}
	}

	return []float64{float64(total) / p.s}
}

func (p *ExampleSimple) Validate(r []float64) bool {

	if r[0] < 0.5 {
		return false
	}

	return true
}

func (p *ExampleSimple) Compare(a, b []float64) bool {

	return a[0] > b[0]
}

func (p *ExampleSimple) Best(v []float64) string {

	best := ""
	if len(v) > 0 {
		best = fmt.Sprintf("%.1f%%", v[0]*100.0)
	}

	return best
}

func (p *ExampleSimple) Goal(v []float64) bool {

	if v[0] > 0.95 {
		return true
	}

	return false
}

func try() {

	f, err := os.Open("./examples/simple/ev.lab")
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	d, err := zlib.NewReader(f)
	if err != nil {
		log.Panicln(err)
	}
	defer d.Close()

	b, err := io.ReadAll(d)
	if err != nil {
		log.Panicln(err)
	}

	l := lab.New(&ExampleSimple{})

	if err = l.Import(b); err != nil {
		log.Panicln(err)
	}

	fmt.Println(l.Value([]float64{0.1, 0.2, 0.3, 0.4}))

}

func main() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")
	runTry := flag.Bool("try", false, "try example result")

	flag.Parse()

	if *runTry {
		try()
	} else {
		app.Run(&ExampleSimple{}, *cfgFile)
	}

}
