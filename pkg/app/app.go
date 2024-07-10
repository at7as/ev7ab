package app

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui

var app application = application{
	state: &State{},
	run:   false,
	view:  0,
	help:  false,
	setup: newItemList([][2]string{
		{"LabFile", "./ev.lab"},
		{"InputFile", ""},
		{"In", "2"},
		{"Out", "2"},
		{"Target", ""},
		{"Limit", ""},
		{"Goal", "false"},
		{"Size", "1000"},
		{"Seed", "0 0"},
		{"Aggr", "avg"},
		{"Proc", "linear"},
	}),
	setupitem: false,
	result:    newProjectList(),
	edit:      nil,
	nodesize:  false,
	cursor:    projectModelSource{0, 0},
	link:      projectModelSource{0, 0},
	linkedit:  false,
	invalid:   false,
	scroll:    projectModelSource{0, 0},
}

type application struct {
	state     *State
	run       bool
	view      int
	help      bool
	setup     *itemList
	setupitem bool
	result    *projectList
	edit      *project
	nodesize  bool
	cursor    projectModelSource
	link      projectModelSource
	linkedit  bool
	invalid   bool
	scroll    projectModelSource
}

// State ...
type State struct {
	id     int
	result *projectList
}

type Cursor struct {
	x, y int
}

// Run ...
func Run() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	gui = g
	defer g.Close()

	g.SetManager(
		NewKeybarWidget(),
		NewHelpWindowWidget(),
		NewSetupViewWidget(),
		NewResultViewWidget(),
		NewEditViewWidget(),
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
