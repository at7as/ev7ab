package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

// ResultViewWidget ...
type ResultViewWidget struct{}

// NewResultViewWidget ...
func NewResultViewWidget() *ResultViewWidget {
	return &ResultViewWidget{}
}

// Layout ...
func (w *ResultViewWidget) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("result", -1, 0, maxX, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	v.Frame = true
	v.Title = " Result "

	fmt.Fprintf(v, "  %s  %s│\n", space(5), space(5))
	fmt.Fprintf(v, "  %s##%s│", space(5), space(5))
	fmt.Fprintf(v, "%sSIZE%s│", space(4), space(4))
	fmt.Fprintf(v, "%sPOP%s│", space(4), space(5))
	fmt.Fprintf(v, "%sTRY%s│", space(4), space(5))
	fmt.Fprintf(v, "%sBEST%s│", space(4), space(4))
	if app.result.holded {
		fmt.Fprintf(v, "%s\033[33;4mSTATUS\033[0m%s\n", space(3), space(3))
	} else {
		fmt.Fprintf(v, "%sSTATUS%s\n", space(3), space(3))
	}

	cnt := 0
	for i, p := range app.result.l {
		if (!app.result.holded && p.status == psActive) || (app.result.holded && p.status != psTerminated) {
			cnt++
			fmt.Fprintf(v, "  %s│%s┼%s┼%s┼%s┼%s\n", strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12))

			selected := " "
			if p.ui.selected {
				selected = "█"
			}
			if i == app.result.ai && !app.run {
				fmt.Fprintf(v, "%s \033[37;7m %v\033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m %v\033[0m\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

			} else {
				if p.status == psActive {
					fmt.Fprintf(v, "%s  %v│%v │%v │%v │%v │ %v\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

				} else if p.status == psHolded {
					fmt.Fprintf(v, "%s  \033[30;1m%v\033[0m│\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │ \033[33;4m%v\033[0m\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

				} else {
					fmt.Fprintf(v, "%s  \033[30;1m%v\033[0m│\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │\033[30;1m%v\033[0m │ \033[31;1m%v\033[0m\n", selected, untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))

				}
			}
		}
	}

	if cnt == 0 {
		fmt.Fprintf(v, "  %s│%s┼%s┼%s┼%s┼%s\n", strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12))
	}

	for range maxY - cnt*2 - 5 {
		fmt.Fprintf(v, "  %s  %s│\n", space(5), space(5))
	}

	return nil
}

type projectList struct {
	l      []*project
	a      *project
	ai     int
	holded bool
}

func newProjectList() *projectList {

	return &projectList{l: make([]*project, 0)}
}

func (pl *projectList) add(p *project) {

	pl.l = append(pl.l, p)
	if pl.a == nil {
		pl.a = p
		pl.ai = 0
	}

}

func (pl *projectList) assignActive(ai int) {

	pl.ai = ai
	if len(pl.l) > 0 {
		pl.a = pl.l[pl.ai]
	}

}

func (pl *projectList) selectUp(start int) {

	ai := pl.ai - start
	for ai >= 0 {
		if !pl.holded && pl.l[ai].status == psActive {
			pl.assignActive(ai)
			break
		}
		if pl.holded && pl.l[ai].status != psTerminated {
			pl.assignActive(ai)
			break
		}
		ai--
	}

}

func (pl *projectList) selectDown(start int) {

	ai := pl.ai + start
	for ai < len(pl.l) {
		if !pl.holded && pl.l[ai].status == psActive {
			pl.assignActive(ai)
			break
		}
		if pl.holded && pl.l[ai].status != psTerminated {
			pl.assignActive(ai)
			break
		}
		ai++
	}

}

func (pl *projectList) next() {

	ai := pl.ai
	pl.selectDown(0)
	if ai != pl.ai {
		return
	}
	pl.selectUp(0)

}

func (pl *projectList) selected() bool {

	for _, v := range app.result.l {
		if v.ui.selected {
			return true
		}
	}

	return false
}

func (pl *projectList) count() int {

	cnt := 0
	for _, v := range app.result.l {
		if v.status == psActive || v.status == psHolded {
			cnt++
		}
	}

	return cnt
}

func (pl *projectList) activated() int {

	cnt := 0
	for _, v := range app.result.l {
		if v.status == psActive {
			cnt++
		}
	}

	return cnt
}

type project struct {
	id     int
	status projectStatus
	ui     projectUI
	o      *project
	d      *projectData
	m      *projectModel
	n      *projectModel
}

func newProject(o *project) (*project, error) {

	in, err := strconv.Atoi(app.setup.m["In"].value)
	if err != nil {
		in = 0
	}
	out, err := strconv.Atoi(app.setup.m["Out"].value)
	if err != nil {
		out = 0
	}

	n := newProjectModel(o, in, out)
	p := project{
		id:     0,
		status: psNew,
		o:      o,
		d: &projectData{
			size: n.size,
			pop:  0,
			try:  0,
			best: 0.0,
		},
		n: n,
	}

	return &p, nil
}

func (p *project) edit() {

	p.n = cloneModel(p.m)
	p.n.dirty = true
	p.n.measure()

}

