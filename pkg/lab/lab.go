package lab

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"sync"
)

// Lab ...
type Lab struct {
	prod Producer
	c    Config
	s    state
}

// Config ...
type Config struct {
	Size int
	Aggr string
	Proc string
	Goal bool
	Duel bool
}

type state struct {
	aggr Aggregator
	proc Processor
	run  bool
	exec bool
	id   int
	ev   map[int]*project
	goal *entity
	wg   *sync.WaitGroup
}

// New ...
func New(prod Producer) *Lab {

	c := Config{
		Size: 1000,
		Aggr: "avg",
		Proc: "linear",
		Goal: false,
		Duel: false,
	}

	s := state{
		aggr: nil,
		proc: nil,
		run:  false,
		id:   0,
		ev:   make(map[int]*project),
		wg:   &sync.WaitGroup{},
	}

	return &Lab{prod, c, s}
}

// Setup ...
func (l *Lab) Setup(c Config) error {

	newSize := l.c.Size == c.Size

	l.c = c

	if l.c.Aggr != "avg" {
		aggr, ok := aggrMap[l.c.Aggr]
		if !ok {
			return fmt.Errorf("%s not exists in aggr map", l.c.Aggr)
		}
		l.s.aggr = aggr
	}

	if l.c.Proc != "linear" {
		proc, ok := procMap[l.c.Proc]
		if !ok {
			return fmt.Errorf("%s not exists in proc map", l.c.Proc)
		}
		l.s.proc = proc
	}

	if newSize {
		l.resize()
	}

	return nil
}

// Load ...
func (l *Lab) Load(s map[string]string) error {

	return l.prod.Load(s)
}

// SetAggregator ...
func (l *Lab) SetAggregator(code string, aggr Aggregator) error {

	aggrMap[code] = aggr
	l.c.Aggr = code
	return l.Setup(l.c)
}

// SetProcessor ...
func (l *Lab) SetProcessor(code string, proc Processor) error {

	procMap[code] = proc
	l.c.Proc = code
	return l.Setup(l.c)
}

// AddProject ...
func (l *Lab) AddProject(layout [][]Node) int {

	l.SetProject(l.s.id, layout)
	l.s.id++

	return l.s.id - 1
}

// SetProject ...
func (l *Lab) SetProject(id int, layout [][]Node) {

	l.s.ev[id] = newProject(l, id, layout)

}

// StatProject ...
func (l *Lab) StatProject(id int) (int, int, int, []float64, bool) {

	return l.s.ev[id].stat()
}

// DelProject ...
func (l *Lab) DelProject(id int) {

	delete(l.s.ev, id)

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

	l.achieve()

	l.s.exec = false

}

// Value ...
func (l *Lab) Value(in []float64) []float64 {

	return l.s.goal.project.exec(l.s.goal, in).project.value(l.s.goal)
}

// Volume ...
func (l *Lab) Volume(in []float64) [][]float64 {

	vol := make([][]float64, 0)
	for _, p := range l.s.ev {
		if p.goal != nil {
			vol = append(vol, p.exec(p.goal, in).project.value(p.goal))
		}
	}
	return vol
}

func (l *Lab) resize() {

	for _, p := range l.s.ev {
		p.resize()
	}

}

func (l *Lab) achieve() {

	g := make([]*entity, 0)
	for _, p := range l.s.ev {
		if p.goal != nil {
			g = append(g, p.goal)
		}
	}

	sort.Slice(g, func(i, j int) bool {
		return l.prod.Compare(g[i].last(0), g[j].last(0))
	})

	l.s.goal = g[0]

}

// Export ...
func (l *Lab) Export() ([]byte, error) {

	if l.s.exec {
		return []byte{}, fmt.Errorf("lab is not finish execution")
	}

	m := memory{
		Cfg: l.c,
		ID:  l.s.id,
		Ev:  make([]memory, len(l.s.ev)),
	}
	for i, p := range l.s.ev {
		m.Ev[i] = memory{
			ID:     p.id,
			Ev:     make([]memory, len(p.ev)),
			Layout: p.layout,
			Age:    p.age,
		}
		for ii, e := range p.ev {
			m.Ev[i].Ev[ii] = memory{
				Mod:    e.mod.v,
				Out:    e.out.v,
				Result: e.result,
			}
		}
	}

	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

// Import ...
func (l *Lab) Import(data []byte) error {

	if l.s.exec {
		return fmt.Errorf("lab is not finish execution")
	}

	b := bytes.Buffer{}
	_, err := b.Write(data)
	if err != nil {
		return err
	}

	m := memory{}
	err = gob.NewDecoder(&b).Decode(&m)
	if err != nil {
		return err
	}

	err = l.Setup(m.Cfg)
	if err != nil {
		return err
	}
	l.s.id = m.ID
	for _, p := range m.Ev {
		l.SetProject(p.ID, p.Layout)
		l.s.ev[p.ID].age = p.Age
		ev := make([]*entity, len(p.Ev))
		for i, e := range p.Ev {
			mod := l.s.ev[p.ID].pool.mod.Get().(*atom)
			for _, v := range e.Mod {
				mod.v = append(mod.v, v)
			}
			ev[i] = l.s.ev[p.ID].spawn(mod)
			for _, v := range e.Out {
				ev[i].out.v = append(ev[i].out.v, v)
			}
			ev[i].result = e.Result
		}
		l.s.ev[p.ID].ev = ev
		l.s.ev[p.ID].achieve()
	}
	l.achieve()

	return nil
}

type memory struct {
	Cfg      Config
	ID       int
	Ev       []memory
	Layout   [][]Node
	Age      int
	Mod, Out []float64
	Result   [][]float64
}
