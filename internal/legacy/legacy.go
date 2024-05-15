package legacy

// import (
// 	"math"
// 	"math/rand"
// 	"time"
// )

// type chaos struct {
// 	r *rand.Rand
// }

// func (c *chaos) rand101() float64 {
// 	return float64(c.r.Intn(1000001))/1000000.0*2.0 - 1.0
// }

// func (c *chaos) randUniqueInt(n int, l int) []int {
// 	a := make([]int, n)
// 	for i := range a {
// 		a[i] = i
// 	}
// 	c.r.Shuffle(n, func(i, j int) { a[i], a[j] = a[j], a[i] })
// 	return a[:l]
// }

// func (c *chaos) intn(n int) int {
// 	return c.r.Intn(n)
// }

// func clamp(v, lo, hi float64) float64 {
// 	return math.Min(math.Max(v, lo), hi)
// }

// func clamp01(v float64) float64 {
// 	return clamp(v, 0.0, 1.0)
// }

// func clamp101(v float64) float64 {
// 	return clamp(v, -1.0, 1.0)
// }

// func contains[E comparable](s []E, v E) bool {
// 	return index(s, v) >= 0
// }

// func index[E comparable](s []E, v E) int {
// 	for i, vs := range s {
// 		if v == vs {
// 			return i
// 		}
// 	}
// 	return -1
// }

// type producer interface {
// 	Produce(e *Entity) []float64
// 	Sort(a, b *Entity) bool
// }

// type Genome struct { // config for state of lab
// 	design []layout
// 	volume [][][]float64
// }

// type layout struct {
// 	ctrl [3]string
// 	out  [][][2]int
// }

// func (e *entity) compile() {
// e.model = make([]*grid, len(e.design))
// for il, l := range e.design {
// 	g := grid{}
// 	g.aggregator = e.lab.gen.aggr[l.ctrl[0]]
// 	g.processor = e.lab.gen.proc[l.ctrl[1]]
// 	g.consolidator = e.lab.gen.cons[l.ctrl[2]]
// 	g.stack = make([]*atom, len(l.out))
// 	for io, o := range l.out {
// 		a := atom{}
// 		switch l.ctrl[2] {
// 		case "input":
// 			a.in = make([]float64, 1)
// 			a.out = make([]*atom, len(o))
// 		case "link":
// 			a.in = make([]float64, 1)
// 			a.out = make([]*atom, 1)
// 		case "node":
// 			a.out = make([]*atom, len(o))
// 		case "output":
// 			a.in = make([]float64, 1)
// 			a.value = make([]float64, 1)
// 		}
// 		g.stack[io] = &a
// 	}
// 	e.model[il] = &g
// }
// for ig, g := range e.model {
// 	for ia, a := range g.stack {
// 		for i, ref := range e.design[ig].out[ia] {
// 			a.out[i] = e.model[ref[0]].stack[ref[1]]
// 		}
// 	}
// }
// for ig, g := range e.volume {
// 	for ia, a := range g {
// 		for _, v := range a {
// 			e.model[ig].stack[ia].value = append(e.model[ig].stack[ia].value, v)
// 		}
// 	}
// }
// }

// func (e *entity) value(input []float64) []float64 {
// if len(e.model) == 0 {
// 	e.compile()
// }
// ingrid := e.model[0]
// for i, v := range input {
// 	ingrid.stack[i].in[0] = v
// }
// for _, g := range e.model {
// 	for _, a := range g.stack {
// 		g.consolidator.exec(a, g.aggregator, g.processor)
// 	}
// }
// outgrid := e.model[len(e.model)-1]
// output := make([]float64, len(outgrid.stack))
// for i, a := range outgrid.stack {
// 	output[i] = a.value[0]
// }
// 	return nil // output
// }

// func (e *entity) replicate() *entity {
// 	c := entity{
// 		// pt: e.pt,
// 	}
// 	return &c
// }

// func newLayout(a, p, c string, size int) layout {
// 	return layout{
// 		ctrl: [3]string{a, p, c},
// 		out:  make([][][2]int, size),
// 	}
// }

// // Lab ...
// type Lab struct {
// 	cfg *config
// 	mob []*Entity
// 	top *Entity
// 	sel map[string][]*Entity
// }

// // NewLab ...
// func NewLab(total int, origin, recombFrom, recombTo, mutateFrom, mutateTo float64, prod producer, debug int, toptotrash bool) *Lab {
// 	l := &Lab{cfg: &config{total, int(origin * float64(total)), int(recombFrom * float64(total)), int(recombTo * float64(total)), int(mutateFrom * float64(total)), int(mutateTo * float64(total)), prod, debug, &mind{}, &chaos{rand.New(rand.NewSource(time.Now().UnixNano()))}, toptotrash}}
// 	l.mob = make([]*Entity, 0, total)
// 	l.sel = make(map[string][]*Entity)
// 	return l
// }

// // Add ...
// func (l *Lab) Add(node int, evaluater string) {
// 	l.cfg.primal.add(node, evaluater)
// }

// type config struct {
// 	total      int
// 	origin     int
// 	recombFrom int
// 	recombTo   int
// 	mutateFrom int
// 	mutateTo   int
// 	producer   producer
// 	debug      int
// 	primal     *mind
// 	rand       *chaos
// 	toptotrash bool
// }

// Configure ...
// func (l *Lab) Configure(in, out []int, aggr, proc []string, target []float64, bias []float64, goal bool) {

// 	l.cfg = Config{in, out, 1000, aggr, proc, target, bias, goal}
// 	l.gen.prepare(l)

// }

// AddAggregator ...
// func (l *Lab) AddAggregator(key string, aggr Aggregator) {

// 	l.gen.aggr[key] = aggr
// 	l.gen.aggrListBuild()

// }

// AddProcessor ...
// func (l *Lab) AddProcessor(key string, proc Processor) {

// 	l.gen.proc[key] = proc
// 	l.gen.procListBuild()

// }
