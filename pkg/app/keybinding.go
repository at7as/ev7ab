package app

import (
	"fmt"
	"strconv"

	"github.com/jroimartin/gocui"
)

func setKeybinding(g *gocui.Gui) error {

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, toggleHelp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, showSetup); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyF3, gocui.ModNone, showResult); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyF4, gocui.ModNone, showEdit); err != nil {
		return err
	}

	if err := g.SetKeybinding("setup", gocui.KeyCtrlS, gocui.ModNone, onSetupSave); err != nil {
		return err
	}

	if err := g.SetKeybinding("setup", gocui.KeyCtrlL, gocui.ModNone, onSetupLoad); err != nil {
		return err
	}

	if err := g.SetKeybinding("setup", gocui.KeyArrowUp, gocui.ModNone, onSetupKeyArrowUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("setup", gocui.KeyArrowDown, gocui.ModNone, onSetupKeyArrowDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("setup", gocui.KeyEnter, gocui.ModNone, openSetupItem); err != nil {
		return err
	}

	if err := g.SetKeybinding("setupitem", gocui.KeyEsc, gocui.ModNone, closeSetupItem); err != nil {
		return err
	}

	if err := g.SetKeybinding("setupitem", gocui.KeyEnter, gocui.ModNone, closeEnterSetupItem); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyEnter, gocui.ModNone, onResultRunToggle); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyArrowUp, gocui.ModNone, onResultKeyArrowUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyArrowDown, gocui.ModNone, onResultKeyArrowDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlN, gocui.ModNone, onResultNewProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlE, gocui.ModNone, onResultEditProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlD, gocui.ModNone, onResultDubProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeySpace, gocui.ModNone, onResultSelectProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlSpace, gocui.ModNone, onResultSelectAllProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlH, gocui.ModNone, onResultHoldProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlA, gocui.ModNone, onResultActivateProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlT, gocui.ModNone, onResultTerminateProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyTab, gocui.ModNone, onResultShowToggle); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyCtrlN, gocui.ModNone, onEditNewProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyCtrlE, gocui.ModNone, onEditEditProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyCtrlD, gocui.ModNone, onEditDubProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyCtrlS, gocui.ModNone, onEditSaveProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyCtrlV, gocui.ModNone, onEditValidateProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyArrowLeft, gocui.ModNone, onEditKeyArrowLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyArrowRight, gocui.ModNone, onEditKeyArrowRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyArrowUp, gocui.ModNone, onEditKeyArrowUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyArrowDown, gocui.ModNone, onEditKeyArrowDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyInsert, gocui.ModNone, onEditInsertStage); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyDelete, gocui.ModNone, onEditDeleteStage); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyEnter, gocui.ModNone, onEditInsertNode); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeySpace, gocui.ModNone, onEditEditNodeSource); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyCtrlSpace, gocui.ModNone, onEditEditNode); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyBackspace, gocui.ModNone, onEditDeleteNode); err != nil {
		return err
	}

	if err := g.SetKeybinding("edit", gocui.KeyEsc, gocui.ModNone, onEditEsc); err != nil {
		return err
	}

	if err := g.SetKeybinding("nodesize", gocui.KeyEsc, gocui.ModNone, closeNodeSize); err != nil {
		return err
	}

	if err := g.SetKeybinding("nodesize", gocui.KeyEnter, gocui.ModNone, closeEnterNodeSize); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {

	return gocui.ErrQuit
}

func clearWindow(g *gocui.Gui, help bool) {

	g.Cursor = false
	if app.help && !help {
		toggleHelp(g, nil)
	}
	if app.setupItem {
		closeSetupItem(g, nil)
	}
	if app.nodeSize {
		closeNodeSize(g, nil)
	}

}