func (p *project) validate() bool {

	for i := range p.n.model {
		for ii := range p.n.model[i].stage {
			p.n.model[i].stage[ii].valid = i == len(p.n.model)-1
		}
	}

	for i := range p.n.model {
		for ii := range p.n.model[i].stage {
			for _, src := range p.n.model[i].stage[ii].source {
				p.n.model[src.s].stage[src.n].valid = true
			}
		}
	}

	for i := range p.n.model {
		for ii := range p.n.model[i].stage {
			if !p.n.model[i].stage[ii].valid {
				continue
			}
			p.n.model[i].stage[ii].valid = p.n.model[i].stage[ii].size > 0
			if !p.n.model[i].stage[ii].valid {
				continue
			}
			if i > 0 {
				p.n.model[i].stage[ii].valid = len(p.n.model[i].stage[ii].source) > 0
			}
		}
	}

	for _, m := range p.n.model {
		if len(m.stage) == 0 {
			return false
		}
		for _, n := range m.stage {
			if !n.valid {
				return false
			}
		}
	}

	return true
}

func (p *project) save() {

	if p.status == psNew {
		p.id = app.state.id
		app.state.id++
		app.result.add(p)
	}

	if app.run {
		p.status = psHolded
	} else {
		p.status = psActive
	}

	if !p.validate() {
		p.status = psInvalid
		app.invalid = true
		go hideInvalid()
	}

	p.m = p.n
	p.m.dirty = false

}

type projectUI struct {
	selected bool
}

type projectStatus int

const (
	psNew projectStatus = iota
	psActive
	psHolded
	psInvalid
	psTerminated
)

func (ps projectStatus) text() string {

	switch ps {
	case 0:
		return "new"
	case 1:
		return "active"
	case 2:
		return "holded"
	case 3:
		return "invalid"
	case 4:
		return "terminated"
	}

	return ""
}

type projectData struct {
	size int
	pop  int
	try  int
	best float64
}

type projectModel struct {
	model  []*projectModelStage
	dirty  bool
	size   int
	volume int
}

type projectModelStage struct {
	stage  []*projectModelNode
	size   int
	volume int
}

type projectModelNode struct {
	source []projectModelSource
	size   int
	volume int
	valid  bool
}

type projectModelSource struct {
	s, n int
}

func newProjectModel(o *project, in, out int) *projectModel {

	m := &projectModel{}
	if o != nil {
		m = cloneModel(o.m)
	} else {
		m.model = make([]*projectModelStage, 2)
		m.model[0] = &projectModelStage{stage: make([]*projectModelNode, 1)}
		m.model[0].stage[0] = &projectModelNode{source: make([]projectModelSource, 0), size: in}
		m.model[1] = &projectModelStage{stage: make([]*projectModelNode, 1)}
		m.model[1].stage[0] = &projectModelNode{source: []projectModelSource{{0, 0}}, size: out}
	}
	m.dirty = true
	m.measure()

	return m
}

func (m *projectModel) measure() {

	m.size = 0
	m.volume = 0
	for i := range m.model {
		m.model[i].size = 0
		m.model[i].volume = 0
		for ii := range m.model[i].stage {
			m.model[i].stage[ii].volume = 0
		}
	}

	for i, s := range m.model {
		for ii, n := range s.stage {
			m.model[i].size += n.size
			m.size += n.size
			for _, src := range n.source {
				v := m.model[src.s].stage[src.n].size * n.size
				m.model[i].stage[ii].volume += v
				m.model[i].volume += v
				m.volume += v
			}
		}
	}

}

func (m *projectModel) addStage(index int) {

	m.model = append(m.model[:index+1], m.model[index:]...)
	m.model[index] = &projectModelStage{stage: make([]*projectModelNode, 0)}
	for i := range m.model {
		for ii := range m.model[i].stage {
			for iii, src := range m.model[i].stage[ii].source {
				if index <= src.s {
					m.model[i].stage[ii].source[iii].s++
				}
			}
		}
	}

}

func (m *projectModel) removeStage(index int) {

	if index != 0 && index != len(m.model)-1 {
		m.model = append(m.model[:index], m.model[index+1:]...)
		for i := range m.model {
			for ii := range m.model[i].stage {
				for iii, src := range m.model[i].stage[ii].source {
					if index <= src.s {
						m.model[i].stage[ii].source[iii].s--
					}
				}
			}
		}
	}

}

func (s *projectModelStage) addNode() {

	s.stage = append(s.stage, &projectModelNode{source: make([]projectModelSource, 0)})
}

func (s *projectModelStage) removeNode(index int) {

	if len(s.stage) > index {
		s.stage = append(s.stage[:index], s.stage[index+1:]...)
	}

}

func space(n int) string {

	return strings.Repeat(" ", n)
}

func untrimRight(v string, n int) string {

	return v + strings.Repeat(" ", n-len(v))
}

func untrimLeft(v string, n int) string {

	return strings.Repeat(" ", n-len(v)) + v
}

func cloneModel(m *projectModel) *projectModel {

	n := &projectModel{}
	n.model = make([]*projectModelStage, len(m.model))
	for i, s := range m.model {
		n.model[i] = &projectModelStage{stage: make([]*projectModelNode, len(s.stage))}
		for ii, node := range s.stage {
			n.model[i].stage[ii] = &projectModelNode{
				source: make([]projectModelSource, len(node.source)),
				size:   node.size,
			}
			for iii, src := range node.source {
				n.model[i].stage[ii].source[iii] = projectModelSource{src.s, src.n}
			}
		}
	}

	return n
}
