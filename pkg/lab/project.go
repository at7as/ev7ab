package lab

import (
	"math/rand/v2"
	"sort"
	"sync"
)

// Node ...
type Node struct {
	src [][2]int
	out int
}

type project struct {
	lab    *Lab
	id     int
	active bool
	layout
	model
	mod, out int
	pool     struct {
		mod *sync.Pool
		out *sync.Pool
	}
	size  int
	rand  []*rand.Rand
	wg    *sync.WaitGroup
	gen   []*entity
	ev    []*entity
	exec  func(*entity, []float64) *entity
	value func(*entity) []float64
}

type layout [][]Node

type model []node

type node struct {
	Node
	in   float64
	mods int
	modc int
}

func newProject(lab *Lab, layout [][]Node) *project {

	p := &project{
		lab:    lab,
		id:     lab.s.id,
		active: false,
		layout: layout,
		size:   lab.c.size,
		wg:     &sync.WaitGroup{},
	}
	p.compile()

	p.pool.mod = &sync.Pool{
		New: func() any { return newAtom(p.mod) },
	}
	p.pool.out = &sync.Pool{
		New: func() any { return newAtom(p.out) },
	}

	p.rand = make([]*rand.Rand, 0, p.size)
	p.resize()

	p.exec = execDefault
	if p.lab.c.aggr != nil {
		p.exec = execCustom
	}
	p.value = valueDefault
	if p.lab.c.proc != nil {
		p.value = valueCustom
	}

	return p
}

func (p *project) compile() {

	p.mod = 0
	p.out = 0
	p.model = make(model, 0)
	index := make([][]int, len(p.layout))

	for i, stage := range p.layout {
		index[i] = make([]int, len(stage))

		for ii, n := range stage {

			index[i][ii] = p.out

			src := make([][2]int, len(n.src))
			in := 0
			for iii, s := range n.src {
				is := index[s[0]][s[1]]
				o := p.layout[s[0]][s[1]].out
				src[iii] = [2]int{is, is + o}
				in += o
			}

			mods := p.mod
			modc := in * n.out
			p.mod += modc
			p.out += n.out

			p.model = append(p.model, node{
				Node: Node{
					out: n.out,
					src: src,
				},
				in:   float64(in),
				mods: mods,
				modc: modc,
			})

		}
	}

}

func (p *project) resize() {

	d := p.lab.c.size - len(p.rand)
	if d > 0 {
		for range d {
			p.rand = append(p.rand, rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())))
		}
	} else if d < 0 {
		p.rand = p.rand[:p.lab.c.size]
	}

}

func (p *project) examine() {

	for p.lab.s.run {
		p.generation()
		p.evolution()
	}

	p.lab.s.wg.Done()

}

func (p *project) generation() {

	h := &house{make([]*entity, p.size)}

	p.wg.Add(len(h.e))
	for i := range len(h.e) {
		go p.generate(p.rand[i], h, i)
	}
	p.wg.Wait()

	p.production(h)

	p.gen = h.e

}

func (p *project) generate(r *rand.Rand, h *house, index int) {

	mod := p.pool.mod.Get().(*atom)
	for range p.mod {
		mod.v = append(mod.v, r.Float64())
	}

	h.e[index] = p.spawn(mod)

	p.wg.Done()

}

func (p *project) spawn(mod *atom) *entity {

	return &entity{
		project: p,
		model:   &p.model,
		mod:     mod,
		out:     p.pool.out.Get().(*atom),
	}
}

func (p *project) evolution() {

	lenev := len(p.ev)
	lengen := len(p.gen)
	lenev1 := lenev - 1
	lengen1 := lengen - 1

	l := lenev*3 + lengen*3 + p.size*2
	c := l + lenev + lengen
	h := &house{make([]*entity, l, c)}
	index := 0

	p.wg.Add(l)
	for i, e := range p.ev {
		go p.mutate(p.rand[i], e, h, index)
		index++
		go p.variate(p.rand[i], e, h, index)
		index++
		go p.combine(p.rand[i], e, p.ev[lenev1-i], h, index)
		index++
	}
	for i, e := range p.gen {
		go p.mutate(p.rand[i], e, h, index)
		index++
		go p.variate(p.rand[i], e, h, index)
		index++
		go p.combine(p.rand[i], e, p.gen[lengen1-i], h, index)
		index++
	}
	for i := range p.size {
		go p.combine(p.rand[i], p.ev[rand.IntN(len(p.ev))], p.gen[rand.IntN(len(p.gen))], h, index)
		index++
		go p.mediate(p.ev[:rand.IntN(len(p.ev))], h, index)
		index++
	}
	p.wg.Wait()

	p.production(h)
	p.reduce(h)

	h.e = append(h.e, p.ev...)
	h.e = append(h.e, p.gen...)

	sort.Slice(h.e, func(i, j int) bool {
		return p.lab.prod.Compare(h.e[i].last(0), h.e[j].last(0))
	})

	for i := p.size; i < len(h.e); i++ {
		h.e[i].atomize()
	}

	p.ev = h.e[:p.size]

}

func (p *project) mutate(r *rand.Rand, e *entity, h *house, index int) {

	mod := e.mod.clone(p.pool.mod)

	num := r.Perm(len(p.model) - 1)[:1+r.IntN(len(p.model)-2)]
	for _, numi := range num {
		n := p.model[numi]
		for i := n.mods; i < n.mods+n.modc; i++ {
			mod.v[i] = clamp(mod.v[i] + 0.1*(r.Float64()-0.5))
		}
	}

	h.e[index] = p.spawn(mod)

	p.wg.Done()

}

