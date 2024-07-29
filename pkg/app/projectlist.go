package app

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

	// for _, v := range app.result.l {
	// 	if v.ui.selected {
	// 		return true
	// 	}
	// }

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
