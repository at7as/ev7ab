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
	"strconv"

	"github.com/at7as/ev7ab/pkg/app"
	"github.com/at7as/ev7ab/pkg/lab"
)

const sizeBezier int = 10

type ExampleBezier struct {
	t      [][]float64
	r      [][2]point
	valid  float64
	target float64
}

func (p *ExampleBezier) Load(setup map[string]string) error {

	p.t = make([][]float64, sizeBezier)
	p.r = make([][2]point, sizeBezier)

	for i := range sizeBezier {

		p0 := point{rand.Float64(), rand.Float64()}
		p1 := point{rand.Float64(), rand.Float64()}
		p2 := point{rand.Float64(), rand.Float64()}
		p3 := point{rand.Float64(), rand.Float64()}

		r0 := bezier(p0, p1, p2, p3, 0.0)
		r1 := bezier(p0, p1, p2, p3, 0.2)
		r2 := bezier(p0, p1, p2, p3, 0.4)
		r3 := bezier(p0, p1, p2, p3, 0.6)
		r4 := bezier(p0, p1, p2, p3, 0.8)
		r5 := bezier(p0, p1, p2, p3, 1.0)

		t := []float64{
			r0.x, r0.y,
			r1.x, r1.y,
			r2.x, r2.y,
			r3.x, r3.y,
			r4.x, r4.y,
			r5.x, r5.y,
		}
		p.t[i] = t

		p.r[i] = [2]point{p1, p2}

	}

	for k, v := range setup {
		if err := p.Setup(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (p *ExampleBezier) Setup(key, value string) error {

	switch key {
	case "Valid":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		p.valid = v
	case "Target":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		p.target = v
	}

	return nil
}

func (p *ExampleBezier) Produce(n lab.Next, _ lab.Next, _ []float64) []float64 {

	d := 0.0
	for i, v := range p.t {
		r := n(v)
		d += point{r[0], r[1]}.distance(p.r[i][0])
		d += point{r[2], r[3]}.distance(p.r[i][1])
	}

	return []float64{d / float64(sizeBezier*2)}
}

func (p *ExampleBezier) Validate(r []float64) bool {

	if r[0] > p.valid {
		return false
	}

	return true
}

func (p *ExampleBezier) Compare(a, b []float64) bool {

	return a[0] < b[0]
}

func (p *ExampleBezier) Best(v []float64) string {

	best := ""
	if len(v) > 0 {
		best = fmt.Sprintf("%.3f", v[0])
	}

	return best
}

func (p *ExampleBezier) Goal(v []float64) bool {

	if v[0] < p.target {
		return true
	}

	return false
}

type point struct {
	x, y float64
}

func bezier(p0, p1, p2, p3 point, t float64) point {

	q0 := p0.lerp(p1, t)
	q1 := p1.lerp(p2, t)
	q2 := p2.lerp(p3, t)
	r0 := q0.lerp(q1, t)
	r1 := q1.lerp(q2, t)
	return r0.lerp(r1, t)
}

func (p point) lerp(d point, l float64) point {

	return point{p.x + (d.x-p.x)*l, p.y + (d.y-p.y)*l}
}

func (p point) distance(d point) float64 {

	return math.Sqrt(math.Pow(d.x-p.x, 2) + math.Pow(d.y-p.y, 2))
}

func ExampleBezierApp() {

	cfgFile := flag.String("config", "./app.config.json", "path to app config file")

	flag.Parse()

	app.Run(&ExampleBezier{}, *cfgFile, true)

}

func ExampleBezierTry() {

	f, err := os.Open("./test/example_bezier/ev.lab")
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

	l := lab.New(&ExampleBezier{}, false)

	if err = l.Import(b); err != nil {
		log.Panicln(err)
	}

	p0 := point{rand.Float64(), rand.Float64()}
	p1 := point{rand.Float64(), rand.Float64()}
	p2 := point{rand.Float64(), rand.Float64()}
	p3 := point{rand.Float64(), rand.Float64()}

	r0 := bezier(p0, p1, p2, p3, 0.0)
	r1 := bezier(p0, p1, p2, p3, 0.2)
	r2 := bezier(p0, p1, p2, p3, 0.4)
	r3 := bezier(p0, p1, p2, p3, 0.6)
	r4 := bezier(p0, p1, p2, p3, 0.8)
	r5 := bezier(p0, p1, p2, p3, 1.0)

	t := []float64{
		r0.x, r0.y,
		r1.x, r1.y,
		r2.x, r2.y,
		r3.x, r3.y,
		r4.x, r4.y,
		r5.x, r5.y,
	}

	fmt.Println(l.Value(t))

}
