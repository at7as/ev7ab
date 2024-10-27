// Copyright 2024 The ev7ab Authors.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lab

import (
	"sync"
)

type entity struct {
	*project
	*model
	mod, act, out *atom
	result        [][]float64
	origin        int
}

type house struct {
	e []*entity
}

func (e *entity) exec(op *entity) {

	e.result = append(e.result, e.project.lab.prod.Produce(next(e), next(op), e.last(0)))

	e.project.wg.Done()

}

func (e *entity) last(shift int) []float64 {

	result := []float64{}
	if len(e.result) > shift {
		result = e.result[len(e.result)-1-shift]
	}

	return result
}

func (e *entity) atomize() {

	e.mod.v = e.mod.v[:0]
	e.project.pool.mod.Put(e.mod)

	e.act.v = e.act.v[:0]
	e.project.pool.act.Put(e.act)

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
	act := 0
	for _, n := range *e.model {
		if len(n.Src) > 0 {
			for range n.Out {
				v := 0.0
				for _, index := range n.Src {
					for _, value := range e.out.v[index[0]:index[1]] {
						v += value * e.mod.v[mod]
						mod++
					}
				}
				e.out.v = append(e.out.v, qlinear(v/n.in, e.act.v[act]))
				act++
			}
		}
	}

	return e
}

func execCustom(e *entity, in []float64) *entity {

	e.out.v = e.out.v[:0]
	e.out.v = append(e.out.v, in...)
	mod := 0
	for _, n := range *e.model {
		if len(n.Src) > 0 {
			for range n.Out {
				v := make([]float64, 0, int(n.in))
				for _, index := range n.Src {
					for _, value := range e.out.v[index[0]:index[1]] {
						v = append(v, qlinear(value, e.mod.v[mod]))
						mod++
					}
				}
				e.out.v = append(e.out.v, e.project.lab.s.aggr(v))
			}
		}
	}

	return e
}

func qlinear(v, m float64) float64 {

	return max(1-((v-m)/0.5)*((v-m)/0.5), 0.0)
}

func valueDefault(e *entity) []float64 {

	out := make([]float64, 0)
	out = append(out, e.out.v[len(e.out.v)-(*e.model)[len(*e.model)-1].Out:]...)

	return out
}

func valueCustom(e *entity) []float64 {

	out := make([]float64, 0)
	for _, v := range e.out.v[len(e.out.v)-(*e.model)[len(*e.model)-1].Out:] {
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
	c.v = append(c.v, a.v...)

	return c
}
