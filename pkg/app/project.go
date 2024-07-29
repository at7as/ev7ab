package app

import (
	"strconv"

	"github.com/at7as/ev7ab/pkg/lab"
)

type project struct {
	id     int
	status projectStatus
	o      *project
	d      *projectStat
	m      *projectModel
	n      *projectModel
	ed     bool
	sel    bool
}

func newProject(o *project) *project {

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
		d: &projectStat{
			size:   0,
			volume: 0,
			gen:    0,
			ev:     0,
			age:    0,
			best:   "",
			goal:   false,
		},
		m:   nil,
		n:   n,
		ed:  false,
		sel: false,
	}

	return &p
}

func (p *project) edit() {

	p.n = cloneModel(p.m)
	p.n.measure()
	p.ed = true

}

func (p *project) validate(m *projectModel) bool {

	for i := range m.model {
		for ii := range m.model[i].stage {
			m.model[i].stage[ii].valid = i == len(m.model)-1
		}
	}

	for i := range m.model {
		for ii := range m.model[i].stage {
			for _, src := range m.model[i].stage[ii].source {
				m.model[src.s].stage[src.n].valid = true
			}
		}
	}

	for i := range m.model {
		for ii := range m.model[i].stage {
			if !m.model[i].stage[ii].valid {
				continue
			}
			m.model[i].stage[ii].valid = m.model[i].stage[ii].size > 0
			if !m.model[i].stage[ii].valid {
				continue
			}
			if i > 0 {
				m.model[i].stage[ii].valid = len(m.model[i].stage[ii].source) > 0
			}
		}
	}

	for _, s := range m.model {
		if len(s.stage) == 0 {
			return false
		}
		for _, n := range s.stage {
			if !n.valid {
				return false
			}
		}
	}

	return true
}

func (p *project) save() {

	if p.status == psNew {
		p.id = app.lab.AddProject([][]lab.Node{})
		app.result.add(p)
	}

	if app.run {
		p.status = psHolded
	} else {
		p.status = psActive
	}

	if !p.validate(p.n) {
		p.status = psInvalid
		app.invalid = true
		go hideInvalid()
	}

	p.m = p.n

	p.ed = false

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

type projectStat struct {
	size   int
	volume int
	gen    int
	ev     int
	age    int
	best   string
	goal   bool
}

func insertStage() {

	selectStageRight()
	app.edit.n.addStage(app.cursor.s)

}

func deleteStage() {

	app.edit.n.removeStage(app.cursor.s)

}
