package app

// type projectList struct {
// 	list   []*project
// 	show   []*project
// 	a      *project
// 	ai     int
// 	holded bool
// }

// func newProjectList() *projectList {

// 	return &projectList{list: make([]*project, 0), show: make([]*project, 0)}
// }

// func (pl *projectList) add(p *project) {

// 	pl.list = append(pl.list, p)
// 	if pl.a == nil {
// 		pl.a = p
// 		pl.ai = 0
// 	}

// }

// func (pl *projectList) assignActive(ai int) {

// 	pl.ai = ai
// 	if len(pl.list) > 0 {
// 		pl.a = pl.list[pl.ai]
// 	}

// }

// func (pl *projectList) selectUp(start int) {

// 	ai := pl.ai - start
// 	if len(pl.list) > 0 {
// 		for ai >= 0 {
// 			if !pl.holded && pl.list[ai].status == psActive {
// 				pl.assignActive(ai)
// 				break
// 			}
// 			if pl.holded && pl.list[ai].status != psTerminated {
// 				pl.assignActive(ai)
// 				break
// 			}
// 			ai--
// 		}
// 	}

// }

// func (pl *projectList) selectDown(start int) {

// 	ai := pl.ai + start
// 	for ai < len(pl.list) {
// 		if !pl.holded && pl.list[ai].status == psActive {
// 			pl.assignActive(ai)
// 			break
// 		}
// 		if pl.holded && pl.list[ai].status != psTerminated {
// 			pl.assignActive(ai)
// 			break
// 		}
// 		ai++
// 	}

// }

// func (pl *projectList) next() {

// 	ai := pl.ai
// 	pl.selectDown(0)
// 	if ai != pl.ai {
// 		return
// 	}
// 	pl.selectUp(0)

// }

// func (pl *projectList) selected() bool {

// for _, v := range app.result.l {
// 	if v.ui.selected {
// 		return true
// 	}
// }

// 	return false
// }

// func (pl *projectList) count() int {

// 	cnt := 0
// for _, v := range app.result.list {
// 	if v.status == psActive || v.status == psHolded {
// 		cnt++
// 	}
// }

// 	return cnt
// }

// func (pl *projectList) activated() int {

// 	cnt := 0
// for _, v := range app.result.list {
// 	if v.status == psActive {
// 		cnt++
// 	}
// }

// 	return cnt
// }

// func (pl *projectList) update() {

// l := make([]*project, 0, len(app.result.list))
// for _, p := range app.result.list {
// 	if (!app.result.holded && p.status == psActive) || (app.result.holded && p.status != psTerminated) {
// 		l = append(l, p)
// 	}
// }
// app.result.show = l
// goal first
// running first

// }
