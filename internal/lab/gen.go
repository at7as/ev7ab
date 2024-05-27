package lab

import (
	"sort"
)

type generator struct {
	lab *Lab
}

func (g *generator) prepare(lab *Lab) {

	g.lab = lab

	// in := make([][]int, 1)
	// in[0] = []int{g.lab.cfg.In}
	// out := make([][]int, 1)
	// out[0] = []int{g.lab.cfg.Out}
	// link := make([][]int, 1)
	// link[0] = []int{g.lab.cfg.In * g.lab.cfg.Out}
	// links := 0 // []int{sum(link[0])}
	g.lab.s.ev = &project{}
	g.lab.s.ev.origin = &model{pt: g.lab.s.ev} //, in, out, link, links}

}

func (g *generator) next() {

	g.project()

	pop := make([]*entity, 0)

	for range g.lab.cfg.Size {
		e := g.spawn(g.lab.s.ev.origin)
		e.result = append(e.result, g.lab.prod.Produce(next(e)))
		if g.lab.prod.Validate(e.last()) {
			pop = append(pop, e)
		}
	}

	g.lab.s.ev.pop = append(g.lab.s.ev.pop, pop...)

	sort.Slice(g.lab.s.ev.pop, func(i, j int) bool {
		return g.lab.prod.Compare(g.lab.s.ev.pop[i].last(), g.lab.s.ev.pop[j].last())
	})

	// s.result = append(s.result, s.pop[0].last())

	// sort.Slice(l.s.ev, func(i, j int) bool {
	// 	return l.prod.Compare(l.s.ev[i].last(), l.s.ev[j].last())
	// })

	// for _, p := range g.lab.s.ev {
	// 	for i := 0; i < 10; i++ {
	// 		p.pop = append(p.pop, g.spawn(p.origin))
	// 	}
	// }

}

func (g *generator) project() {

	// if g.lab.s.ev == nil {

	// 	g.lab.s.ev = &project{}
	// 	in := make([][]int, 1)
	// 	in[0] = []int{}
	// 	out := make([][]int, 1)

	// 	// l[0] = make([][2]int, 0)
	// 	g.lab.s.ev.origin = &model{g.lab.s.ev, in, out}

	// 	// st := []*stage{{
	// 	// 	layout: layout{
	// 	// 		in: []int{0},
	// 	// 		n:  g.out,
	// 	// 	},
	// 	// }}

	// 	// g.lab.s.ev.origin = &model{g.lab.s.ev, st}

	// 	return
	// }

	// if g.lab.s.next {

	// 	l := make([]int, len(g.lab.s.ev.origin.st))
	// 	for i, s := range g.lab.s.ev.origin.st {
	// 		l[i] = s.level
	// 	}
	// 	lmin := slices.Min(l)
	// 	if lmin < levelMax {
	// 		for i, v := range g.lab.s.ev.origin.st {
	// 			if v.level == lmin {
	// 				g.lab.s.ev.origin.st[i].level++
	// 				return
	// 			}
	// 		}
	// 	}

	// 	// new project

	// }

}

func (g *generator) spawn(o *model) *entity {

	e := &entity{origin: o}

	e.link = make([]float64, e.origin.link)

	for i := range e.origin.link {
		e.link[i] = g.lab.r.Float64()
	}

	// e.link = make([][][]float64, len(o.st))
	// e.node = make([][]*atom, len(o.st))
	// for i, s := range o.st {
	// 	e.link[i] = make([][]float64, len(s.in))
	// 	ini := s.in[i]
	// 	for ii := range e.link[i] {
	// 		if ini > 0 {
	// 			e.link[i][ii] = make([]float64, o.st[ini].n*s.n)
	// 		} else {
	// 			// e.link[i][ii] = make([]float64, g.in*s.n)
	// 		}
	// 		for iii := range e.link[i][ii] {
	// 			e.link[i][ii][iii] = g.random(s.level)
	// 		}
	// 	}
	// 	e.node[i] = make([]*atom, s.n)
	// 	for ii := range e.node[i] {
	// 		e.node[i][ii] = &atom{in: make([]float64, 0)}
	// 	}
	// }

	return e
}