func toggleAppRun(g *gocui.Gui) {

	if !app.run && app.result.activated() == 0 {
		return
	}

	app.run = !app.run

	if app.run {
		showRun(g)
	} else {
		hideRun(g)
	}

}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {

	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	v, err := g.SetViewOnTop(name)
	if err != nil {
		return nil, err
	}

	if app.help {
		_, err = g.SetViewOnTop("help")
		if err != nil {
			return nil, err
		}
	}

	if app.run {
		_, err = g.SetViewOnTop("run")
		if err != nil {
			return nil, err
		}
	}

	return v, err
}

func showSetup(g *gocui.Gui, v *gocui.View) error {

	clearWindow(g, false)

	app.view = 0
	_, err := setCurrentViewOnTop(g, "setup")

	return err
}

func showResult(g *gocui.Gui, v *gocui.View) error {

	clearWindow(g, false)

	app.view = 1
	_, err := setCurrentViewOnTop(g, "result")

	return err
}

func showEdit(g *gocui.Gui, v *gocui.View) error {

	clearWindow(g, false)

	app.view = 2
	_, err := setCurrentViewOnTop(g, "edit")

	return err
}

func toggleHelp(g *gocui.Gui, v *gocui.View) error {

	clearWindow(g, true)

	app.help = !app.help
	if app.help {
		if _, err := g.SetViewOnTop("help"); err != nil {
			return err
		}
	} else {
		if _, err := g.SetViewOnBottom("help"); err != nil {
			return err
		}
	}

	return nil
}

func onSetupSave(g *gocui.Gui, v *gocui.View) error {

	return app.save()
}

func onSetupLoad(g *gocui.Gui, v *gocui.View) error {

	return app.load()
}

func onSetupKeyArrowUp(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	app.setup.selectUp()

	return nil
}

func onSetupKeyArrowDown(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	app.setup.selectDown()

	return nil
}