func (p *project) variate(r *rand.Rand, e *entity, h *house, index int) {

	mod := e.mod.clone(p.pool.mod)

	num := r.Perm(len(p.model) - 1)[:1+r.IntN(len(p.model)-2)]
	for _, numi := range num {
		n := p.model[numi]
		for i := n.mods; i < n.mods+n.modc; i++ {
			mod.v[i] = r.Float64()
		}
	}

	h.e[index] = p.spawn(mod)

	p.wg.Done()

}

func (p *project) combine(r *rand.Rand, e1 *entity, e2 *entity, h *house, index int) {

	mod := e1.mod.clone(p.pool.mod)

	num := r.Perm(len(p.model) - 1)[:1+r.IntN(len(p.model)-2)]
	for _, numi := range num {
		n := p.model[numi]
		for i := n.mods; i < n.mods+n.modc; i++ {
			mod.v[i] = e2.mod.v[i]
		}
	}

	h.e[index] = p.spawn(mod)

	p.wg.Done()

}

func (p *project) mediate(e []*entity, h *house, index int) {

	mod := p.pool.mod.Get().(*atom)

	for range p.mod {
		mod.v = append(mod.v, 0)
	}

	for i, esrc := range e {
		w := float64((len(e) - i))
		for ii, v := range esrc.mod.v {
			mod.v[ii] += v * w
		}
	}

	tri := triangular(len(e))
	for i := range p.mod {
		mod.v[i] /= float64(tri)
	}

	h.e[index] = p.spawn(mod)

	p.wg.Done()

}

func (p *project) production(h *house) {

	p.wg.Add(len(h.e))
	for _, e := range h.e {
		go e.exec()
	}
	p.wg.Wait()

	he := make([]*entity, 0, len(h.e))
	for _, e := range h.e {
		if p.lab.prod.Validate(e.last(0)) {
			he = append(he, e)
		} else {
			e.atomize()
		}
	}
	h.e = he

}

func (p *project) reduce(h *house) {

	he := make([]*entity, 0, len(h.e))
	for _, e := range h.e {
		if len(e.result) > 1 && p.lab.prod.Compare(e.last(1), e.last(0)) {
			he = append(he, e)
		} else {
			e.atomize()
		}
	}
	h.e = he

}

func (p *project) set(layout [][]Node) {

	p.layout = layout
	p.compile()
	// not ready

}

func clamp(v float64) float64 {

	return max(0.0, min(1.0, v))
}

func triangular(n int) int {

	return n * (n + 1) / 2
}

type house struct {
	e []*entity
}

type entity struct {
	*project
	*model
	mod, out *atom
	result   [][]float64
}

func (e *entity) exec() {

	e.result = append(e.result, e.project.lab.prod.Produce(next(e)))

	e.project.wg.Done()

}

func (e *entity) last(shift int) []float64 {

	return e.result[len(e.result)-1-shift]
}

func (e *entity) atomize() {

	e.mod.v = e.mod.v[:0]
	e.project.pool.mod.Put(e.mod)

	e.out.v = e.out.v[:0]
	e.project.pool.out.Put(e.out)

}

func next(e *entity) Next {

	return func(in []float64) []float64 {

		return e.project.exec(e, in).project.value(e)
	}
}

func execDefault(e *entity, in []float64) *entity {

	e.out.v = e.out.v[:0]
	e.out.v = append(e.out.v, in...)
	mod := 0
	for _, n := range *e.model {
		for range n.out {
			v := 0.0
			for _, index := range n.src {
				for _, value := range e.out.v[index[0]:index[1]] {
					v += qlinear(value, e.mod.v[mod])
					mod++
				}
			}
			e.out.v = append(e.out.v, v/n.in)
		}
	}

	return e
}

func execCustom(e *entity, in []float64) *entity {

	e.out.v = e.out.v[:0]
	e.out.v = append(e.out.v, in...)
	mod := 0
	for _, n := range *e.model {
		for range n.out {
			v := make([]float64, 0, int(n.in))
			for _, index := range n.src {
				for _, value := range e.out.v[index[0]:index[1]] {
					v = append(v, qlinear(value, e.mod.v[mod]))
					mod++
				}
			}
			e.out.v = append(e.out.v, e.project.lab.c.aggr(v))
		}
	}

	return e
}

func qlinear(v, m float64) float64 {

	return max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
}

func valueDefault(e *entity) []float64 {

	out := make([]float64, 0)
	for _, v := range e.out.v[len(e.out.v)-(*e.model)[len(*e.model)-1].out:] {
		out = append(out, v)
	}

	return out
}

func valueCustom(e *entity) []float64 {

	out := make([]float64, 0)
	for _, v := range e.out.v[len(e.out.v)-(*e.model)[len(*e.model)-1].out:] {
		out = append(out, e.project.lab.c.proc(v))
	}

	return out
}

type atom struct {
	v []float64
}

func newAtom(size int) *atom {

	return &atom{make([]float64, 0, size)}
}

func (a *atom) clone(pool *sync.Pool) *atom {

	c := pool.Get().(*atom)
	for _, v := range a.v {
		c.v = append(c.v, v)
	}

	return c
}
