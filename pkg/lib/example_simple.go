package lib

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lab"
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

func (p *ExampleSimple) Produce(n lab.Next, op lab.Next) []float64 {

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

func (p *ExampleSimple) Compare(a, b []float64) bool {

	return a[0] > b[0]
}

func (p *ExampleSimple) Validate(r []float64) bool {

	if r[0] < 0.5 {
		return false
	}

	return true
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

func ExampleSimpleApp() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&ExampleSimple{}, *cfgFile, true)

}

func ExampleSimpleTry() {

	f, err := os.Open("./test/example_simple/ev.lab")
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

	l := lab.New(&ExampleSimple{}, false)

	if err = l.Import(b); err != nil {
		log.Panicln(err)
	}

	fmt.Println(l.Value([]float64{0.1, 0.2, 0.3, 0.4}))

}