func openSetupItem(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	if err := g.DeleteView("setupitem"); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	app.setupItem = true
	g.Cursor = !app.help

	maxX, maxY := g.Size()
	view, err := g.SetView("setupitem", maxX/2-30, maxY/2-2, maxX/2+30, maxY/2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	view.Editable = true
	view.Frame = true
	view.Title = fmt.Sprintf(" %v ", app.setup.a.key)
	fmt.Fprintf(view, "%v", app.setup.a.value)
	if err := view.SetCursor(len(app.setup.a.value), 0); err != nil {
		return err
	}

	_, err = setCurrentViewOnTop(g, "setupitem")

	return err
}

func closeSetupItem(g *gocui.Gui, v *gocui.View) error {

	app.setupItem = false
	g.Cursor = false

	if err := g.DeleteView("setupitem"); err != nil {
		return err
	}

	_, err := setCurrentViewOnTop(g, "setup")

	return err
}

func closeEnterSetupItem(g *gocui.Gui, v *gocui.View) error {

	app.setupItem = false
	g.Cursor = false

	value, err := v.Line(0)
	if err != nil {
		value = ""
	}
	app.setup.l[app.setup.ai].value = value

	if err := g.DeleteView("setupitem"); err != nil {
		return err
	}

	_, err = setCurrentViewOnTop(g, "setup")

	app.lab.Setup([][2]string{{app.setup.l[app.setup.ai].key, value}})

	return err
}

func onResultRunToggle(g *gocui.Gui, v *gocui.View) error {

	toggleAppRun(g)

	return nil
}

func onResultKeyArrowUp(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	app.result.selectUp(1)

	return nil
}

func onResultKeyArrowDown(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	app.result.selectDown(1)

	return nil
}

func onResultNewProject(g *gocui.Gui, v *gocui.View) error {

	app.edit = newProject(nil)
	app.cursor = projectModelSource{0, 0}
	app.scroll = projectModelSource{0, 0}

	return showEdit(g, v)
}

func onResultEditProject(g *gocui.Gui, v *gocui.View) error {

	if app.result.a != nil {
		app.edit = app.result.a
		app.edit.edit()
		app.cursor = projectModelSource{0, 0}
		app.scroll = projectModelSource{0, 0}
		return showEdit(g, v)
	}

	return nil
}

func onResultDubProject(g *gocui.Gui, v *gocui.View) error {

	if app.result.a != nil {
		app.edit = newProject(app.result.a)
		app.cursor = projectModelSource{0, 0}
		app.scroll = projectModelSource{0, 0}
		return showEdit(g, v)
	}

	return nil
}

func onResultSelectProject(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	if !app.run && app.result.a != nil {
		app.result.a.sel = !app.result.a.sel
	}

	return nil
}

func onResultSelectAllProject(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	selected := !app.result.selected()
	for i := range app.result.l {
		app.result.l[i].sel = selected
	}

	return nil
}

func onResultHoldProject(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	if app.result.selected() {
		for i, v := range app.result.l {
			if v.status == psActive && v.sel {
				app.result.l[i].status = psHolded
			}
		}
	} else {
		if app.result.a.status == psActive {
			app.result.a.status = psHolded
		}
	}

	app.result.next()

	return nil
}

func onResultActivateProject(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	if !app.result.holded {
		return nil
	}

	if app.result.selected() {
		for i, v := range app.result.l {
			if v.status == psHolded && v.sel {
				app.result.l[i].status = psActive
			}
		}
	} else {
		if app.result.a.status == psHolded {
			app.result.a.status = psActive
		}
	}

	return nil
}

func onResultTerminateProject(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	if app.result.selected() {
		for i, v := range app.result.l {
			if v.sel {
				app.result.l[i].sel = false
				app.result.l[i].status = psTerminated
			}
		}
	} else {
		app.result.a.sel = false
		app.result.a.status = psTerminated
	}

	app.result.next()

	return nil
}

func onResultShowToggle(g *gocui.Gui, v *gocui.View) error {

	if app.run {
		return nil
	}

	app.result.holded = !app.result.holded

	app.result.next()

	return nil
}

func onEditNewProject(g *gocui.Gui, v *gocui.View) error {

	app.linkEdit = false
	app.edit = newProject(nil)
	app.cursor = projectModelSource{0, 0}
	app.scroll = projectModelSource{0, 0}

	return nil
}

func onEditEditProject(g *gocui.Gui, v *gocui.View) error {

	app.linkEdit = false

	if app.edit != nil {
		app.edit.edit()
	}

	return nil
}

func onEditDubProject(g *gocui.Gui, v *gocui.View) error {

	app.linkEdit = false

	if app.result.a != nil {
		app.edit = newProject(app.result.a)
		app.cursor = projectModelSource{0, 0}
		app.scroll = projectModelSource{0, 0}
	}

	return nil
}

func onEditSaveProject(g *gocui.Gui, v *gocui.View) error {

	app.linkEdit = false

	if app.edit != nil {
		app.edit.save()
	}

	return nil
}

func onEditValidateProject(g *gocui.Gui, v *gocui.View) error {

	app.linkEdit = false

	if app.edit != nil {
		if !app.edit.validate(app.edit.n) {
			app.invalid = true
			go hideInvalid()
		}
	}

	return nil
}

func onEditKeyArrowLeft(g *gocui.Gui, v *gocui.View) error {

	selectStageLeft()
	shiftCursorY()

	return nil
}

func onEditKeyArrowRight(g *gocui.Gui, v *gocui.View) error {

	selectStageRight()
	shiftCursorY()

	return nil
}

func onEditKeyArrowUp(g *gocui.Gui, v *gocui.View) error {

	selectNodeUp()

	return nil
}

func onEditKeyArrowDown(g *gocui.Gui, v *gocui.View) error {

	selectNodeDown()

	return nil
}

func onEditInsertStage(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	app.linkEdit = false
	insertStage()
	shiftCursorY()

	return nil
}

func onEditDeleteStage(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	app.linkEdit = false
	deleteStage()
	shiftCursorY()
	app.edit.n.measure()

	return nil
}

func onEditInsertNode(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	app.linkEdit = false

	if app.cursor.s > 0 && app.cursor.s < len(app.edit.n.model)-1 {
		app.edit.n.model[app.cursor.s].addNode()
		app.cursor.n = len(app.edit.n.model[app.cursor.s].stage) - 1
		return openNodeSize(g, v)
	}

	app.edit.n.measure()

	return nil
}

func onEditEditNodeSource(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	if app.linkEdit {
		if app.cursor.s != app.link.s || app.cursor.n != app.link.n {

			found := false
			for i, src := range app.edit.n.model[app.cursor.s].stage[app.cursor.n].source {
				if src.s == app.link.s && src.n == app.link.n {
					app.edit.n.model[app.cursor.s].stage[app.cursor.n].source = append(app.edit.n.model[app.cursor.s].stage[app.cursor.n].source[:i], app.edit.n.model[app.cursor.s].stage[app.cursor.n].source[i+1:]...)
					found = true
					break
				}
			}

			if !found && app.link.s < app.cursor.s && len(app.edit.n.model[app.link.s].stage) > app.link.n {
				app.edit.n.model[app.cursor.s].stage[app.cursor.n].source = append(app.edit.n.model[app.cursor.s].stage[app.cursor.n].source, projectModelSource{s: app.link.s, n: app.link.n})
			}

		}
	} else {
		if len(app.edit.n.model[app.cursor.s].stage) == 0 {
			return nil
		}
		app.link.s = app.cursor.s
		app.link.n = app.cursor.n
	}

	app.linkEdit = !app.linkEdit
	app.edit.n.measure()

	return nil
}

func onEditEditNode(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	app.linkEdit = false
	app.edit.n.measure()

	return openNodeSize(g, v)
}

func onEditDeleteNode(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	app.linkEdit = false

	if app.cursor.s > 0 && app.cursor.s < len(app.edit.n.model)-1 {
		app.edit.n.model[app.cursor.s].removeNode(app.cursor.n)
		for i := range app.edit.n.model {
			for ii := range app.edit.n.model[i].stage {
				for iii, src := range app.edit.n.model[i].stage[ii].source {
					if src.s == app.cursor.s && src.n == app.cursor.n {
						app.edit.n.model[i].stage[ii].source = append(app.edit.n.model[i].stage[ii].source[:iii], app.edit.n.model[i].stage[ii].source[iii+1:]...)
						break
					}
				}
			}
		}

	}
	shiftCursorY()
	app.edit.n.measure()

	return nil
}

func openNodeSize(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		return nil
	} else {
		if !app.edit.ed {
			return nil
		}
	}

	app.linkEdit = false

	if err := g.DeleteView("nodesize"); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	app.nodeSize = true
	g.Cursor = !app.help

	maxX, maxY := g.Size()
	view, err := g.SetView("nodesize", maxX/2-30, maxY/2-2, maxX/2+30, maxY/2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	view.Editable = true
	view.Frame = true
	view.Title = " Size "
	value := strconv.Itoa(app.edit.n.model[app.cursor.s].stage[app.cursor.n].size)
	fmt.Fprintf(view, "%v", value)
	if err := view.SetCursor(len(value), 0); err != nil {
		return err
	}

	_, err = setCurrentViewOnTop(g, "nodesize")

	return err
}

func closeNodeSize(g *gocui.Gui, v *gocui.View) error {

	app.nodeSize = false
	g.Cursor = false

	if err := g.DeleteView("nodesize"); err != nil {
		return err
	}

	_, err := setCurrentViewOnTop(g, "edit")

	return err
}

func closeEnterNodeSize(g *gocui.Gui, v *gocui.View) error {

	app.nodeSize = false
	g.Cursor = false

	value, err := v.Line(0)
	if err != nil {
		value = ""
	}
	size, err := strconv.Atoi(value)
	if err != nil {
		size = 0
	}
	app.edit.n.model[app.cursor.s].stage[app.cursor.n].size = size
	app.edit.n.measure()

	if err := g.DeleteView("nodesize"); err != nil {
		return err
	}

	_, err = setCurrentViewOnTop(g, "edit")

	return err
}

func onEditEsc(g *gocui.Gui, v *gocui.View) error {

	app.linkEdit = false

	return nil
}
