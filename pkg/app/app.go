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

type application struct {
	lab *lab.Lab
	// state     *State
	run       bool
	view      int
	help      bool
	setup     *itemList
	setupItem bool
	result    *projectList
	edit      *project
	nodeSize  bool
	cursor    projectModelSource
	link      projectModelSource
	linkEdit  bool
	invalid   bool
	scroll    projectModelSource
}

func newApplication(prod lab.Producer, labFile string) *application {

	return &application{
		lab: lab.New(prod),
		// state: &State{},
		run:  false,
		view: 0,
		help: false,
		setup: newItemList([][2]string{
			{"LabFile", labFile},
			{"In", "2"},
			{"Out", "2"},
			{"Size", "1000"},
			{"Aggr", "avg"},
			{"Proc", "linear"},
			{"Goal", "false"},
			{"Target", ""},
			{"Limit", ""},
			{"InputFile", ""},
		}),
		setupItem: false,
		result:    newProjectList(),
		edit:      nil,
		nodeSize:  false,
		cursor:    projectModelSource{0, 0},
		link:      projectModelSource{0, 0},
		linkEdit:  false,
		invalid:   false,
		scroll:    projectModelSource{0, 0},
	}
}

func (a *application) save() error {

	s := &State{}
	// app.state.ID = app.lab.s.id
	s.Setup = make([][2]string, len(app.setup.l))
	for i, v := range app.setup.l {
		s.Setup[i][0] = v.key
		s.Setup[i][1] = v.value
	}
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

// Run ...
func Run(prod lab.Producer, labFile string) {

	app = newApplication(prod, labFile)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	gui = g
	defer g.Close()

	g.SetManager(
		newKeybarWidget(),
		newHelpWindowWidget(),
		newSetupViewWidget(),
		newResultViewWidget(),
		newEditViewWidget(),
	)

	g.InputEsc = true
	if err := setKeybinding(g); err != nil {
		log.Panicln(err)
	}

	g.Update(func(g *gocui.Gui) error { return showSetup(g, nil) })

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func showRun(g *gocui.Gui) error {

	if err := g.DeleteView("run"); err != nil && err != gocui.ErrUnknownView {
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
type State struct {
	Setup [][2]string
	Lab   []string
}
