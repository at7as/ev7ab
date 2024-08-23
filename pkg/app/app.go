package app

import (
	"log"

	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui
var app *application
var err error

type application struct {
	c config
	s state
	v view
}

type config struct {
	file string
}

type state struct {
	status applicationStatus
	lab    *lab.Lab
	list   []*project
}

type view struct {
	modal     *widget
	keybar    *keybarWidget
	statusbar *statusbarWidget
	setup     *setupWidget
	result    *resultWidget
	edit      *editWidget
}

type applicationStatus int

const (
	appIdle applicationStatus = iota
	appRun
	appWait
)

func newApplication(prod lab.Producer, cfgFile string) *application {

	c := config{
		file: cfgFile,
	}

	s := state{
		status: appIdle,
		lab:    lab.New(prod),
		list:   make([]*project, 0),
	}

	v := view{
		keybar:    newKeybarWidget(),
		statusbar: newStatusbarWidget(),
		setup:     newSetupWidget(),
		result:    newResultWidget(),
		edit:      newEditWidget(),
	}

	return &application{c, s, v}
}

func (a *application) keybinding() error {

	if err = gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, app.quit); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, app.toggleHelp); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF2, gocui.ModNone, app.showSetup); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF3, gocui.ModNone, app.showResult); err != nil {
		return err
	}

	if err = gui.SetKeybinding("", gocui.KeyF4, gocui.ModNone, app.showEdit); err != nil {
		return err
	}

	if err = app.v.keybar.keybinding(); err != nil {
		return err
	}

	if err = app.v.statusbar.keybinding(); err != nil {
		return err
	}

	if err = app.v.setup.keybinding(); err != nil {
		return err
	}

	if err = app.v.result.keybinding(); err != nil {
		return err
	}

	if err = app.v.edit.keybinding(); err != nil {
		return err
	}

	return nil
}

func (a *application) quit(_ *gocui.Gui, _ *gocui.View) error {

	return gocui.ErrQuit
}

func (a *application) toggleHelp(_ *gocui.Gui, _ *gocui.View) error {

	if app.v.modal == nil {

		if err = app.openModal(newHelpBox()); err != nil {
			return err
		}

	} else {

		open := app.v.modal.name != "help"

		if err = app.setTabCurrent(app.v.keybar.tab); err != nil {
			return err
		}

		if err = app.closeModal(); err != nil {
			return err
		}

		if open {
			if err = app.openModal(newHelpBox()); err != nil {
				return err
			}
		}

	}

	return nil
}

func (a *application) showSetup(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.closeModal(); err != nil {
		return err
	}
	app.v.keybar.setTab(tabSetup)

	return app.setTabCurrent(tabSetup)
}

func (a *application) showResult(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.closeModal(); err != nil {
		return err
	}
	app.v.keybar.setTab(tabResult)

	return app.setTabCurrent(tabResult)
}

func (a *application) showEdit(_ *gocui.Gui, _ *gocui.View) error {

	if err = app.closeModal(); err != nil {
		return err
	}
	app.v.keybar.setTab(tabEdit)

	return app.setTabCurrent(tabEdit)
}

func (a *application) idle() bool {

	return a.s.status == appIdle
}

func (a *application) setCurrent(name string) error {

	if _, err = gui.SetCurrentView(name); err != nil {
		return err
	}
	if _, err = gui.SetViewOnTop(name); err != nil {
		return err
	}

	return nil
}

func (a *application) setTabCurrent(tab keybarTab) error {

	if _, err = gui.SetCurrentView(string(tab)); err != nil {
		return err
	}
	if _, err = gui.SetViewOnTop(string(tab)); err != nil {
		return err
	}

	return nil
}

func (a *application) openModal(w *widget) error {

	if app.v.modal != nil {
		app.v.modal.clean()
	}
	app.v.modal = w

	if err = app.v.modal.Layout(gui); err != nil {
		return err
	}

	if err = app.v.modal.keybinding(); err != nil {
		return err
	}

	if err = app.setCurrent(w.name); err != nil {
		return err
	}

	return nil
}

func (a *application) closeModal() error {

	gui.Cursor = false

	if app.v.modal != nil {
		if err = app.v.modal.clean(); err != nil {
			return err
		}
		app.v.modal = nil
	}

	return nil
}

func (a *application) save() error {

	// s := &State{}
	// app.state.ID = app.lab.s.id
	// s.Setup = make([][2]string, len(app.setup.l))
	// for i, v := range app.setup.l {
	// 	s.Setup[i][0] = v.key
	// 	s.Setup[i][1] = v.value
	// }
	// prepare result

	// save to file

	return nil
}

func (a *application) load() error {

	// load to state

	// app.id = app.state.ID
	// app.setup = newItemList(app.state.Setup)
	// app.lab.Setup(app.state.Setup)

	// prepare result

	return nil
}

func (a *application) apply(s applicationStatus, t string) {

	a.s.status = s
	a.v.statusbar.setColor(s)
	a.v.statusbar.setText(t)

}

// Run ...
func Run(prod lab.Producer, appConfigFile string) {

	if gui, err = gocui.NewGui(gocui.OutputNormal); err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	app = newApplication(prod, appConfigFile)

	gui.SetManager(
		app.v.keybar.widget,
		app.v.statusbar.widget,
		app.v.setup.widget,
		app.v.result.widget,
		app.v.edit.widget,
	)

	app.keybinding()

	gui.InputEsc = true

	gui.Update(func(_ *gocui.Gui) error { return app.showSetup(nil, nil) })

	if err = gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
