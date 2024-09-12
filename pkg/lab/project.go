package lab

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"sort"
	"sync"
)

// Node ...
type Node struct {
	Src [][2]int
	Out int
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
	age   int
	goal  *entity
	// step  int
}

type layout [][]Node

type model []node

type node struct {
	Node
	in   float64
	mods int
	modc int
}

func newProject(lab *Lab, id int, layout [][]Node) *project {

	p := &project{
		lab:    lab,
		id:     id,
		active: false,
		layout: layout,
		size:   lab.c.Size,
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
	if p.lab.s.aggr != nil {
		p.exec = execCustom
	}
	p.value = valueDefault
	if p.lab.s.proc != nil {
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

			src := make([][2]int, len(n.Src))
			in := 0
			for iii, s := range n.Src {
				is := index[s[0]][s[1]]
				o := p.layout[s[0]][s[1]].Out
				src[iii] = [2]int{is, is + o}
				in += o
			}

			mods := p.mod
			modc := in * n.Out
			p.mod += modc
			p.out += n.Out

			p.model = append(p.model, node{
				Node: Node{
					Out: n.Out,
					Src: src,
				},
				in:   float64(in),
				mods: mods,
				modc: modc,
			})

		}
	}

}

func (p *project) resize() {

	d := p.lab.c.Size - len(p.rand)
	if d > 0 {
		for range d {
			p.rand = append(p.rand, rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())))
		}
	} else if d < 0 {
		p.rand = p.rand[:p.lab.c.Size]
	}

	p.size = p.lab.c.Size

}

func (p *project) examine() {

	for p.lab.s.run {
		p.generation()
		if len(p.gen) > 0 {
			p.evolution()
		}
		p.age++
	}

	p.lab.s.wg.Done()

}

