package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

// total (results)
// config
// state
// projects and his props
// project = model + []entity

var app application = application{
	state: &State{},
	run:   false,
	view:  0,
	help:  false,
	setup: newItemList([][2]string{
		{"LabFile", "./ev.lab"},
		{"InputFile", ""},
		{"In", ""},
		{"Out", ""},
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
}

// State ...
type State struct {
	id int
}

func init() {
	// temp
	app.result.add(&project{0, psActive, projectUI{}, nil, &projectData{4, 1000, 2154, 0.25}, nil, nil})
	app.result.add(&project{3, psActive, projectUI{}, nil, &projectData{12, 1000, 123, 0.33}, nil, nil})
	app.result.add(&project{11, psActive, projectUI{}, nil, &projectData{64, 1000, 54, 0.48}, nil, nil})
}

// Run ...
func Run() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManager(
		NewKeybarWidget(),
		NewHelpWindowWidget(),
		NewSetupViewWidget(),
		NewResultViewWidget(),
		NewEditViewWidget(),
		// NewRunWidget(),
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

// // RunWidget ...
// type RunWidget struct{}

// // NewRunWidget ...
// func NewRunWidget() *RunWidget {
// 	return &RunWidget{}
// }

// // Layout ...
// func (w *RunWidget) Layout(g *gocui.Gui) error {

// 	maxX, maxY := g.Size()
// 	v, err := g.SetView("run", -1, maxY-3, maxX, maxY-1)
// 	if err != nil && err != gocui.ErrUnknownView {
// 		return err
// 	}

// 	v.Frame = false
// 	fmt.Fprintf(v, "\033[35;4m%s\033[0m", strings.Repeat("█", maxX))

// 	return nil
// }

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
