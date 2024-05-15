package lab

import (
	"math/rand"
	"slices"
)

type generator struct {
	lab  *Lab
	in   int
	out  int
	aggr []Aggregator
	proc []Processor
}

func (g *generator) prepare(lab *Lab) {

	g.lab = lab

	g.in = multiply(g.lab.cfg.In)
	g.out = multiply(g.lab.cfg.Out)

	g.aggr = make([]Aggregator, 0)
	if len(g.lab.cfg.Aggr) == 0 || slices.Contains(g.lab.cfg.Aggr, "sum") {
		g.aggr = append(g.aggr, sumAggr)
	}
	if len(g.lab.cfg.Aggr) == 0 || slices.Contains(g.lab.cfg.Aggr, "avg") {
		g.aggr = append(g.aggr, avgAggr)
	}
	if len(g.lab.cfg.Aggr) == 0 || slices.Contains(g.lab.cfg.Aggr, "min") {
		g.aggr = append(g.aggr, minAggr)
	}
	if len(g.lab.cfg.Aggr) == 0 || slices.Contains(g.lab.cfg.Aggr, "max") {
		g.aggr = append(g.aggr, maxAggr)
	}

	g.proc = make([]Processor, 0)
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "linear") {
		g.proc = append(g.proc, linearProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "binary") {
		g.proc = append(g.proc, binaryProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "relu") {
		g.proc = append(g.proc, reluProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "lrelu") {
		g.proc = append(g.proc, lreluProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "silu") {
		g.proc = append(g.proc, siluProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "sigmoid") {
		g.proc = append(g.proc, sigmoidProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "softplus") {
		g.proc = append(g.proc, softplusProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "softsign") {
		g.proc = append(g.proc, softsignProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "tanh") {
		g.proc = append(g.proc, tanhProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "arctan") {
		g.proc = append(g.proc, arctanProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "sinusoid") {
		g.proc = append(g.proc, sinusoidProc)
	}
	if len(g.lab.cfg.Proc) == 0 || slices.Contains(g.lab.cfg.Proc, "gaussian") {
		g.proc = append(g.proc, gaussianProc)
	}

	// for a := range g.aggr {
	// 	for p := range g.proc {
	// 		pt := &project{}
	// 		pt.origin = g.basic(a, p, pt)
	// 		g.lab.s.ev = append(g.lab.s.ev, pt)
	// 	}
	// }

}

// func (g *generator) basic(a, p int, pt *project) *model {

// 	nest := []*stage{{
// 		layout: layout{
// 			a:  a,
// 			p:  p,
// 			in: []int{0},
// 			n:  g.out,
// 		},
// 		// Aggregator: g.aggr[g.aggrList[a]],
// 		// Processor:  g.proc[g.procList[p]],
// 	}}

// 	return &model{pt, nest}
// }

func (g *generator) next() {

	g.project()

	// for _, p := range g.lab.s.ev {
	// 	for i := 0; i < 10; i++ {
	// 		p.pop = append(p.pop, g.spawn(p.origin))
	// 	}
	// }

}

func (g *generator) project() {

	if g.lab.s.ev == nil {

		g.lab.s.ev = &project{}

		st := []*stage{{
			layout: layout{
				a:  0,
				p:  0,
				in: []int{0},
				n:  g.out,
			},
			Aggregator: g.aggr[0],
			Processor:  g.proc[0],
		}}

		g.lab.s.ev.origin = &model{g.lab.s.ev, st}

		return
	}

	if g.lab.s.next {

		l := make([]int, len(g.lab.s.ev.origin.st))
		for i, s := range g.lab.s.ev.origin.st {
			l[i] = s.level
		}
		lmin := slices.Min(l)
		if lmin < levelMax {
			for i, v := range g.lab.s.ev.origin.st {
				if v.level == lmin {
					g.lab.s.ev.origin.st[i].level++
					return
				}
			}
		}

		// new project

	}

}

func (g *generator) spawn(o *model) *entity {

	e := &entity{origin: o}

	e.link = make([][][]float64, len(o.st))
	e.node = make([][]*atom, len(o.st))
	for i, s := range o.st {
		e.link[i] = make([][]float64, len(s.in))
		ini := s.in[i]
		for ii := range e.link[i] {
			if ini > 0 {
				e.link[i][ii] = make([]float64, o.st[ini].n*s.n)
			} else {
				e.link[i][ii] = make([]float64, g.in*s.n)
			}
			for iii := range e.link[i][ii] {
				e.link[i][ii][iii] = g.random(s.level)
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
