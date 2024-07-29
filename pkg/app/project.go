package app

import (
	"strconv"

	"github.com/at7as/ev7ab/pkg/lab"
)

type project struct {
	id     int
	status projectStatus
	ui     projectUI
	o      *project
	d      *projectData
	m      *projectModel
	n      *projectModel
	ed     bool
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
	p.n.measure()
	p.ed = true

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
		p.id = app.lab.AddProject([][]lab.Node{})
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

type projectData struct {
	size int
	pop  int
	try  int
	best float64
}

func insertStage() {

	selectStageRight()
	app.edit.n.addStage(app.cursor.s)

}

func deleteStage() {

	app.edit.n.removeStage(app.cursor.s)

}
