package app

import (
	"strconv"

	"github.com/at7as/ev7ab/pkg/lab"
)

type project struct {
	id     int
	status projectStatus
	origin *project
	stat
	*model
	draft    *model
	selected bool
}

type projectStatus int

const (
	projectEdit projectStatus = iota
	projectActive
	projectHolded
	projectInvalid
	projectTerminated
)

type stat struct {
	size   int
	volume int
	gen    int
	ev     int
	age    int
	best   string
	goal   bool
}

func newProject(o *project) *project {

	s := stat{
		size:   0,
		volume: 0,
		gen:    0,
		ev:     0,
		age:    0,
		best:   "",
		goal:   false,
	}

	in := 0
	out := 0
	for _, obj := range app.v.setup.list {
		if obj.key == "In" {
			in, _ = strconv.Atoi(obj.value)
		}
		if obj.key == "Out" {
			out, _ = strconv.Atoi(obj.value)
		}
	}

	return &project{
		id:       app.s.lab.AddProject([][]lab.Node{}),
		status:   projectEdit,
		origin:   o,
		stat:     s,
		model:    nil,
		draft:    newProjectModel(o, in, out),
		selected: false,
	}
}

func (p *project) edit() {

	p.draft = p.model.clone()
	p.draft.measure()
	p.status = projectEdit

}

func (p *project) cancel() {

	if !p.model.validate() {
		p.status = projectInvalid
	} else {
		p.status = projectHolded
	}

	p.draft = p.model
	p.model.measure()

}

func (p *project) validate() bool {

	return p.draft.validate()
}

func (p *project) measure() {

	p.draft.measure()
}

func (p *project) save() {

	if p.status == projectEdit {

		if !p.validate() {
			p.status = projectInvalid
			app.v.edit.setInvalid(true)
			go app.v.edit.hideInvalid()
		} else if app.idle() {
			p.status = projectActive
		} else {
			p.status = projectHolded
		}

		if p.model == nil {
			app.s.list = append(app.s.list, p)
		}

		p.model = p.draft

	}

}

func (p *project) setNodeSize(x, y, v int) {

	p.draft.model[x].model[y].size = v

}

func (p *project) getNodeSize(x, y int) int {

	return p.draft.model[x].model[y].size
}

func (p *project) getModel() []*model {

	return p.draft.model
}

func (p *project) getStage(x int) *model {

	return p.draft.model[x]
}

func (p *project) getNode(pos position) *model {

	return p.draft.model[pos.x].model[pos.y]
}

func (p *project) getSize() int {

	return p.draft.size
}

func (p *project) getVolume() int {

	return p.draft.volume
}

func (p *project) insertStage(index int) {

	p.draft.insertStage(index)

}

func (p *project) deleteStage(index int) {

	p.draft.deleteStage(index)

}

func (ps projectStatus) text() string {

	switch ps {
	case 0:
		return "edit"
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
