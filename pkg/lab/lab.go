package lab

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Lab ...
type Lab struct {
	prod Producer
	c    config
	s    state
}

type config struct {
	target []float64
	limit  []float64
	goal   bool
	size   int
	aggr   Aggregator
	proc   Processor
}

type state struct {
	run  bool
	exec bool
	id   int
	ev   map[int]*project
	wg   *sync.WaitGroup
	goal *project
}

// New ...
func New(prod Producer) *Lab {

	c := config{
		target: make([]float64, 0),
		limit:  make([]float64, 0),
		goal:   false,
		size:   1000,
		aggr:   nil,
		proc:   nil,
	}

	s := state{
		run: false,
		id:  0,
		ev:  make(map[int]*project),
		wg:  &sync.WaitGroup{},
	}

	return &Lab{prod, c, s}
}

// Setup ...
func (l *Lab) Setup(s [][2]string) error {

	for _, v := range s {
		switch v[0] {

		case "Target":
			for _, str := range strings.Split(v[1], " ") {
				value, err := strconv.ParseFloat(str, 64)
				if err != nil {
					return err
				}
				l.c.target = append(l.c.target, value)
			}

		case "Limit":
			for _, str := range strings.Split(v[1], " ") {
				value, err := strconv.ParseFloat(str, 64)
				if err != nil {
					return err
				}
				l.c.limit = append(l.c.limit, value)
			}

		case "Goal":
			value, err := strconv.ParseBool(v[1])
			if err != nil {
				return err
			}
			l.c.goal = value

		case "Size":
			value, err := strconv.Atoi(v[1])
			if err != nil {
				return err
			}
			l.c.size = value
			l.resize()

		case "Aggr":
			if v[1] == "avg" {
				return nil
			}
			value, ok := aggrMap[v[1]]
			if !ok {
				return fmt.Errorf("%s not exists in aggr map", v[1])
			}
			l.c.aggr = value

		case "Proc":
			if v[1] == "linear" {
				return nil
			}
			value, ok := procMap[v[1]]
			if !ok {
				return fmt.Errorf("%s not exists in proc map", v[1])
			}
			l.c.proc = value

		}
	}

	return nil
}

// SetAggregator ...
func (l *Lab) SetAggregator(code string, aggr Aggregator) error {

	aggrMap[code] = aggr
	return l.Setup([][2]string{{"Aggr", code}})
}

// SetProcessor ...
func (l *Lab) SetProcessor(code string, proc Processor) error {

	procMap[code] = proc
	return l.Setup([][2]string{{"Proc", code}})
}

// AddProject ...
func (l *Lab) AddProject(layout [][]Node) int {

	l.s.ev[l.s.id] = newProject(l, layout)
	l.s.id++

	return l.s.id
}

// SetProject ...
func (l *Lab) SetProject(id int, layout [][]Node) {

	// l.s.pop[id] = newProject(id, layout) // edit project with  struct
	// truncate pop
	// clear stat

}

// Run ...
func (l *Lab) Run() {

	l.s.run = true
	l.examine()

}

// Stop ...
func (l *Lab) Stop() {

	l.s.run = false

}

func (l *Lab) examine() {

	l.s.exec = true

	for _, p := range l.s.ev {
		if p.active {
			l.s.wg.Add(1)
			go p.examine()
		}
	}
	l.s.wg.Wait()

	l.s.exec = false

}

func (l *Lab) resize() {

	for _, p := range l.s.ev {
		p.resize()
	}

}

// Export ...
func (l *Lab) Export() {}

// Import ...
func (l *Lab) Import() {
	// spawn to top to continue examine or get value
}

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

// обучение с противником -- challange mode
// продумать значения параметров: начальные, конечные, шаг, динамический шаг
// debug
// массовое мышление
// старость ноды
// best/top -- cut