func (p *project) generation() {

	h := &house{make([]*entity, p.size)}

	p.wg.Add(p.size)
	for i := range p.size {
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

	h.e[index] = p.spawn(mod, make([][]float64, 0), 0)

	p.wg.Done()

}

func (p *project) spawn(mod *atom, result [][]float64, origin int) *entity {

	r := [][]float64{}
	if len(result) > 0 {
		rr := make([]float64, len(result[len(result)-1]))
		copy(rr, result[len(result)-1])
		r = append(r, rr)
	}

	return &entity{
		project: p,
		model:   &p.model,
		mod:     mod,
		out:     p.pool.out.Get().(*atom),
		result:  r,
		origin:  origin,
	}
}

func (p *project) evolution() {

	lenev := len(p.ev)
	lengen := len(p.gen)

	l := lenev*2 + lengen*2 + min(lenev, lengen)
	if lenev > 1 {
		l += lenev
	}
	if lengen > 1 {
		l += lengen
	}
	he := &house{make([]*entity, l)}
	index := 0

	p.wg.Add(l)
	for i := range lenev {
		e := p.ev[i]
		go p.mutate(p.rand[i], e, he, index, 2)
		index++
		go p.variate(p.rand[i], e, he, index, 3)
		index++
		if lenev > 1 {
			ii := 0
			for {
				ii = p.rand[i].IntN(lenev)
				if i != ii {
					break
				}
			}
			go p.combine(p.rand[i], e, p.ev[ii], he, index, 4)
			index++
		}
	}
	for i := range lengen {
		e := p.gen[i]
		go p.mutate(p.rand[i], e, he, index, 5)
		index++
		go p.variate(p.rand[i], e, he, index, 6)
		index++
		if lengen > 1 {
			ii := 0
			for {
				ii = p.rand[i].IntN(lengen)
				if i != ii {
					break
				}
			}
			go p.combine(p.rand[i], e, p.gen[ii], he, index, 7)
			index++
		}
	}
	for i := range min(lenev, lengen) {
		go p.combine(p.rand[i], p.ev[i], p.gen[i], he, index, 8)
		index++
	}
	p.wg.Wait()

	p.production(he)
	p.reduce(he)

	he.e = append(he.e, p.gen...)

	for _, e := range p.ev {
		e.origin = 1
	}

	if len(he.e) < p.size/2 {

		for _, e := range p.ev {
			if len(he.e) < p.size {
				he.e = append(he.e, e)
			} else {
				e.atomize()
			}
		}

		sort.SliceStable(he.e, func(i, j int) bool {
			return p.lab.prod.Compare(he.e[i].last(0), he.e[j].last(0))
		})

		p.ev = he.e

	} else {

		he.e = append(he.e, p.ev...)

		sort.SliceStable(he.e, func(i, j int) bool {
			return p.lab.prod.Compare(he.e[i].last(0), he.e[j].last(0))
		})

		counter := [9]int{}
		for _, e := range he.e {
			counter[e.origin]++
		}
		for i, v := range counter {
			counter[i] = int(math.Floor(float64(v) / float64(len(he.e)) * float64(min(len(he.e), p.size))))
		}

		h := &house{make([]*entity, 0, p.size)}
		for _, e := range he.e {
			if counter[e.origin] > 0 {
				h.e = append(h.e, e)
				counter[e.origin]--
			} else {
				e.atomize()
			}
		}

		p.ev = h.e

	}

	p.achieve()

	if p.lab.c.debug {

		s := make([][]int, 9)
		for i := range s {
			s[i] = make([]int, 2)
			s[i][0] = p.size
		}
		for i, e := range p.ev {
			s[e.origin][0] = min(s[e.origin][0], i)
			s[e.origin][1] = s[e.origin][1] + 1
		}
		gen, ev, age, best, goal := p.stat()
		text := fmt.Sprintf("id:%v gen:%v ev:%v age:%v best:%s goal:%v", p.id, gen, ev, age, best, goal)
		for i := range s {
			if s[i][1] > 0 {
				text += fmt.Sprintf(" %v:%v:%v", i, s[i][0], s[i][1])
			}
		}
		text += "\n"
		if _, err := p.lab.s.debugfile.WriteString(text); err != nil {
			log.Panicln(err)
		}

	}

}

func (p *project) mutate(r *rand.Rand, e *entity, h *house, index int, origin int) {

	mod := e.mod.clone(p.pool.mod)

	num := r.Perm(len(p.model) - 1)[:1+rintn(r, len(p.model)-1)]
	for _, numi := range num {
		n := p.model[numi+1]
		numv := r.Perm(n.modc)[:1+rintn(r, n.modc-1)]
		for _, i := range numv {
			mod.v[n.mods+i] = clamp(mod.v[n.mods+i] + 0.1*(r.Float64()-0.5))
		}
	}

	h.e[index] = p.spawn(mod, e.result, origin)

	p.wg.Done()

}

func (p *project) variate(r *rand.Rand, e *entity, h *house, index int, origin int) {

	mod := e.mod.clone(p.pool.mod)

	num := r.Perm(len(p.model) - 1)[:1+rintn(r, len(p.model)-1)]
	for _, numi := range num {
		n := p.model[numi+1]
		numv := r.Perm(n.modc)[:1+rintn(r, n.modc-1)]
		for _, i := range numv {
			mod.v[n.mods+i] = r.Float64()
		}
	}

	h.e[index] = p.spawn(mod, e.result, origin)

	p.wg.Done()

}

func (p *project) combine(r *rand.Rand, e1 *entity, e2 *entity, h *house, index int, origin int) {

	mod := e1.mod.clone(p.pool.mod)

	num := r.Perm(len(p.model) - 1)[:1+rintn(r, len(p.model)-1)]
	for _, numi := range num {
		n := p.model[numi+1]
		numv := r.Perm(n.modc)[:1+rintn(r, n.modc-1)]
		for _, i := range numv {
			mod.v[n.mods+i] = e2.mod.v[n.mods+i]
		}
	}

	hr := &house{[]*entity{e1, e2}}
	sort.SliceStable(hr.e, func(i, j int) bool {
		return p.lab.prod.Compare(hr.e[i].last(0), hr.e[j].last(0))
	})

	h.e[index] = p.spawn(mod, hr.e[0].result, origin)

	p.wg.Done()

}

func (p *project) production(h *house) {

	p.wg.Add(len(h.e))
	if p.lab.c.Duel {
		if len(h.e) > 1 {
			for i, e := range h.e {
				last := len(h.e) - 1
				if i != last {
					go e.exec(h.e[i+1])
				} else {
					go e.exec(h.e[0])
				}
			}
		}
	} else {
		for _, e := range h.e {
			go e.exec(nil)
		}
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
		if p.lab.prod.Compare(e.last(0), e.last(1)) {
			he = append(he, e)
		} else {
			e.atomize()
		}
	}
	h.e = he

}

func (p *project) activate() {

	p.active = true

}

func (p *project) deactivate() {

	p.active = false

}

func (p *project) achieve() {

	if len(p.ev) > 0 && p.lab.prod.Goal(p.ev[0].last(0)) {

		p.goal = p.ev[0]

		if p.lab.c.Goal {
			p.lab.s.run = false
		}

	}

}

func (p *project) stat() (int, int, int, string, bool) {

	result := []float64{}
	if p.goal != nil {
		result = p.goal.last(0)
	} else if len(p.ev) > 0 {
		result = p.ev[0].last(0)
	}

	return len(p.gen), len(p.ev), p.age, p.lab.prod.Best(result), p.goal != nil
}

func clamp(v float64) float64 {

	return max(0.0, min(1.0, v))
}

func rintn(r *rand.Rand, v int) int {

	if v < 2 {
		return 0
	}
	return r.IntN(v)

}
