package lab

import (
	"math/rand"
)

type generator struct {
	lab  *Lab
	in   int
	out  int
	aggr map[string]Aggregator
	proc map[string]Processor
}

func (g *generator) prepare(lab *Lab) {

	g.lab = lab

	g.in = multiply(g.lab.cfg.in)
	g.out = multiply(g.lab.cfg.out)

	g.aggr = make(map[string]Aggregator)
	g.aggr["sum"] = &sumAggr{}
	g.aggr["avg"] = &avgAggr{}
	g.aggr["min"] = &minAggr{}
	g.aggr["max"] = &maxAggr{}
	if len(g.lab.cfg.aggr) > 0 {
		aggr := make(map[string]Aggregator)
		for _, key := range g.lab.cfg.aggr {
			aggr[key] = g.aggr[key]
		}
		g.aggr = aggr
	}

	g.proc = make(map[string]Processor)
	g.proc["linear"] = &linearProc{}
	g.proc["binary"] = &binaryProc{}
	g.proc["relu"] = &reLUProc{}
	g.proc["lrelu"] = &lReLUProc{}
	g.proc["silu"] = &siLUProc{}
	g.proc["sigmoid"] = &sigmoidProc{}
	g.proc["softplus"] = &softPlusProc{}
	g.proc["softsign"] = &softSignProc{}
	g.proc["tanh"] = &tanhProc{}
	g.proc["arctan"] = &arcTanProc{}
	g.proc["sinusoid"] = &sinusoidProc{}
	g.proc["gaussian"] = &gaussianProc{}
	if len(g.lab.cfg.proc) > 0 {
		proc := make(map[string]Processor)
		for _, key := range g.lab.cfg.proc {
			proc[key] = g.proc[key]
		}
		g.proc = proc
	}

	for a := range g.aggr {
		for p := range g.proc {
			pt := &project{}
			pt.origin = g.basic(a, p, pt)
			g.lab.s.ev = append(g.lab.s.ev, pt)
		}
	}

}

func (g *generator) basic(a, p string, pt *project) *model {

	nest := []*stage{{
		layout: layout{
			a:  a,
			p:  p,
			in: []int{0},
			n:  g.out,
		},
		Aggregator: g.aggr[a],
		Processor:  g.proc[p],
	}}

	return &model{pt, nest}
}

func (g *generator) next() {

	for _, p := range g.lab.s.ev {
		for i := 0; i < 10; i++ {
			p.pop = append(p.pop, g.spawn(p.origin))

		}
	}

}

func (g *generator) spawn(o *model) *entity {

	e := &entity{origin: o}

	e.link = make([][][]float64, len(o.nest))
	e.node = make([][]*atom, len(o.nest))
	for i, s := range o.nest {
		e.link[i] = make([][]float64, len(s.in))
		ini := s.in[i]
		for ii := range e.link[i] {
			if ini > 0 {
				e.link[i][ii] = make([]float64, o.nest[ini].n*s.n)
			} else {
				e.link[i][ii] = make([]float64, g.in*s.n)
			}
			for iii := range e.link[i][ii] {
				e.link[i][ii][iii] = g.random(o.pt.level)
			}
		}
		e.node[i] = make([]*atom, s.n)
		for ii := range e.node[i] {
			e.node[i][ii] = &atom{in: make([]float64, 0)}
		}
	}

	return e
}

func (g *generator) random(level int) float64 {

	return -1.0 + float64(rand.Intn(randomValues[level].n))*randomValues[level].step
}

type randomValue struct {
	n    int
	step float64
}

var randomValues = []randomValue{
	{3, 1.0},
	{5, 0.5},
	{9, 0.25},
	{21, 0.1},
	{41, 0.05},
	{101, 0.02},
	{201, 0.01},
	{401, 0.005},
	{1001, 0.002},
	{2001, 0.001},
}

// 0:
// create project
// gen only random ???
// produce and sort result
// check result for goal
// if goal s.goal = e, if cfg.goal then end
// if not goal gen ???
// 1:
// produce and sort result
// check result for goal
// if goal and result better s.goal = e, if cfg.goal then end
// if not goal
// if result better or (not better or worse less than ??? times) prev then gen ??? then 1:
// if result not better or worse prev more than ??? times then 2:
// 2:

// * gen(random, mutate, recomb)
// ??? * count -- init when new
// ??? * attempts w\o rise
// 9 level -- 0 when new

// new (
// 4 aggr
// 12 proc
// * size
// )
