// Package lab provides another genetic algorithm.
package lab

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

// Lab provides Lab data.
type Lab struct {
	prod Producer
	c    Config
	s    state
}

// Config provides Lab configuration.
type Config struct {
	// Size is maximum size of generation pool and every evolution method.
	Size int
	// Aggr is code of Aggregator lab function.
	// Builtin functions is "sum", "avg", "min", "max".
	// Default is "avg".
	Aggr string
	// Proc is code of Processor lab function.
	// Builtin functions is "linear", "binary", "relu", "lrelu", "silu", "sigmoid", "softplus", "softsign", "tanh", "arctan", "sinusoid", "gaussian".
	// Default is "linear".
	Proc string
	// Goal option if true, than Lab will stop examine when target is reached.
	Goal bool
	// Duel option if true, than Lab will use challenge mode when produce result.
	Duel  bool
	debug bool
}

type state struct {
	aggr      Aggregator
	proc      Processor
	run       bool
	exec      bool
	id        int
	ev        map[int]*project
	goal      *entity
	wg        *sync.WaitGroup
	debugfile *os.File
}

// New returns new *Lab.
func New(prod Producer, debug bool) *Lab {

	c := Config{
		Size:  1000,
		Aggr:  "avg",
		Proc:  "linear",
		Goal:  false,
		Duel:  false,
		debug: debug,
	}

	s := state{
		aggr:      nil,
		proc:      nil,
		run:       false,
		id:        0,
		ev:        make(map[int]*project),
		wg:        &sync.WaitGroup{},
		debugfile: nil,
	}

	if debug {
		os.Remove("./lab.log")
		f, err := os.OpenFile("./lab.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Panicln(err)
		}
		s.debugfile = f
	}

	return &Lab{prod, c, s}
}

// Close soft closes Lab.
func (l *Lab) Close() error {

	if l.s.debugfile != nil {
		return l.s.debugfile.Close()
	}

	return nil
}

// Setup applies new Config to Lab.
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

	if !newSize {
		l.resize()
	}

	return nil
}

// GetConfig returns current Config of Lab.
func (l *Lab) GetConfig() Config {

	return l.c
}

// GetExec returns exec state of Lab.
func (l *Lab) GetExec() bool {

	return l.s.exec
}

// GetProjects returns slice of projects id.
func (l *Lab) GetProjects() []int {

	list := make([]int, 0, len(l.s.ev))
	for id := range l.s.ev {
		list = append(list, id)
	}

	return list
}

// ProjectAdd adds new project. Returns id of new project.
func (l *Lab) ProjectAdd(layout [][]Node) int {

	l.ProjectSet(l.s.id, layout)
	l.s.id++

	return l.s.id - 1
}

// ProjectSet updates project layout by id.
// Reset all data of project when update layout.
func (l *Lab) ProjectSet(id int, layout [][]Node) {

	l.s.ev[id] = newProject(l, id, layout)

}

// ProjectStatus status of project.
func (l *Lab) ProjectStatus(id int) bool {

	return l.s.ev[id].active
}

// ProjectStat returns current stats of project:
// generated count, evoluted count, age, best, goal state
func (l *Lab) ProjectStat(id int) (int, int, int, string, bool) {

	return l.s.ev[id].stat()
}

// ProjectLayout returns layout project by id.
func (l *Lab) ProjectLayout(id int) [][]Node {

	return l.s.ev[id].layout
}

// ProjectDelete delete project by id.
func (l *Lab) ProjectDelete(id int) {

	delete(l.s.ev, id)

}

// ProjectActivate activate project by id.
func (l *Lab) ProjectActivate(id int) {

	l.s.ev[id].activate()

}

// ProjectDeactivate deactivate project by id.
func (l *Lab) ProjectDeactivate(id int) {

	l.s.ev[id].deactivate()

}

// ProjectValue returns out node values of top entity project selected by id.
func (l *Lab) ProjectValue(id int, in []float64) []float64 {

	if len(l.s.ev[id].ev) > 0 {
		e := l.s.ev[id].ev[0]
		return l.s.ev[id].exec(e, in).project.value(e)
	}

	return []float64{}
}

// Run used for run Lab examine.
func (l *Lab) Run() {

	l.s.run = true
	go l.examine()

}

// Stop used for soft stop Lab examine.
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

// Value returns out node values of goal entity.
func (l *Lab) Value(in []float64) []float64 {

	return l.s.goal.project.exec(l.s.goal, in).project.value(l.s.goal)
}

// Volume returns out node values of goal entities all projects.
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

	if len(g) > 0 {

		sort.SliceStable(g, func(i, j int) bool {
			return l.prod.Compare(g[i].last(0), g[j].last(0))
		})

		l.s.goal = g[0]
	}

}

// Export returns current state of Lab encoded by gob.
func (l *Lab) Export() ([]byte, error) {

	if l.s.exec {
		return []byte{}, fmt.Errorf("Lab is not finish execution")
	}

	m := memory{
		Cfg: l.c,
		ID:  l.s.id,
		Ev:  make([]memory, len(l.s.ev)),
	}

	i := 0
	for _, p := range l.s.ev {
		m.Ev[i] = memory{
			ID:     p.id,
			Ev:     make([]memory, len(p.ev)),
			Layout: p.layout,
			Age:    p.age,
			Active: p.active,
		}
		for ii, e := range p.ev {
			m.Ev[i].Ev[ii] = memory{
				Mod:    e.mod.v,
				Out:    e.out.v,
				Result: e.result,
				Age:    e.origin,
			}
		}
		i++
	}

	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

// Import used for importing data created by Export function.
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
		l.ProjectSet(p.ID, p.Layout)
		l.s.ev[p.ID].active = p.Active
		l.s.ev[p.ID].age = p.Age
		ev := make([]*entity, len(p.Ev))
		for i, e := range p.Ev {
			mod := l.s.ev[p.ID].pool.mod.Get().(*atom)
			for _, v := range e.Mod {
				mod.v = append(mod.v, v)
			}
			ev[i] = l.s.ev[p.ID].spawn(mod, e.Result, 0)
			for _, v := range e.Out {
				ev[i].out.v = append(ev[i].out.v, v)
			}
			ev[i].origin = e.Age
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
	Active   bool
	Ev       []memory
	Layout   [][]Node
	Age      int
	Mod, Out []float64
	Result   [][]float64
}
