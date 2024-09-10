package lib

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lab"
)

type ExampleSimple struct{}

func (p *ExampleSimple) Load(setup map[string]string) error {

	return nil
}

func (p *ExampleSimple) Setup(key, value string) error {

	return nil
}

func (p *ExampleSimple) Produce(n lab.Next, op lab.Next) []float64 {
	r := n([]float64{0.1, 0.2})
	d := math.Abs((r[0] / 0.3) - 1.0)
	return []float64{d, r[0]}
}

func (p *ExampleSimple) Challange(n1 lab.Next, n2 lab.Next) []float64 {

	return []float64{}
}

func (p *ExampleSimple) Compare(a, b []float64) bool {
	return a[0] < b[0]
}

func (p *ExampleSimple) Validate(r []float64) bool {

	if r[0] > 1.0 {
		return false
	}

	return true
}

func (p *ExampleSimple) Best(v []float64) string {

	best := ""
	if len(v) > 0 {
		best = fmt.Sprintf("%.2f", v[0])
	}

	return best
}

func (p *ExampleSimple) Goal(v []float64) bool {

	if v[0] < 0.01 {
		return true
	}

	return false
}

func ExampleSimpleApp() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&ExampleSimple{}, *cfgFile)

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

	l := lab.New(&ExampleSimple{})

	if err = l.Import(b); err != nil {
		log.Panicln(err)
	}

	fmt.Println(l.Value([]float64{0.1, 0.2}))
	fmt.Println(l.Volume([]float64{0.1, 0.2}))

}
