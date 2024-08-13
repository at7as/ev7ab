package lab

import "sync"

type entity struct {
	*project
	*model
	mod, out *atom
	result   [][]float64
}

type house struct {
	e []*entity
}

func (e *entity) exec(op *entity) {

	e.result = append(e.result, e.project.lab.prod.Produce(next(e), next(op)))

	e.project.wg.Done()

}

func (e *entity) last(shift int) []float64 {

	result := []float64{}
	if len(e.result) > 0 {
		result = e.result[len(e.result)-1-shift]
	}

	return result
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
			e.out.v = append(e.out.v, e.project.lab.s.aggr(v))
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
		out = append(out, e.project.lab.s.proc(v))
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
