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

// Lab ...
type Lab struct {
	prod Producer
	c    Config
	s    state
}

// Config ...
type Config struct {
	Size  int
	Aggr  string
	Proc  string
	Goal  bool
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

// New ...
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

// Close ...
func (l *Lab) Close() error {

	if l.s.debugfile != nil {
		return l.s.debugfile.Close()
	}

	return nil
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

	if !newSize {
		l.resize()
	}

	return nil
}

// GetConfig ...
func (l *Lab) GetConfig() Config {

	return l.c
}

// GetExec ...
func (l *Lab) GetExec() bool {

	return l.s.exec
}

// GetProjects ...
func (l *Lab) GetProjects() []int {

	list := make([]int, 0, len(l.s.ev))
	for id := range l.s.ev {
		list = append(list, id)
	}

	return list
}

// ProjectAdd ...
func (l *Lab) ProjectAdd(layout [][]Node) int {

	l.ProjectSet(l.s.id, layout)
	l.s.id++

	return l.s.id - 1
}

// ProjectSet ...
func (l *Lab) ProjectSet(id int, layout [][]Node) {

	l.s.ev[id] = newProject(l, id, layout)

}

// ProjectStatus ...
func (l *Lab) ProjectStatus(id int) bool {

	return l.s.ev[id].active
}

// ProjectStat ...
func (l *Lab) ProjectStat(id int) (int, int, int, string, bool) {

	return l.s.ev[id].stat()
}

// ProjectLayout ...
func (l *Lab) ProjectLayout(id int) [][]Node {

	return l.s.ev[id].layout
}

// ProjectDelete ...
func (l *Lab) ProjectDelete(id int) {

	delete(l.s.ev, id)

}

// ProjectActivate ...
func (l *Lab) ProjectActivate(id int) {

	l.s.ev[id].activate()

}

// ProjectDeactivate ...
func (l *Lab) ProjectDeactivate(id int) {

	l.s.ev[id].deactivate()

}

// ProjectValue ...
func (l *Lab) ProjectValue(id int, in []float64) []float64 {

	if len(l.s.ev[id].ev) > 0 {
		e := l.s.ev[id].ev[0]
		return l.s.ev[id].exec(e, in).project.value(e)
	}

	return []float64{}
}

// Run ...
func (l *Lab) Run() {

	l.s.run = true
	go l.examine()

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

	if len(g) > 0 {

		sort.SliceStable(g, func(i, j int) bool {
			return l.prod.Compare(g[i].last(0), g[j].last(0))
		})

		l.s.goal = g[0]
	}

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
