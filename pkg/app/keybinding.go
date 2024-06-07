package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func setKeybinding(g *gocui.Gui) error {

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

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

	if err := g.SetKeybinding("setup", gocui.KeyF4, gocui.ModNone, showEdit); err != nil {
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

	if err := g.SetKeybinding("result", gocui.KeyCtrlD, gocui.ModNone, onResultDubProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyF4, gocui.ModNone, onResultEditProject); err != nil {
		return err
	}

	if err := g.SetKeybinding("result", gocui.KeyCtrlE, gocui.ModNone, onResultEditProject); err != nil {
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

	if err := g.SetKeybinding("result", gocui.KeyCtrlTilde, gocui.ModNone, onResultShowToggle); err != nil {
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
	if app.setupitem {
		closeSetupItem(g, nil)
	}

}

func toggleAppRun(g *gocui.Gui) {

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

	if app.edit == nil {
		app.edit = newProject(nil)
	}

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

func onSetupKeyArrowUp(g *gocui.Gui, v *gocui.View) error {

	app.setup.selectUp()

	return nil
}

func onSetupKeyArrowDown(g *gocui.Gui, v *gocui.View) error {

	app.setup.selectDown()

	return nil
}

func openSetupItem(g *gocui.Gui, v *gocui.View) error {

	if err := g.DeleteView("setupitem"); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	app.setupitem = true
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

	app.setupitem = false
	g.Cursor = false

	if err := g.DeleteView("setupitem"); err != nil {
		return err
	}

	_, err := setCurrentViewOnTop(g, "setup")

	return err
}

func closeEnterSetupItem(g *gocui.Gui, v *gocui.View) error {

	app.setupitem = false
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

	return err
}

func onResultRunToggle(g *gocui.Gui, v *gocui.View) error {

	toggleAppRun(g)

	return nil
}

func onResultKeyArrowUp(g *gocui.Gui, v *gocui.View) error {

	app.result.selectUp()

	return nil
}

func onResultKeyArrowDown(g *gocui.Gui, v *gocui.View) error {

	app.result.selectDown()

	return nil
}

func onResultNewProject(g *gocui.Gui, v *gocui.View) error {

	app.edit = newProject(nil)

	return showEdit(g, v)
}

func onResultDubProject(g *gocui.Gui, v *gocui.View) error {

	if len(app.result.l) > 0 {
		app.edit = newProject(app.result.a)
	} else {
		return nil
	}

	return showEdit(g, v)
}

func onResultEditProject(g *gocui.Gui, v *gocui.View) error {

	if app.edit == nil {
		if !app.run && len(app.result.l) > 0 {
			app.edit = app.result.a
		} else {
			app.edit = newProject(nil)
		}
	}

	return showEdit(g, v)
}

func onResultSelectProject(g *gocui.Gui, v *gocui.View) error {

	if !app.run {
		// if len(app.result.l)
		// app.result.a
	}

	return nil
}

func onResultSelectAllProject(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func onResultHoldProject(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func onResultActivateProject(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func onResultTerminateProject(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func onResultShowToggle(g *gocui.Gui, v *gocui.View) error {

	return nil
}

// select project
// do with selected
// hold ctrlH
// activate ctrlA
// show all ctrlW
// delete holded (delete all ents) ctrlT (terminate)

// what do when IL or PL is changed?!?
