package lab

import (
	"fmt"
	"sort"
)

// Lab ...
type Lab struct {
	cfg  config
	prod Producer
	gen  *generator
	s    state
}

type config struct {
	in     []int
	out    []int
	aggr   []string
	proc   []string
	target []float64
	bias   []float64
	goal   bool
}

type state struct {
	ev []*project
	// pop []*entity
	// top []*entity
	run  bool
	goal *entity
}

// New ...
func New(prod Producer) *Lab {

	l := &Lab{prod: prod, gen: &generator{}}

	return l
}

// Configure ...
func (l *Lab) Configure(in, out []int, aggr, proc []string, target []float64, bias []float64, goal bool) {

	l.cfg = config{in, out, aggr, proc, target, bias, goal}
	l.gen.prepare(l)

}

// AddAggregator ...
func (l *Lab) AddAggregator(key string, aggr Aggregator) {

	l.gen.aggr[key] = aggr

}

// AddProcessor ...
func (l *Lab) AddProcessor(key string, proc Processor) {

	l.gen.proc[key] = proc

}

// Examine ...
func (l *Lab) Examine() {

	l.generate()

	l.s.run = true
	for l.s.run {

		l.production()

		if l.s.run {
			l.evolution()
		}

	}

	// >>>
	for _, p := range l.s.ev {
		for _, s := range p.origin.nest {
			fmt.Printf("s: %v %v %v\n", s.a, s.p, s.n)
		}
		fmt.Printf("p.last(): %v\n", p.last())
	}
	// <<<
}

func (l *Lab) generate() {

	l.gen.next()

}

func (l *Lab) production() {

	for _, s := range l.s.ev {
		for _, e := range s.pop {
			e.result = append(e.result, l.prod.Produce(next(e)))
		}
		sort.Slice(s.pop, func(i, j int) bool {
			return l.prod.Compare(s.pop[i].last(), s.pop[j].last())
		})
		s.result = append(s.result, s.pop[0].last())
	}

	sort.Slice(l.s.ev, func(i, j int) bool {
		return l.prod.Compare(l.s.ev[i].last(), l.s.ev[j].last())
	})

	// >>>
	// l.s.top = l.s.pop
	// <<<
}

func next(e *entity) Next {

	return func(in []float64) []float64 {

		return e.exec(in).value()
	}
}

func (l *Lab) evolution() {

	l.selection()
	l.origination()
	l.recombination()
	l.mutation()
	l.generate()

	// >>>
	l.s.run = false
	// <<<
}

// Value ...
func (l *Lab) Value(in []float64) [][]float64 {

	out := make([][]float64, 0)

	// out := make([][]float64, len(l.s.top))
	// for i, e := range l.s.top {
	// 	e.exec(in)
	// 	out[i] = e.value()
	// }

	return out
}

// Export ...
func (l *Lab) Export() {}

// Import ...
func (l *Lab) Import() {
	// spawn to top to continue examine or get value
}

// func (l *Lab) Spawn(gen Genome) *Entity {
// 	return &Entity{Genome: gen, lab: l}
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

func (l *Lab) selection() {
	// 	l.sel["origin"] = l.mob[:l.cfg.origin]
	// 	l.sel["recomb"] = l.mob[:l.cfg.recombFrom]
	// 	l.sel["mutate"] = l.mob[:l.cfg.mutateFrom]
	// 	l.mob = l.mob[:0]
}

func (l *Lab) origination() {
	// 	for _, e := range l.sel["origin"] {
	// 		l.mob = append(l.mob, e)
	// 	}
	// 	if l.cfg.toptotrash {
	// 		l.mob[0] = l.mob[len(l.mob)-1]
	// 	}
}

func (l *Lab) recombination() {
	// 	for ei := 0; ei < l.cfg.recombTo; ei++ {
	// 		e := cloneEntity(l.sel["recomb"][l.cfg.rand.intn(len(l.sel["recomb"]))])
	// 		es := l.sel["recomb"][l.cfg.rand.intn(len(l.sel["recomb"]))]
	// 		for i, v := range *es.node {
	// 			s := l.cfg.rand.randUniqueInt(len(v), l.cfg.rand.intn(len(v)))
	// 			for ii, vv := range v {
	// 				if contains(s, ii) {
	// 					copy((*e.node)[i][ii], vv)
	// 				}
	// 			}
	// 		}
	// 		l.mob = append(l.mob, e)
	// }
}

func (l *Lab) mutation() {
	// 	for ei := 0; ei < l.cfg.mutateTo; ei++ {
	// 		e := cloneEntity(l.sel["mutate"][l.cfg.rand.intn(len(l.sel["mutate"]))])
	// 		for i, v := range *e.node {
	// 			s := l.cfg.rand.randUniqueInt(len(v), l.cfg.rand.intn(len(v)))
	// 			for ii, vv := range v {
	// 				if contains(s, ii) {
	// 					ss := l.cfg.rand.randUniqueInt(len(vv), l.cfg.rand.intn(len(vv)))
	// 					for iii, vvv := range vv {
	// 						if contains(ss, iii) {
	// 							(*e.node)[i][ii][iii] = clamp101(vvv + (l.cfg.rand.rand101() * 0.1))
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 		l.mob = append(l.mob, e)
	// 	}
}

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

// обучение с противником
// продумать значения параметров: начальные, конечные, шаг, динамический шаг
// debug
// top 1 del if result > max result
