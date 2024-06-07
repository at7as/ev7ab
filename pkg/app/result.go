package app

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
	fmt.Fprintf(v, "%sSTATUS%s\n", space(3), space(3))

	for i, p := range app.result.l {
		fmt.Fprintf(v, "  %s│%s┼%s┼%s┼%s┼%s\n", strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12), strings.Repeat("─", 12))
		if i == app.result.ai && !app.run {
			fmt.Fprintf(v, "  \033[37;7m %v\033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m%v \033[0m│\033[37;7m %v\033[0m\n", untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))
		} else {
			fmt.Fprintf(v, "   %v│%v │%v │%v │%v │ %v\n", untrimRight(fmt.Sprint(p.id), 11), untrimLeft(fmt.Sprint(p.d.size), 11), untrimLeft(fmt.Sprint(p.d.pop), 11), untrimLeft(fmt.Sprint(p.d.try), 11), untrimLeft(fmt.Sprint(p.d.best), 11), untrimRight(p.status.text(), 11))
		}
	}

	for range maxY - len(app.result.l)*2 - 5 {
		fmt.Fprintf(v, "  %s  %s│\n", space(5), space(5))
	}

	return nil
}

type projectList struct {
	l  []*project
	a  *project
	ai int
}

func newProjectList() *projectList {

	return &projectList{l: make([]*project, 0)}
}

func (pl *projectList) assignActive() {

	if len(pl.l) > 0 {
		pl.a = pl.l[pl.ai]
	}

}

func (pl *projectList) selectUp() {

	if pl.ai > 0 {
		pl.ai--
	}
	pl.assignActive()

}

func (pl *projectList) selectDown() {

	if pl.ai < len(pl.l)-1 {
		pl.ai++
	}
	pl.assignActive()

}

func (pl *projectList) add(p *project) {

	pl.l = append(pl.l, p)
	if pl.a == nil {
		pl.a = p
	}

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

func newProject(o *project) *project {

	m := newProjectModel(o)

	p := project{
		id:     0,
		status: psNew,
		o:      o,
		d: &projectData{
			size: m.size(),
			pop:  0,
			try:  0,
			best: 0.0,
		},
		m: m,
	}

	return &p
}

func (p *project) save() {

	if p.status == psNew {
		p.id = app.state.id
		app.state.id++
		p.status = psActive
		p.m.dirty = false
		app.result.l = append(app.result.l, p)

		app.run = false
	}

}

type projectUI struct {
	selected bool
}

type projectStatus int

const (
	psNew projectStatus = iota
	psActive
	psHolded
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
	dirty bool
}

func newProjectModel(o *project) *projectModel {

	m := &projectModel{}
	if o != nil {
		clone(o.m, m)
	}
	m.dirty = true

	return m
}

func (m *projectModel) size() int {

	return 0 // temp
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

func clone(src, dist interface{}) error {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dist)
}
