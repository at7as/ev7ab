package app

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/at7as/ev7ab/pkg/lab"
	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui
var app *application
var err error

type application struct {
	c      config
	s      state
	status applicationStatus

	cfgFile string
	lab     *lab.Lab
	view    int
	run     bool
	exec    bool
	help    bool
	// setup     *itemList
	// setupItem bool
	result   *projectList
	edit     *project
	nodeSize bool
	cursor   projectModelSource
	scroll   projectModelSource
	link     projectModelSource
	linkEdit bool
	invalid  bool
	position int
	offset   int
}

type config struct{}

type state struct {
	modal     *widget
	keybar    *keybarWidget
	statusbar *statusbarWidget
	setup     *setupWidget
	widgets   map[string]*widget
	managers  []gocui.Manager
}

func newApplication(prod lab.Producer, cfgFile string) *application {

	c := config{}

	s := state{
		keybar:    newKeybarWidget(),
		statusbar: newStatusbarWidget(),
		setup:     newSetupWidget(),
		widgets:   make(map[string]*widget),
		managers:  make([]gocui.Manager, 0),
	}

	s.widgets[s.keybar.name] = s.keybar.widget
	s.widgets[s.statusbar.name] = s.statusbar.widget
	s.widgets[s.setup.name] = s.setup.widget

	for _, w := range s.widgets {
		s.managers = append(s.managers, w)
	}

	return &application{
		c:      c,
		s:      s,
		status: appIdle,

		cfgFile: cfgFile,
		lab:     lab.New(prod),
		run:     false,
		exec:    false,
		view:    0,
		help:    false,
		// setup: newItemList([][2]string{
		// 	{"LabFile", cfgFile},
		// 	{"In", "2"},
		// 	{"Out", "2"},
		// 	{"Size", "1000"},
		// 	{"Aggr", "avg"},
		// 	{"Proc", "linear"},
		// 	{"Goal", "false"},
		// 	{"Target", ""},
		// 	{"Limit", ""},
		// 	{"InputFile", ""},
		// }),
		// setupItem: false,
		result:   newProjectList(),
		edit:     nil,
		nodeSize: false,
		cursor:   projectModelSource{0, 0},
		scroll:   projectModelSource{0, 0},
		link:     projectModelSource{0, 0},
		linkEdit: false,
		invalid:  false,
		position: 0,
		offset:   0,
	}
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

	return nil
}

func (a *application) quit(_ *gocui.Gui, _ *gocui.View) error {

	return gocui.ErrQuit
}

func (a *application) toggleHelp(_ *gocui.Gui, _ *gocui.View) error {

	if app.s.modal == nil {

		print("== nil")

		if err = app.openModal(newHelpBox()); err != nil {
			return err
		}

	} else {

		print("!= nil")

		open := app.s.modal.name != "help"

		if err = app.setTabCurrent(app.s.keybar.getTab()); err != nil {
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
	app.s.keybar.setTab(tabSetup)

	return app.setTabCurrent(tabSetup)
}

func (a *application) showResult(_ *gocui.Gui, _ *gocui.View) error {

	// clearWindow(g, false)

	app.s.keybar.setTab(tabResult)
	// app.s.statusbar.setText(gui.CurrentView().Name())

	// _, err := setCurrentViewOnTop(g, "result")

	return nil
}

func (a *application) showEdit(_ *gocui.Gui, _ *gocui.View) error {

	// clearWindow(g, false)

	app.s.keybar.setTab(tabEdit)
	// app.s.statusbar.setText(gui.CurrentView().Name())

	// _, err := setCurrentViewOnTop(g, "edit")

	return nil
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

	if app.s.modal != nil {
		app.s.modal.clean()
	}
	app.s.modal = w

	if err = app.s.modal.Layout(gui); err != nil {
		return err
	}

	if err = app.s.modal.keybinding(); err != nil {
		return err
	}

	if err = app.setCurrent(w.name); err != nil {
		return err
	}

	return nil
}

func (a *application) closeModal() error {

	gui.Cursor = false

	if app.s.modal != nil {
		if err = app.s.modal.clean(); err != nil {
			return err
		}
		app.s.modal = nil
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

type applicationStatus int

const (
	appIdle applicationStatus = iota
	appRun
	appWait
)

// Run ...
func Run(prod lab.Producer, appConfigFile string) {

	if gui, err = gocui.NewGui(gocui.OutputNormal); err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	app = newApplication(prod, appConfigFile)

	// widgets := []*widget{

	// 	newWidget(app.s.keybar, "keybar"),
	// 	// newHelpWindowWidget(),
	// 	// newEditViewWidget(),
	// 	// newResultViewWidget(),
	// 	// newSetupViewWidget(),
	// 	// newWidget(g, newSetupWidget(), "setup"),
	// 	// statusbar
	// }

	gui.SetManager(app.s.managers...)

	app.keybinding()
	for _, w := range app.s.widgets {
		if err = w.keybinding(); err != nil {
			log.Panicln(err)
		}
	}

	gui.InputEsc = true

	gui.Update(func(_ *gocui.Gui) error { return app.showSetup(nil, nil) })

	if err = gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func showRun(g *gocui.Gui) error {

	if err = g.DeleteView("run"); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	maxX, maxY := g.Size()
	v, err := g.SetView("run", -1, maxY-3, maxX, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Frame = false

	fmt.Fprintf(v, "\033[35;4m%s\033[0m", strings.Repeat("█", maxX))

	_, err = g.SetViewOnTop("run")

	return err
}

func hideRun(g *gocui.Gui) error {

	return g.DeleteView("run")
}

func hideInvalid() {

	time.Sleep(1 * time.Second)
	app.invalid = false
	gui.Update(func(g *gocui.Gui) error {
		return nil
	})

}

// State ...
// type State struct {
// 	Setup [][2]string
// 	Lab   []string
// }
