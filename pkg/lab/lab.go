package lab

import (
	"fmt"
	"strconv"
	"sync"
)

// Lab ...
type Lab struct {
	prod Producer
	c    config
	s    state
}

type config struct {
	goal bool
	size int
	aggr Aggregator
	proc Processor
}

type state struct {
	run  bool
	exec bool
	id   int
	ev   map[int]*project
	wg   *sync.WaitGroup
}

// New ...
func New(prod Producer) *Lab {

	c := config{
		goal: false,
		size: 1000,
		aggr: nil,
		proc: nil,
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

// Load ...
func (l *Lab) Load(s [][2]string) error {

	return l.prod.Load(s)
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

	l.SetProject(l.s.id, layout)
	l.s.id++

	return l.s.id
}

// SetProject ...
func (l *Lab) SetProject(id int, layout [][]Node) {

	l.s.ev[id] = newProject(l, layout)

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

// Producer ...
type Producer interface {
	Load([][2]string) error
	Produce(Next) []float64
	Compare([]float64, []float64) bool
	Validate([]float64) bool
	Best([]float64) string
	Goal([]float64) bool
}

// обучение с противником -- challange mode
// продумать значения параметров: начальные, конечные, шаг, динамический шаг
// debug
// массовое мышление
// старость ноды
// best/top -- cut
