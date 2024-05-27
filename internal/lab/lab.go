package lab

import "math/rand"

const levelMax int = 9

// Lab ...
type Lab struct {
	cfg  Config
	prod Producer
	gen  *generator
	s    state
	r    *rand.Rand
}

// Config ...
type Config struct {
	In     int
	Out    int
	Target []float64
	Limit  []float64
	Goal   bool
	Size   int
	Seed   int64
	Aggr   Aggregator
	Proc   Processor
}

type state struct {
	run  bool
	ev   *project
	next bool
	goal *entity
	// pop []*entity
	// top []*entity
}

// New ...
func New(cfg Config, prod Producer) *Lab {

	var r *rand.Rand
	if cfg.Seed != 0 {
		r = rand.New(rand.NewSource(cfg.Seed))
	} else {
		r = rand.New(rand.NewSource(rand.Int63()))
	}

	if cfg.Aggr == nil {
		cfg.Aggr = avgAggr
	}

	if cfg.Proc == nil {
		cfg.Proc = linearProc
	}

	l := &Lab{cfg: cfg, prod: prod, gen: &generator{}, r: r}

	return l
}

// Examine ...
func (l *Lab) Examine() {

	l.gen.prepare(l)
	// l.s.run = true

	for l.s.run {
		l.generation()
		// l.evolution()
	}

	// l.generate()

	// l.s.run = true
	// for l.s.run {

	// 	l.production()

	// 	if l.s.run {
	// 		l.evolution()
	// 	}

	// }

	// >>>
	// for _, p := range l.s.ev {
	// 	for _, s := range p.origin.nest {
	// 		fmt.Printf("s: %v %v %v\n", s.a, s.p, s.n)
	// 	}
	// 	fmt.Printf("p.last(): %v\n", p.last())
	// }
	// <<<
}

func (l *Lab) generation() {

	l.gen.next()

}

func (l *Lab) production() {

	// for _, s := range l.s.ev {
	// 	for _, e := range s.pop {
	// 		e.result = append(e.result, l.prod.Produce(next(e)))
	// 	}
	// 	sort.Slice(s.pop, func(i, j int) bool {
	// 		return l.prod.Compare(s.pop[i].last(), s.pop[j].last())
	// 	})
	// 	s.result = append(s.result, s.pop[0].last())
	// }

	// sort.Slice(l.s.ev, func(i, j int) bool {
	// 	return l.prod.Compare(l.s.ev[i].last(), l.s.ev[j].last())
	// })

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
	// l.generate()

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

// обучение с противником
// продумать значения параметров: начальные, конечные, шаг, динамический шаг
// debug
// top 1 del if result > max result
// массовое мышление
// старость ноды
